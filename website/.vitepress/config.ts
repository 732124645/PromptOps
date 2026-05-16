import { defineConfig } from 'vitepress'

// Base path for GitHub Pages project site: https://732124645.github.io/PromptOps/
export default defineConfig({
  base: '/PromptOps/',
  lang: 'zh-CN',
  title: 'PromptOps',
  description: 'Open-source runtime platform for AI prompts, agents, and workflows.',
  cleanUrls: true,
  themeConfig: {
    nav: [
      { text: '指南', link: '/guide/introduction' },
      { text: 'SDK', link: '/sdk/node' },
      { text: 'API', link: '/api' },
    ],
    sidebar: {
      '/guide/': [
        {
          text: '指南',
          items: [
            { text: '介绍', link: '/guide/introduction' },
            { text: '快速上手', link: '/guide/quickstart' },
            { text: '部署', link: '/guide/deployment' },
          ],
        },
      ],
      '/sdk/': [
        {
          text: 'SDK',
          items: [
            { text: 'Node', link: '/sdk/node' },
            { text: 'Python', link: '/sdk/python' },
            { text: 'Java', link: '/sdk/java' },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/732124645/PromptOps' },
    ],
    footer: {
      message: 'Open-source runtime platform for AI prompts, agents, and workflows.',
      copyright: 'PromptOps',
    },
    search: { provider: 'local' },
  },
})
