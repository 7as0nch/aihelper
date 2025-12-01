import { reactive } from 'vue';

export interface RuntimeConfig {
    VITE_APP_TITLE?: string;
    VITE_APP_LOGO?: string;
    VITE_API_BASE_URL?: string;
    VITE_OPENAI_API_KEY?: string;
    VITE_OPENAI_BASE_URL?: string;
    VITE_OPENAI_MODEL?: string;
    VITE_AI_TYPE?: string;
    VITE_ENABLE_QR_LOGIN?: string;
    VITE_BASE_URL?: string;
    VITE_FLOAT_BALL_IMAGE?: string;
    VITE_FLOAT_BALL_WELCOME_CONTENTS?: string;
    [key: string]: string | undefined;
}

// Default config from build-time env
const defaultConfig: RuntimeConfig = {
    VITE_APP_TITLE: import.meta.env.VITE_APP_TITLE,
    VITE_APP_LOGO: import.meta.env.VITE_APP_LOGO,
    VITE_BASE_URL: import.meta.env.VITE_BASE_URL,

    // Sensitive config: Only load from env in development mode
    // In production, these must be passed via initAiChat options
    VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL,
    VITE_OPENAI_API_KEY: import.meta.env.VITE_OPENAI_API_KEY,
    VITE_OPENAI_BASE_URL: import.meta.env.VITE_OPENAI_BASE_URL,
    VITE_OPENAI_MODEL: import.meta.env.VITE_OPENAI_MODEL,
    VITE_AI_TYPE: import.meta.env.VITE_AI_TYPE,
    VITE_ENABLE_QR_LOGIN: import.meta.env.VITE_ENABLE_QR_LOGIN,
    VITE_FLOAT_BALL_WELCOME_CONTENTS: import.meta.env.VITE_FLOAT_BALL_WELCOME_CONTENTS,
    VITE_FLOAT_BALL_IMAGE: import.meta.env.VITE_FLOAT_BALL_IMAGE,
};

export const runtimeConfig = reactive<RuntimeConfig>({ ...defaultConfig });

export function setRuntimeConfig(config: Partial<RuntimeConfig>) {
    Object.assign(runtimeConfig, config);
}

export function getConfig(key: keyof RuntimeConfig, defaultValue: string = ''): string {
    return runtimeConfig[key] || defaultValue;
}
