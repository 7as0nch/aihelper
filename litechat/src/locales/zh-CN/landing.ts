const landing = {
  nav: {
    features: '能力',
    workflow: '架构',
    projects: '实战案例',
    openSource: '开源',
    docs: '文档',
    launch: '进入 LiteChat',
  },
  language: {
    label: '语言',
    zh: '中文',
    en: 'EN',
  },
  hero: {
    badge: 'LiteChat Frontline Build',
    titleLead: 'AI 助手',
    titleHighlight: '产品化前台',
    description: 'LiteChat 把官网、聊天入口、可嵌入组件和后台配置串成一套可复用的前端基础层，帮助团队更快交付自己的 AI 助手，而不是每次都从头搭一遍页面。',
    primary: '启动 LiteChat',
    secondary: '查看 GitHub',
    eyebrow: 'Open-source frontend layer for assistants',
    metrics: [
      {
        label: '编排引擎',
        value: 'Go + Eino',
      },
      {
        label: '后台治理',
        value: 'BackWeb / Vben',
      },
      {
        label: '交付形态',
        value: 'Web + Widget',
      },
    ],
    highlights: [
      '普通聊天与深度搜索共存',
      '截图、网页搜索、语音输入可直接接入',
      '支持继续扩展到多 Agent 与业务工作流',
    ],
  },
  capabilities: {
    eyebrow: 'Capabilities',
    title: '围绕真实助手场景设计，而不是只做一个展示页。',
    description: '官网、聊天入口、嵌入式交付和后台配置放在同一条产品链路里，让前端不再为每个项目重复建设。',
    cards: [
      {
        title: '普通聊天与深度搜索切换',
        description: '在同一个入口里同时处理轻量问答和更复杂的搜索式推理。',
      },
      {
        title: '深度搜索聊天',
        description: '面向复杂问题的多步检索、资料整合与上下文推理。',
      },
      {
        title: '截图理解',
        description: '直接上传页面截图，让 AI 更快理解界面、报错和视觉信息。',
      },
      {
        title: '网页搜索',
        description: '把在线信息补充到回答里，更适合需要时效性的提问场景。',
      },
      {
        title: '语音输入',
        description: '支持语音转文字，降低输入成本，更适合日常记录与发问。',
      },
      {
        title: '多 Agent 协作扩展',
        description: '后端基于 Go 与 Eino 继续推进协作式能力，为复杂业务预留空间。',
      },
    ],
  },
  workflow: {
    eyebrow: 'Workflow',
    title: '从模型能力到业务交付，保持一条清晰的工程路径。',
    description: 'LiteChat 不只是聊天 UI，而是把模型调用、后台配置、页面交付和嵌入式接入组织成更稳定的前端基础层。',
    panels: [
      {
        title: 'Eino 编排引擎',
        description: '负责模型协同、工作流编排和多 Agent 能力延展。',
      },
      {
        title: 'BackWeb 治理后台',
        description: '承接模型配置和后续治理，让官网与后台形成闭环。',
      },
      {
        title: 'Reusable Delivery Layer',
        description: '同一套思路可以落到官网、聊天页和 Widget 植入场景。',
      },
    ],
    callouts: [
      '统一路由与主题基座',
      '官网与产品页共享状态策略',
      '为后续场景化小项目提供样板',
    ],
  },
  projectLab: {
    eyebrow: 'Project Lab',
    title: '让 LiteChat 不只讲能力，也展示怎么落地。',
    description: '这里会持续放入小项目实战案例，用真实场景说明 LiteChat 适合做可交付、可复用、可迭代的 AI 助手产品。',
    featured: {
      tag: '家庭可定制菜谱',
      title: '今晚吃什么，直接问 AI。',
      description: '把家庭成员口味、冰箱食材、人数、时间和忌口交给 AI，随时生成更贴近家里真实情况的做饭建议。',
      bullets: [
        '家庭成员口味可定制',
        '随时问 AI 今天怎么做',
        '可根据食材、人数、忌口、时间生成菜谱建议',
        '适合作为 LiteChat 小项目实战样板',
      ],
      promptLabel: '今晚我家有西红柿、鸡蛋和豆腐，30 分钟内能做什么？',
      repoLabel: '查看 AICook 仓库',
    },
    backlog: [
      {
        title: '家庭计划台',
        description: '把日程、提醒和琐事整理成一个更轻量的家用助手入口。',
      },
      {
        title: '健康饮食助手',
        description: '围绕热量、营养偏好和长期习惯做更细的建议生成。',
      },
      {
        title: 'Personal Jarvis',
        description: '面向个人工作和生活管理的长期助手场景延展。',
      },
    ],
  },
  openSource: {
    eyebrow: 'Open Source',
    title: '开源，但不是一次性的 Demo。',
    description: 'LiteChat 的价值在于减少 AI 产品落地过程中的重复前端开发。它可以先从官网和聊天入口起步，再沿着 Widget、后台治理和场景项目一路扩展。',
    bullets: [
      '当前面向所有用户免费开源',
      '支持 Web 与插件式接入',
      '适合作为多场景 AI 助手的前端基础层',
    ],
    cta: '查看项目路线',
  },
  footer: {
    description: 'AI 助手前端交付基础层。',
    ecosystem: '产品',
    resources: '资源',
    status: '状态',
    statusLabel: '当前状态',
    statusValue: '持续迭代中',
    officialSite: '官网',
    github: 'GitHub',
    gitee: 'Gitee',
    contact: '联系我',
    docs: '文档',
    builtWith: 'BUILT WITH EINO',
    copyright: 'Built for reusable AI delivery.',
  },
};
export default landing;