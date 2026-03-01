# NovaBlog 使用指南

欢迎使用 NovaBlog！这是一款面向程序员和极客的轻量级混合架构博客系统。本指南将帮助您快速上手使用 NovaBlog 的各项功能。

---

## 目录

1. [快速开始](#快速开始)
2. [文章管理](#文章管理)
3. [MDX 组件使用](#mdx-组件使用)
4. [React 动效组件](#react-动效组件)
5. [动效 HTML 块](#动效-html-块)
6. [评论系统](#评论系统)
7. [微语功能](#微语功能)
8. [用户注册与登录](#用户注册与登录)
9. [主题定制](#主题定制)
10. [附件管理](#附件管理)
11. [常见问题](#常见问题)

---

## 快速开始

### 环境要求

- Node.js 18+
- Go 1.21+ (仅后端开发需要)

### 启动开发服务器

```bash
# 启动前端开发服务器
npm run dev

# 启动后端 API 服务器
cd server && go run cmd/server/main.go
```

访问 `http://localhost:4321` 即可预览博客。

---

## 文章管理

### 创建新文章

所有文章存放在 `src/content/blog/` 目录下，支持 `.md` 和 `.mdx` 格式。

#### 文章命名规范

推荐使用以下命名格式：
- `my-first-post.md` - 英文，小写，连字符分隔
- `使用指南.md` - 中文也可以，但 URL 会进行编码

#### Frontmatter 字段说明

每篇文章需要在开头添加 Frontmatter 元数据：

```yaml
---
title: 文章标题
description: 文章描述，用于 SEO 和社交分享
pubDate: 2024-01-15
updatedDate: 2024-01-20  # 可选，更新日期
heroImage: /images/cover.jpg  # 可选，封面图
heroAlt: 封面图描述  # 可选
author: 作者名  # 可选，默认为站点名称
tags:  # 可选，标签列表
  - JavaScript
  - 教程
category: 技术  # 可选，分类
---

文章内容从这里开始...
```

#### 示例文章

```markdown
---
title: 我的第一篇博客
description: 这是我使用 NovaBlog 发布的第一篇文章
pubDate: 2024-01-15
heroImage: /images/hello-world.jpg
tags:
  - 随笔
  - 博客
category: 生活
---

## 欢迎来到我的博客

这是正文内容。支持标准 Markdown 语法。

### 列表

- 项目 1
- 项目 2
- 项目 3

### 代码块

```javascript
console.log('Hello, NovaBlog!');
```

### 引用

> 这是一段引用文字。
```

### 文章状态

目前文章发布即公开，后续版本将支持：
- `draft: true` - 草稿状态
- `published: false` - 隐藏文章

### 删除文章

直接删除 `src/content/blog/` 目录下的对应文件即可。

---

## MDX 组件使用

NovaBlog 原生支持 MDX，允许在 Markdown 中嵌入交互式组件。

### 什么是 MDX？

MDX = Markdown + JSX。它让您可以在 Markdown 中直接使用 React/Vue 组件。

### 内置组件

#### Counter 计数器组件

```mdx
import Counter from '../components/Counter.vue';

<Counter initialCount={0} />
```



### 自定义组件

您可以在 `src/components/` 目录下创建自己的组件：

```vue
<!-- src/components/MyButton.vue -->
<template>
  <button class="my-btn" @click="handleClick">
    <slot />
  </button>
</template>

<script setup>
const handleClick = () => {
  console.log('Button clicked!');
};
</script>

<style scoped>
.my-btn {
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
}
</style>
```

然后在文章中使用：

```mdx
import MyButton from '../components/MyButton.vue';

<MyButton>点击我</MyButton>
```

### 组件交互性

使用 `client:*` 指令控制组件的客户端行为：

| 指令 | 说明 |
|------|------|
| `client:load` | 页面加载时立即激活 |
| `client:visible` | 组件进入视口时激活 |
| `client:idle` | 浏览器空闲时激活 |
| `client:media` | 满足媒体查询时激活 |

示例：
```mdx
<MyButton client:visible>进入视口时激活</MyButton>
```

---

## React 动效组件

NovaBlog 内置了多个 React 动效组件，可以在 MDX 文章中使用，为内容增添生动的视觉效果。

### 内置 React 动效组件

#### AnimatedCard 悬停动画卡片

一个带有悬停效果的卡片组件，当鼠标悬停时会产生上浮和阴影变化效果。

```mdx
import AnimatedCard from '../components/react/AnimatedCard';

<AnimatedCard 
  title="特色功能" 
  description="NovaBlog 提供丰富的动效组件，让你的博客更加生动" 
  color="#3b82f6" 
  client:visible
/>
```

**属性说明**：
- `title`：卡片标题
- `description`：卡片描述
- `color`：卡片颜色（可选，默认为 `#3b82f6`）

#### FlipCard 翻转卡片

一个可以点击翻转的卡片组件，展示两面不同的内容。

```mdx
import FlipCard from '../components/react/FlipCard';

<FlipCard 
  frontTitle="前端技术" 
  frontDescription="Astro + Vue + React"
  backTitle="后端技术" 
  backDescription="Go + Gin + SQLite"
  frontColor="#3b82f6"
  backColor="#10b981"
  client:visible
/>
```

**属性说明**：
- `frontTitle`：正面标题
- `frontDescription`：正面描述
- `backTitle`：背面标题
- `backDescription`：背面描述
- `frontColor`：正面颜色（可选，默认为 `#3b82f6`）
- `backColor`：背面颜色（可选，默认为 `#10b981`）

#### ParticleBackground 粒子背景

一个带有动态粒子效果的背景组件，可以在粒子上显示自定义内容。

```mdx
import ParticleBackground from '../components/react/ParticleBackground';

<ParticleBackground 
  particleCount={100} 
  color="#6366f1" 
  speed={1.5}
  client:visible
>
  <div style={{ textAlign: 'center' }}>
    <h2 style={{ fontSize: '2.5rem', marginBottom: '1rem' }}>✨ 欢迎来到 NovaBlog</h2>
    <p style={{ fontSize: '1.25rem', opacity: 0.9 }}>探索无限可能</p>
  </div>
</ParticleBackground>
```

**属性说明**：
- `particleCount`：粒子数量（可选，默认为 50）
- `color`：粒子颜色（可选，默认为 `#3b82f6`）
- `speed`：粒子移动速度（可选，默认为 1）
- `children`：要显示的内容（可选）

#### TypewriterText 打字机效果

一个模拟打字机效果的文本组件，逐字显示文本内容。

```mdx
import TypewriterText from '../components/react/TypewriterText';

<TypewriterText 
  text="NovaBlog 是一个极简、高效的程序员博客系统" 
  speed={50} 
  loop={true}
  client:visible
/>
```

**属性说明**：
- `text`：要显示的文本
- `speed`：打字速度（可选，默认为 100ms）
- `loop`：是否循环播放（可选，默认为 false）
- `style`：自定义样式（可选）

#### MathFlipCard 数学公式翻转卡片

一个专门用于展示 LaTeX 数学公式的翻转卡片组件。正面显示渲染后的公式，背面显示 LaTeX 源码。

```mdx
import MathFlipCard from '../components/react/MathFlipCard';

<MathFlipCard 
  latex="E = mc^2" 
  client:visible
/>
```

**属性说明**：
- `latex`：LaTeX 公式字符串
- `displayMode`：是否为块级公式（可选，默认为 true）
- `className`：自定义 CSS 类名（可选）

**示例**：

```mdx
---
title: 数学公式展示
description: 使用翻转卡片展示数学公式
pubDate: 2024-01-20
tags: [数学, LaTeX]
---

import MathFlipCard from '../components/react/MathFlipCard';

# 数学公式展示

## 质能方程

<MathFlipCard latex="E = mc^2" client:load />

## 麦克斯韦方程组

<MathFlipCard latex="\nabla \cdot \mathbf{E} = \frac{\rho}{\varepsilon_0}" client:load />

## 薛定谔方程

<MathFlipCard latex="i\hbar\frac{\partial}{\partial t}\Psi(\mathbf{r},t) = \left[-\frac{\hbar^2}{2m}\nabla^2 + V(\mathbf{r},t)\right]\Psi(\mathbf{r},t)" client:load />

## 矩阵运算

<MathFlipCard latex="\begin{pmatrix} a & b \\ c & d \end{pmatrix} \begin{pmatrix} x \\ y \end{pmatrix} = \begin{pmatrix} ax + by \\ cx + dy \end{pmatrix}" client:load />
```

**特点**：
- 点击"显示 LaTeX"按钮可查看公式源码
- 点击"显示公式"按钮返回渲染结果
- 支持 KaTeX 的所有语法
- 适合教学和技术文档

### 在文章中使用

在 MDX 文章中导入并使用这些组件：

```mdx
---
title: 动效组件展示
description: 展示 NovaBlog 中的 React 动效组件
pubDate: 2024-01-20
tags: [React, 动效, 组件]
---

import AnimatedCard from '../components/react/AnimatedCard';
import FlipCard from '../components/react/FlipCard';
import ParticleBackground from '../components/react/ParticleBackground';
import TypewriterText from '../components/react/TypewriterText';

# React 动效组件展示

## 打字机效果

<TypewriterText 
  text="这是一个打字机效果的示例" 
  client:visible
/>

## 悬停卡片

<AnimatedCard 
  title="悬停效果" 
  description="鼠标悬停时会产生动画效果" 
  color="#8b5cf6" 
  client:visible
/>

## 翻转卡片

<FlipCard 
  frontTitle="点击我" 
  frontDescription="查看背面内容" 
  backTitle="翻转效果" 
  backDescription="这是卡片的背面" 
  client:visible
/>

## 粒子背景

<ParticleBackground client:visible>
  <h3>粒子背景效果</h3>
  <p>带有动态粒子的背景</p>
</ParticleBackground>
```

### 自定义样式

这些组件都支持通过 `style` 属性自定义样式，例如：

```mdx
<TypewriterText 
  text="自定义样式示例" 
  style={{
    fontSize: '2rem',
    color: '#ef4444',
    fontWeight: 'bold'
  }}
  client:visible
/>
```

---

## 动效 HTML 块

在 MDX 中可以直接使用 HTML 标签和内联样式，创建带动画的内容块。

### 示例：渐变背景卡片

```mdx
<div style={{
  padding: '2rem',
  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  borderRadius: '1rem',
  color: 'white',
  textAlign: 'center'
}}>
  <h3 style={{ marginBottom: '0.5rem' }}>✨ 特色卡片</h3>
  <p>这是一个带渐变背景的卡片</p>
</div>
```

### 示例：CSS 动画

```mdx
<style>
  .pulse-box {
    animation: pulse 2s infinite;
  }
  
  @keyframes pulse {
    0% { transform: scale(1); opacity: 1; }
    50% { transform: scale(1.05); opacity: 0.8; }
    100% { transform: scale(1); opacity: 1; }
  }
</style>

<div class="pulse-box" style={{
  padding: '1.5rem',
  background: '#3b82f6',
  borderRadius: '0.5rem',
  color: 'white',
  textAlign: 'center'
}}>
  🔔 脉冲动画效果
</div>
```

### 示例：悬停效果

```mdx
<style>
  .hover-card {
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  .hover-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  }
</style>

<div class="hover-card" style={{
  padding: '1.5rem',
  background: 'white',
  borderRadius: '0.5rem',
  boxShadow: '0 2px 10px rgba(0, 0, 0, 0.1)',
  cursor: 'pointer'
}}>
  🖱️ 鼠标悬停试试
</div>
```

### 示例：交互式组件

结合 Vue 组件创建交互式内容：

```vue
<!-- src/components/InteractiveCard.vue -->
<template>
  <div 
    class="card"
    @click="toggle"
    :style="{ background: isActive ? '#3b82f6' : '#e5e7eb' }"
  >
    <h3>{{ title }}</h3>
    <p>{{ isActive ? '已激活 ✓' : '点击激活' }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  title: String
});

const isActive = ref(false);
const toggle = () => isActive.value = !isActive.value;
</script>

<style scoped>
.card {
  padding: 1.5rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  color: white;
}
</style>
```

---

## 评论系统

NovaBlog 内置评论系统，支持多级嵌套回复和 Markdown 语法。

### 发表评论

1. 在文章页面滚动到评论区
2. 点击登录按钮进行用户认证
3. 输入评论内容（支持 Markdown）
4. 点击发送

### 回复评论

1. 点击评论下方的"回复"按钮
2. 输入回复内容
3. 提交后将显示在原评论下方

### Markdown 支持

评论区支持基础 Markdown 语法：

| 语法 | 效果 |
|------|------|
| `**粗体**` | **粗体** |
| `*斜体*` | *斜体* |
| `` `代码` `` | `代码` |
| `[链接](url)` | [链接](url) |
| `> 引用` | 引用块 |

### 删除评论

用户可以删除自己发表的评论。管理员可以删除任何评论。

---

## 微语功能

微语是一个轻量级的分享空间，让你可以随时记录生活中的点滴、灵感与感悟。类似于社交媒体的动态功能，所有注册用户都可以发布。

### 访问微语

点击导航栏中的"微语"链接，或访问 `/micro` 页面。

### 发布微语

1. 登录你的账号
2. 在微语页面顶部的输入框中输入内容
3. 可选：添加标签（用逗号或空格分隔）
4. 选择是否公开可见
5. 点击"发布"按钮

**内容限制**：
- 单条微语最多 2000 字
- 支持多标签

### 微语列表

- 所有公开的微语都会显示在列表中
- 支持分页加载更多
- 显示发布者头像、昵称和发布时间
- 支持点赞功能

### 热力图

微语页面右侧显示 GitHub 风格的热力图，展示一年内的发布活动：

- **颜色深浅**：表示当天发布的微语数量
- **悬停查看**：鼠标悬停可查看具体日期和数量
- **年度统计**：显示全年发布的微语总数

### 点赞微语

1. 登录后可以给微语点赞
2. 点击心形图标即可点赞或取消点赞
3. 点赞数会实时更新

### 删除微语

- 用户可以删除自己发布的微语
- 管理员可以删除任何微语

### 微语与博客的区别

| 特性 | 博客文章 | 微语 |
|------|----------|------|
| 内容长度 | 无限制 | 最多 2000 字 |
| 格式支持 | Markdown + MDX | 纯文本 |
| 发布权限 | 管理员 | 所有注册用户 |
| 适用场景 | 长篇教程、技术文章 | 随手记录、灵感分享 |
| 互动功能 | 评论 | 点赞 |

---

## 用户注册与登录

### 注册账号

1. 点击页面右上角的"登录"按钮
2. 在登录弹窗中点击"注册"
3. 填写用户名、邮箱和密码
4. 提交注册

### 登录账号

支持两种登录方式：
- 用户名 + 密码
- 邮箱 + 密码

### 用户角色

| 角色 | 权限 |
|------|------|
| `user` | 发表评论、删除自己的评论 |
| `admin` | 所有权限 + 删除任意评论 |

### 修改个人资料

登录后可以修改：
- 昵称
- 头像 URL
- 个人简介

---

## 主题定制

### CSS 变量

NovaBlog 使用 CSS 变量实现主题系统，可在 `src/styles/global.css` 中修改：

```css
:root {
  /* 主色调 */
  --color-primary-50: #eff6ff;
  --color-primary-100: #dbeafe;
  --color-primary-500: #3b82f6;
  --color-primary-600: #2563eb;
  
  /* 背景色 */
  --color-background: #ffffff;
  --color-muted: #f3f4f6;
  
  /* 文字色 */
  --color-foreground: #1f2937;
  --color-muted-foreground: #6b7280;
  
  /* 边框色 */
  --color-border: #e5e7eb;
}
```

### 暗色主题

修改 CSS 变量实现暗色模式：

```css
.dark {
  --color-background: #1a1a2e;
  --color-foreground: #eaeaea;
  --color-muted: #16213e;
  --color-border: #0f3460;
}
```

### Tailwind 配置

在 `tailwind.config.mjs` 中自定义主题：

```javascript
export default {
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          // ... 更多颜色
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      }
    }
  }
}
```

### 布局组件

主要布局文件：
- `src/layouts/BaseLayout.astro` - 基础布局（Header/Footer）
- `src/layouts/PostLayout.astro` - 文章布局（目录/评论）

---

## 附件管理

### 图片存放位置

将图片放在 `public/images/` 目录下：

```
public/
├── images/
│   ├── posts/
│   │   ├── hello-world.jpg
│   │   └── tutorial-cover.png
│   └── avatars/
│       └── default.png
```

### 在文章中引用

```markdown
![图片描述](/images/posts/hello-world.jpg)
```

### 在 Frontmatter 中使用封面图

```yaml
---
heroImage: /images/posts/my-cover.jpg
heroAlt: 这是一张封面图
---
```

### 其他附件

PDF、ZIP 等文件也放在 `public/` 目录：

```
public/
├── downloads/
│   └── source-code.zip
└── documents/
    └── paper.pdf
```

引用方式：
```markdown
[下载源码](/downloads/source-code.zip)
[查看论文](/documents/paper.pdf)
```

### 图片优化建议

- 使用 WebP 格式减少文件大小
- 封面图推荐尺寸：1920x1080
- 文章内图片推荐宽度：800px
- 压缩工具：[Squoosh](https://squoosh.app/)

---

## 常见问题

### Q: 文章修改后没有更新？

A: 开发模式下 Astro 会自动热重载。如果生产构建，需要重新运行 `npm run build`。

### Q: 评论无法发送？

A: 检查后端服务是否正常运行在 `localhost:8080`。

### Q: 如何修改站点信息？

A: 编辑 `src/layouts/BaseLayout.astro` 和 `astro.config.mjs` 中的站点配置。

### Q: 如何添加 Google Analytics？

A: 在 `src/layouts/BaseLayout.astro` 的 `<head>` 中添加 GA 脚本。

---

## 获取帮助

- 📖 查看 [开发文档](./developer-guide.md) 了解技术细节
- 🐛 提交 Issue 反馈问题
- 💬 在 GitHub Discussions 中讨论