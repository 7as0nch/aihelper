/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_APP_LOGO: string
  readonly VITE_BASE_URL: string
  readonly VITE_API_BASE_URL: string
  readonly VITE_OPENAI_API_KEY: string
  readonly VITE_OPENAI_BASE_URL: string
  readonly VITE_OPENAI_MODEL: string
  readonly VITE_AI_TYPE: string
  readonly VITE_ENABLE_QR_LOGIN: string
  readonly VITE_TRACKER_ENABLE: string
  // 更多环境变量...
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}