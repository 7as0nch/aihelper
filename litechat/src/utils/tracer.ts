// 埋点

/**
 * AnalyticsTracker Pro
 * 支持：链式调用、拦截器、批量上报、页面/API/按钮埋点
 */
import { v4 as uuidv4 } from 'uuid';
import type { Router, RouteLocationNormalized } from 'vue-router';
import type { App, DirectiveBinding } from 'vue';
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

interface TrackerOptions {
    baseUrl?: string;
    appId?: string;
    userId?: number;
    batchSize?: number;
    wait?: number;
}

interface LogItem {
    appId: string;
    deviceId: string;
    userId: number;
    timestamp: string;
    userAgent: string;
    pageUrl: string;
    type: string;
    data: any;
}

type BeforeHandler = (log: LogItem) => boolean | Promise<boolean>;
type AfterHandler = (logs: LogItem[]) => void;

// Extend Axios config to include metadata
declare module 'axios' {
    interface InternalAxiosRequestConfig {
        metadata?: {
            startTime: number;
        };
    }
}

class AnalyticsTracker {
    public baseUrl: string;
    public appId: string;
    public userId: number;
    public queue: LogItem[];
    public batchSize: number;
    public wait: number;
    public timer: ReturnType<typeof setInterval> | null;
    public beforeHandlers: BeforeHandler[];
    public afterHandlers: AfterHandler[];
    public deviceId: string;
    public pageStartTime: number;

    constructor(options: TrackerOptions = {}) {
        this.baseUrl = options.baseUrl || '/api/log/batch'; // 注意：通常批量接口不同
        this.appId = options.appId || 'vue3-app';
        this.userId = options.userId || 0;

        // --- 批量配置 ---
        this.queue = []; // 缓冲队列
        this.batchSize = options.batchSize || 10; // 积攒多少条发一次
        this.wait = options.wait || 5000; // 多少毫秒发一次
        this.timer = null;

        // 拦截器
        this.beforeHandlers = [];
        this.afterHandlers = [];

        // 基础数据
        this.deviceId = this._getDeviceId();
        this.pageStartTime = Date.now();

        // 初始化定时器和监听器
        this._startTimer();
        this._listenUnload();
    }

    // --- 内部机制 ---

    private _getDeviceId(): string {
        let id = localStorage.getItem('tracker_device_id');
        if (!id) {
            id = uuidv4();
            localStorage.setItem('tracker_device_id', id);
        }
        return id;
    }

    private _formatDate(timestamp: number): string {
        const date = new Date(timestamp);
        const yyyy = date.getFullYear();
        const MM = String(date.getMonth() + 1).padStart(2, '0');
        const dd = String(date.getDate()).padStart(2, '0');
        const HH = String(date.getHours()).padStart(2, '0');
        const mm = String(date.getMinutes()).padStart(2, '0');
        const ss = String(date.getSeconds()).padStart(2, '0');
        return `${yyyy}-${MM}-${dd} ${HH}:${mm}:${ss}`;
    }

    private _getBaseInfo() {
        return {
            appId: this.appId,
            deviceId: this.deviceId,
            userId: this.userId,
            timestamp: this._formatDate(Date.now()), // 注意：这里是加入队列的时间
            userAgent: navigator.userAgent,
            pageUrl: window.location.href,
        };
    }

    // 启动定时轮询
    private _startTimer() {
        this.timer = setInterval(() => {
            this.flush();
        }, this.wait);
    }

    // 监听页面关闭/隐藏，强制上报
    private _listenUnload() {
        const flushHandler = () => this.flush();
        // 页面不可见或关闭时触发
        window.addEventListener('visibilitychange', () => {
            if (document.visibilityState === 'hidden') {
                flushHandler();
            }
        });
        // 兼容部分浏览器关闭
        window.addEventListener('beforeunload', flushHandler);
    }

    // --- 核心逻辑 ---

    /**
     * 收集日志（并不立即发送，而是入队）
     */
    public async report(type: string, payload: any = {}): Promise<this> {
        const logItem: LogItem = {
            ...this._getBaseInfo(),
            type,
            data: payload,
        };

        // 1. 执行 Before Handlers (针对单条数据)
        for (const handler of this.beforeHandlers) {
            const shouldContinue = await handler(logItem);
            if (shouldContinue === false) return this;
        }

        // 2. 加入队列
        this.queue.push(logItem);

        // 3. 检查是否达到数量阈值
        if (this.queue.length >= this.batchSize) {
            this.flush();
        }

        return this;
    }

    /**
     * 强制发送队列中的数据 (Flush)
     */
    public flush() {
        if (this.queue.length === 0) return;

        // 取出当前所有数据，并清空队列
        const dataToSend = [...this.queue];
        this.queue = [];

        // 发送数据
        this._sendData(dataToSend);

        // 4. 执行 After Handlers (针对这一批数据)
        for (const handler of this.afterHandlers) {
            handler(dataToSend);
        }
    }

    // 发送实现
    private _sendData(dataList: LogItem[]) {
        // 后端要求 data 字段为字符串，因此需要手动序列化
        const payload = dataList.map(item => ({
            ...item,
            data: typeof item.data === 'object' ? JSON.stringify(item.data) : String(item.data)
        }));

        const body = JSON.stringify({ logs: payload }); // 包装成对象

        // 使用 sendBeacon 保证页面关闭时也能发送
        // 注意：使用 type: 'text/plain' 避免触发 CORS 预检 (OPTIONS 请求)
        if (navigator.sendBeacon) {
            const blob = new Blob([body], { type: 'text/plain; charset=UTF-8' });
            navigator.sendBeacon(this.baseUrl, blob);
        } else {
            fetch(this.baseUrl, {
                method: 'POST',
                body: body,
                headers: {
                    'Content-Type': 'text/plain',
                },
                keepalive: true,
            }).catch(err => console.error('[Tracker] Batch send failed', err));
        }
    }

    // --- 链式配置 (保持不变) ---

    public setUserId(id: number): this {
        this.userId = id;
        return this;
    }

    public before(handler: BeforeHandler): this {
        this.beforeHandlers.push(handler);
        return this;
    }

    public after(handler: AfterHandler): this {
        this.afterHandlers.push(handler);
        return this;
    }

    // --- 场景集成 ---

    public installRouter(router: Router): this {
        router.beforeEach((_: RouteLocationNormalized, from: RouteLocationNormalized, next: any) => {
            // 记录上一个页面的停留时间
            if (from.name) {
                const durationMs = Date.now() - this.pageStartTime;
                let durationStr = '';
                if (durationMs < 60000) {
                    durationStr = `${Math.round(durationMs / 1000)}s`;
                } else {
                    durationStr = `${(durationMs / 60000).toFixed(2)}min`;
                }
                this.report('stay', { path: from.fullPath, duration: durationStr });
            }
            this.pageStartTime = Date.now();
            next();
        });

        router.afterEach((to: RouteLocationNormalized, from: RouteLocationNormalized) => {
            this.report('pv', { path: to.fullPath, referrer: from.fullPath });
        });
        return this;
    }

    public installApi(axiosInstance: AxiosInstance): this {
        axiosInstance.interceptors.request.use((config: InternalAxiosRequestConfig) => {
            config.metadata = { startTime: Date.now() };
            return config;
        });
        axiosInstance.interceptors.response.use(
            (res: AxiosResponse) => {
                const config = res.config as InternalAxiosRequestConfig;
                this.report('api', {
                    url: config.url,
                    method: config.method,
                    status: res.status,
                    duration: Date.now() - (config.metadata?.startTime || Date.now()),
                    success: true
                });
                return res;
            },
            (err: any) => {
                const config = err.config as InternalAxiosRequestConfig | undefined;
                this.report('api', {
                    url: config?.url,
                    status: err.response?.status || 0,
                    duration: Date.now() - (config?.metadata?.startTime || Date.now()),
                    success: false
                });
                return Promise.reject(err);
            }
        );
        return this;
    }

    public installDirective(app: App): this {
        app.directive('tracker', {
            mounted: (el: HTMLElement & { handler?: EventListener }, binding: DirectiveBinding) => {
                const eventType = binding.arg || 'click';
                el.handler = () => this.report('button', { eventType, ...(binding.value || {}) });
                el.addEventListener(eventType, el.handler, true);
            },
            unmounted: (el: HTMLElement & { handler?: EventListener }, binding: DirectiveBinding) => {
                if (el.handler) {
                    el.removeEventListener(binding.arg || 'click', el.handler, true);
                }
            }
        });
        return this;
    }
}

export const tracker = new AnalyticsTracker({
    appId: 'litechat',
    baseUrl: 'http://aihelper.chat:6039/tracker/batch', // 优先使用环境变量配置
    batchSize: 5,  // 积累5条发送一次
    wait: 3000     // 或者每3秒发送一次
});