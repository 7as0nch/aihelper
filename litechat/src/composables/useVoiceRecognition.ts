import { ref, onUnmounted } from 'vue';

// Web Speech API Type declarations
interface SpeechRecognitionEvent extends Event {
    resultIndex: number;
    results: SpeechRecognitionResultList;
}

interface SpeechRecognitionErrorEvent extends Event {
    error: string;
    message?: string;
}

interface SpeechRecognitionResultList {
    length: number;
    item(index: number): SpeechRecognitionResult;
    [index: number]: SpeechRecognitionResult;
}

interface SpeechRecognitionResult {
    isFinal: boolean;
    length: number;
    item(index: number): SpeechRecognitionAlternative;
    [index: number]: SpeechRecognitionAlternative;
}

interface SpeechRecognitionAlternative {
    transcript: string;
    confidence: number;
}

interface SpeechRecognition extends EventTarget {
    continuous: boolean;
    interimResults: boolean;
    lang: string;
    maxAlternatives: number;
    onresult: ((event: SpeechRecognitionEvent) => void) | null;
    onerror: ((event: SpeechRecognitionErrorEvent) => void) | null;
    onend: (() => void) | null;
    onstart: (() => void) | null;
    start(): void;
    stop(): void;
    abort(): void;
}

interface SpeechRecognitionConstructor {
    new(): SpeechRecognition;
}

declare global {
    interface Window {
        SpeechRecognition?: SpeechRecognitionConstructor;
        webkitSpeechRecognition?: SpeechRecognitionConstructor;
    }
}

export interface UseVoiceRecognitionOptions {
    lang?: string;
    continuous?: boolean;
    interimResults?: boolean;
}

export interface UseVoiceRecognitionReturn {
    isSupported: boolean;
    isRecording: ReturnType<typeof ref<boolean>>;
    transcript: ReturnType<typeof ref<string>>;
    interimTranscript: ReturnType<typeof ref<string>>;
    error: ReturnType<typeof ref<string | null>>;
    start: () => void;
    stop: () => void;
    setLanguage: (lang: string) => void;
    supportedLanguages: readonly { code: string; label: string }[];
}

// Supported languages for speech recognition
const SUPPORTED_LANGUAGES = [
    { code: 'zh-CN', label: '中文（普通话）' },
    { code: 'en-US', label: 'English (US)' },
    { code: 'en-GB', label: 'English (UK)' },
] as const;

// Error message mapping for Chinese UI
const ERROR_MESSAGES: Record<string, string> = {
    'not-allowed': '麦克风权限被拒绝，请在浏览器设置中允许访问麦克风',
    'network': '网络连接失败，语音识别需要网络支持',
    'no-speech': '未检测到语音，请重试',
    'audio-capture': '无法访问麦克风，请检查设备是否正常',
    'aborted': '语音识别已取消',
    'service-not-allowed': '语音服务不可用',
    'language-not-supported': '不支持该语言',
};

export function useVoiceRecognition(options: UseVoiceRecognitionOptions = {}): UseVoiceRecognitionReturn {
    const {
        lang = 'zh-CN',
        continuous = true,
        interimResults = true
    } = options;

    // Check browser support
    const SpeechRecognitionAPI = window.SpeechRecognition || window.webkitSpeechRecognition;
    const isSupported = !!SpeechRecognitionAPI;

    // State
    const isRecording = ref(false);
    const transcript = ref('');
    const interimTranscript = ref('');
    const error = ref<string | null>(null);

    let recognition: SpeechRecognition | null = null;
    let currentLang = lang;

    // Initialize recognition if supported
    const initRecognition = () => {
        if (!SpeechRecognitionAPI) return null;

        const rec = new SpeechRecognitionAPI();
        rec.continuous = continuous;
        rec.interimResults = interimResults;
        rec.lang = currentLang;
        rec.maxAlternatives = 1;

        rec.onstart = () => {
            isRecording.value = true;
            error.value = null;
        };

        rec.onresult = (event: SpeechRecognitionEvent) => {
            let finalTranscript = '';
            let interim = '';

            for (let i = event.resultIndex; i < event.results.length; i++) {
                const result = event.results[i];
                if (result.isFinal) {
                    finalTranscript += result[0].transcript;
                } else {
                    interim += result[0].transcript;
                }
            }

            if (finalTranscript) {
                transcript.value += finalTranscript;
            }
            interimTranscript.value = interim;
        };

        rec.onerror = (event: SpeechRecognitionErrorEvent) => {
            const errorMessage = ERROR_MESSAGES[event.error] || `语音识别错误: ${event.error}`;
            error.value = errorMessage;
            isRecording.value = false;
        };

        rec.onend = () => {
            isRecording.value = false;
            // Clear interim transcript when recognition ends
            interimTranscript.value = '';
        };

        return rec;
    };

    const start = () => {
        if (!isSupported) {
            error.value = '您的浏览器不支持语音识别。请使用 Chrome、Edge 或 Safari 浏览器。';
            return;
        }

        // Reset state
        transcript.value = '';
        interimTranscript.value = '';
        error.value = null;

        // Create new recognition instance each time
        recognition = initRecognition();

        if (recognition) {
            try {
                recognition.start();
            } catch (e) {
                error.value = '启动语音识别失败，请重试';
                console.error('Speech recognition start error:', e);
            }
        }
    };

    const stop = () => {
        if (recognition) {
            try {
                recognition.stop();
            } catch (e) {
                console.error('Speech recognition stop error:', e);
            }
            recognition = null;
        }
        isRecording.value = false;
    };

    const setLanguage = (newLang: string) => {
        currentLang = newLang;
        if (recognition) {
            recognition.lang = newLang;
        }
    };

    // Cleanup on unmount
    onUnmounted(() => {
        if (recognition) {
            try {
                recognition.abort();
            } catch (e) {
                // Ignore cleanup errors
            }
            recognition = null;
        }
    });

    return {
        isSupported,
        isRecording,
        transcript,
        interimTranscript,
        error,
        start,
        stop,
        setLanguage,
        supportedLanguages: SUPPORTED_LANGUAGES,
    };
}
