import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "RuoYi Web Documentation",
  description: "Documentation for RuoYi Web Application",
  base: "/",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/' }
    ],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Introduction', link: '/guide/' },
          { text: 'Getting Started', link: '/guide/getting-started' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/your-repo/ruoyi-web' }
    ]
  }
})