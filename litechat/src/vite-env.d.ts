/// <reference types="vite/client" />


interface ImportMetaEnv {
    readonly VITE_TRACKER_ENABLE: string
    readonly VITE_APP_TITLE: string
    readonly VITE_APP_LOGO: string
    readonly VITE_BASE_URL: string
    readonly VITE_API_BASE_URL: string
    readonly VITE_OPENAI_API_KEY: string
    readonly VITE_OPENAI_BASE_URL: string
    readonly VITE_OPENAI_MODEL: string
    readonly VITE_AI_TYPE: string
    readonly VITE_ENABLE_QR_LOGIN: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}
