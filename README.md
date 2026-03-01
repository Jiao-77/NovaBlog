# NovaBlog

🚀 一款面向程序员和极客的极轻量级博客系统，采用"静态渲染 + 轻量级微服务"的解耦架构。

## ✨ 特性

- **极致的排版自由**：原生支持 MDX，允许在 Markdown 中直接嵌入 Vue/React 组件和 Typst 复杂学术排版
- **优雅的动态交互**：支持用户注册登录、评论、点赞，但核心页面保持纯静态
- **低资源占用**：Go 后端 + SQLite，内存占用极低（通常十几 MB），可在 2C1G 甚至更低配置运行
- **现代化主题系统**：基于组件化思想构建，CSS 变量换肤，支持暗黑模式
- **Docker 容器化**：一键部署，方便在个人 NAS 或云服务器上运行

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────┐
│                      Nginx                              │
│  (静态文件服务 + API 反向代理)                           │
└─────────────────────┬───────────────────────────────────┘
                      │
        ┌─────────────┴─────────────┐
        │                           │
        ▼                           ▼
┌───────────────┐           ┌───────────────┐
│   Astro 前端   │           │   Go API     │
│  (静态 HTML)   │◄─────────►│  (Gin + SQLite)│
│   Zero-JS     │   HTTP    │   ~15MB 内存  │
└───────────────┘           └───────────────┘
```

## 📁 项目结构

```
NovaBlog/
├── src/                    # Astro 前端源码
│   ├── components/         # Vue/Astro 组件
│   │   ├── CommentSection.vue  # 评论区组件
│   │   ├── LikeButton.vue      # 点赞按钮
│   │   ├── Counter.vue         # MDX 示例组件
│   │   └── TypstBlock.astro    # Typst 渲染组件
│   ├── content/            # 内容集合
│   │   ├── blog/           # 博客文章 (MDX)
│   │   └── config.ts       # 内容配置
│   ├── layouts/            # 页面布局
│   │   ├── BaseLayout.astro
│   │   └── PostLayout.astro
│   ├── pages/              # 页面路由
│   └── styles/             # 全局样式
├── server/                 # Go 后端源码
│   ├── cmd/server/         # 程序入口
│   ├── internal/
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库连接
│   │   ├── handlers/       # HTTP 处理器
│   │   ├── middleware/     # 中间件
│   │   ├── models/         # 数据模型
│   │   └── utils/          # 工具函数
│   └── Dockerfile
├── docker-compose.yml      # Docker 编排
└── nginx.conf              # Nginx 配置
```

## 🚀 快速开始

### 开发环境

**前端开发：**

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:4321
```

**后端开发：**

```bash
cd server

# 下载依赖
go mod download

# 运行服务
go run ./cmd/server

# API 服务于 http://localhost:8080
```

### 生产部署

```bash
# 1. 构建前端
npm run build

# 2. 启动 Docker 容器
docker-compose up -d

# 访问 http://localhost
```

## 📝 写作

在 `src/content/blog/` 目录下创建 `.mdx` 或 `.md` 文件：

```yaml
---
title: 我的第一篇文章
description: 文章描述
pubDate: 2024-01-01
author: NovaBlog
tags: [博客, 技术]
category: 技术
heroImage: /images/hero.jpg
---

# Hello NovaBlog!

这是一篇使用 MDX 编写的文章...

<!-- 嵌入 Vue 组件 -->
<Counter client:load />

<!-- Typst 数学公式 -->
<TypstBlock code="$ sum_(i=1)^n x_i = x_1 + x_2 + dots + x_n $" />
```

## 🔌 API 接口

### 认证

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/auth/register` | 用户注册 |
| POST | `/api/auth/login` | 用户登录 |
| GET | `/api/auth/profile` | 获取当前用户信息 (需认证) |

### 评论

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/comments?post_id=xxx` | 获取文章评论 |
| POST | `/api/comments` | 发表评论 (需认证) |
| DELETE | `/api/comments/:id` | 删除评论 (需认证) |

### 点赞

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/likes?post_id=xxx` | 获取点赞状态 |
| POST | `/api/likes` | 切换点赞状态 (需认证) |

## 🛠️ 技术栈

**前端：**
- [Astro](https://astro.build/) - 静态站点生成
- [Vue 3](https://vuejs.org/) - 交互组件
- [Tailwind CSS](https://tailwindcss.com/) - 样式

**后端：**
- [Go](https://golang.org/) - 编程语言
- [Gin](https://gin-gonic.com/) - Web 框架
- [GORM](https://gorm.io/) - ORM
- [SQLite](https://www.sqlite.org/) - 数据库
- [JWT](https://jwt.io/) - 认证

**部署：**
- [Docker](https://www.docker.com/) - 容器化
- [Nginx](https://nginx.org/) - 反向代理

## 📚 文档

- **[使用指南](./docs/user-guide.md)** - 博客使用教程，包括文章编写、MDX 组件、Typst 排版、主题定制等
- **[开发文档](./docs/developer-guide.md)** - 面向开发者的技术文档，包括 API 接口、数据库结构、深度定制指南等

## 🤖 AI 使用声明

本系统大部分代码由 AI 辅助生成，但所有代码都经过了人工审核和测试。

### 为什么使用 AI 辅助开发？

AI 辅助开发能够显著提高开发效率，让我更专注于系统设计和用户体验。同时，我也深知大家对 AI 生成代码的担忧：

### 常见担忧与应对

| 担忧 | 应对措施 |
|------|----------|
| **隐私泄露风险** | 系统完全本地部署，不会向任何第三方服务器发送数据。所有用户数据（评论、点赞等）仅存储在您自己的 SQLite 数据库中 |
| **代码安全性** | 所有 AI 生成的代码都经过人工逐行审查，重点检查：SQL 注入、XSS 攻击、敏感信息泄露等安全漏洞 |
| **代码质量** | 对 AI 生成的代码进行了重构和优化，确保符合最佳实践 |
| **后门/恶意代码** | 代码完全开源，您可以自行审计每一行代码。没有任何隐藏的网络请求或数据收集行为 |
| **依赖安全** | 所有依赖包都经过筛选，使用知名、活跃维护的开源项目，避免引入有安全隐患的第三方库 |

### 我的承诺

- ✅ **隐私优先**：您的数据属于您自己，我无权也无从访问
- ✅ **透明公开**：坦诚披露 AI 的参与，不隐瞒开发方式
- ✅ **人工把关**：AI 只是辅助工具，最终的决策和审核由人工完成
- ✅ **开源审计**：代码完全开源，欢迎社区审查和改进

如果您对代码有任何疑虑，欢迎：
- 查看完整的 [开发文档](./docs/developer-guide.md) 了解技术细节
- 直接阅读源代码进行审计
- 提交 Issue 反馈问题或建议

## 📜 License

MIT License © 2024
