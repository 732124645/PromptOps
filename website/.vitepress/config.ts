import { defineConfig } from 'vitepress'

// Base path for GitHub Pages project site: https://732124645.github.io/PromptOps/
export default defineConfig({
  base: '/PromptOps/',
  title: 'PromptOps',
  description: 'Open-source runtime platform for AI prompts, agents, and workflows.',
  cleanUrls: true,

  themeConfig: {
    socialLinks: [{ icon: 'github', link: 'https://github.com/732124645/PromptOps' }],
    search: { provider: 'local' },
  },

  locales: {
    root: {
      label: 'English',
      lang: 'en-US',
      themeConfig: {
        nav: [
          { text: 'Guide', link: '/guide/introduction' },
          { text: 'Features', link: '/features/prompts' },
          { text: 'SDK', link: '/sdk/node' },
          { text: 'API', link: '/api' },
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'Guide',
              items: [
                { text: 'Introduction', link: '/guide/introduction' },
                { text: 'Quickstart', link: '/guide/quickstart' },
                { text: 'Deployment', link: '/guide/deployment' },
              ],
            },
          ],
          '/features/': [
            {
              text: 'Features',
              items: [
                { text: 'Prompts', link: '/features/prompts' },
                { text: 'Versioning & Publishing', link: '/features/versioning' },
                { text: 'Hot Reload', link: '/features/hot-reload' },
                { text: 'Playground', link: '/features/playground' },
                { text: 'Agents', link: '/features/agents' },
                { text: 'Workflows', link: '/features/workflows' },
                { text: 'Gray Release', link: '/features/rollout' },
                { text: 'Observability', link: '/features/observability' },
                { text: 'Access Control & Workspaces', link: '/features/access-control' },
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
        footer: {
          message: 'Open-source runtime platform for AI prompts, agents, and workflows.',
          copyright: 'PromptOps',
        },
      },
    },

    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/zh/',
      themeConfig: {
        nav: [
          { text: '指南', link: '/zh/guide/introduction' },
          { text: '功能', link: '/zh/features/prompts' },
          { text: 'SDK', link: '/zh/sdk/node' },
          { text: 'API', link: '/zh/api' },
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '指南',
              items: [
                { text: '介绍', link: '/zh/guide/introduction' },
                { text: '快速上手', link: '/zh/guide/quickstart' },
                { text: '部署', link: '/zh/guide/deployment' },
              ],
            },
          ],
          '/zh/features/': [
            {
              text: '功能',
              items: [
                { text: 'Prompt 管理', link: '/zh/features/prompts' },
                { text: '版本与发布', link: '/zh/features/versioning' },
                { text: '热更新', link: '/zh/features/hot-reload' },
                { text: 'Playground', link: '/zh/features/playground' },
                { text: 'Agent', link: '/zh/features/agents' },
                { text: 'Workflow', link: '/zh/features/workflows' },
                { text: '灰度发布', link: '/zh/features/rollout' },
                { text: '可观测性', link: '/zh/features/observability' },
                { text: '权限与工作区', link: '/zh/features/access-control' },
              ],
            },
          ],
          '/zh/sdk/': [
            {
              text: 'SDK',
              items: [
                { text: 'Node', link: '/zh/sdk/node' },
                { text: 'Python', link: '/zh/sdk/python' },
                { text: 'Java', link: '/zh/sdk/java' },
              ],
            },
          ],
        },
        docFooter: { prev: '上一页', next: '下一页' },
        outline: { label: '本页目录' },
        returnToTopLabel: '回到顶部',
        langMenuLabel: '切换语言',
        footer: {
          message: 'AI Prompt、Agent 与 Workflow 的开源运行时平台。',
          copyright: 'PromptOps',
        },
      },
    },
  },
})
