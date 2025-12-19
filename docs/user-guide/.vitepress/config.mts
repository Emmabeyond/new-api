import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'BigAI Pro',
  description: '企业级 AI API 聚合平台文档',
  lang: 'zh-CN',
  
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3b82f6' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:site_name', content: 'BigAI Pro Docs' }],
  ],

  themeConfig: {
    logo: '/logo.jpg',
    siteTitle: 'BigAI Pro',
    
    nav: [
      { text: '首页', link: '/' },
      { text: '快速入门', link: '/quick-start' },
      { text: '模型', link: '/models/overview' },
      { text: '开发工具', link: '/tools/claude-code' },
      { 
        text: '控制台', 
        link: 'https://api.bigaipro.com/console',
        target: '_blank'
      }
    ],

    sidebar: {
      '/': [
        {
          text: '🚀 快速入门',
          items: [
            { text: '快速入门指南', link: '/quick-start' },
            { text: 'API Key 管理', link: '/api-key-management' },
            { text: '计费与额度', link: '/billing-and-quota' },
          ]
        },
        {
          text: '🤖 模型使用',
          collapsed: false,
          items: [
            { text: '模型总览', link: '/models/overview' },
            { text: 'GPT 系列', link: '/models/gpt-models' },
            { text: 'Claude 系列', link: '/models/claude-models' },
            { text: 'Gemini 系列', link: '/models/gemini-models' },
            { text: '国产模型', link: '/models/chinese-models' },
          ]
        },
        {
          text: '🛠️ 开发工具',
          collapsed: false,
          items: [
            { text: 'Claude Code', link: '/tools/claude-code' },
            { text: 'Codex CLI', link: '/tools/codex-cli' },
            { text: 'Cursor IDE', link: '/tools/cursor-ide' },
            { text: 'Continue 插件', link: '/tools/continue-plugin' },
            { text: 'Kiro IDE', link: '/tools/kiro-ide' },
            { text: 'Cherry Studio', link: '/tools/cherry-studio' },
          ]
        },
        {
          text: '🎨 高级功能',
          collapsed: true,
          items: [
            { text: '多模态使用', link: '/advanced/multimodal' },
            { text: '流式响应', link: '/advanced/streaming' },
            { text: '函数调用', link: '/advanced/function-calling' },
            { text: '错误处理', link: '/advanced/error-handling' },
          ]
        },
        {
          text: '📖 最佳实践',
          collapsed: true,
          items: [
            { text: '性能优化', link: '/best-practices/performance' },
            { text: '安全指南', link: '/best-practices/security' },
            { text: '成本控制', link: '/best-practices/cost-control' },
          ]
        }
      ]
    },

    footer: {
      message: 'BigAI Pro - 让 AI 触手可及',
      copyright: 'Copyright © 2025 BigAI Pro'
    },

    search: {
      provider: 'local',
      options: {
        translations: {
          button: {
            buttonText: '搜索文档',
            buttonAriaLabel: '搜索文档'
          },
          modal: {
            noResultsText: '无法找到相关结果',
            resetButtonTitle: '清除查询条件',
            footer: {
              selectText: '选择',
              navigateText: '切换',
              closeText: '关闭'
            }
          }
        }
      }
    },

    outline: {
      label: '页面导航',
      level: [2, 3]
    },

    docFooter: {
      prev: '上一篇',
      next: '下一篇'
    },

    lastUpdated: {
      text: '最后更新于',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'short'
      }
    },

    editLink: {
      pattern: 'https://github.com/bigaipro/docs/edit/main/:path',
      text: '编辑此页'
    },

    returnToTopLabel: '回到顶部',
    sidebarMenuLabel: '菜单',
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式'
  },

  markdown: {
    lineNumbers: true,
    theme: {
      light: 'github-light',
      dark: 'github-dark'
    }
  },

  vite: {
    server: {
      port: 5173
    }
  }
})
