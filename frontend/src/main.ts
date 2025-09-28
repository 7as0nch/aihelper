import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import pinia from './stores'
import { useUserStore } from './stores/user'

// 导入Element Plus和Vant组件库
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import Vant from 'vant'
import 'vant/lib/index.css'

const app = createApp(App)

// 使用插件
app.use(router)
app.use(pinia)
app.use(ElementPlus)
app.use(Vant)

// 初始化用户信息
const userStore = useUserStore()
userStore.initialize()

app.mount('#app')
