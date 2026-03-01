# NovaBlog 开发者文档

本文档面向希望深度定制 NovaBlog 的开发者，详细介绍系统架构、API 接口、数据库结构以及扩展开发指南。

---

## 目录

1. [系统架构](#系统架构)
2. [项目结构](#项目结构)
3. [API 接口文档](#api-接口文档)
4. [数据库结构](#数据库结构)
5. [前端组件开发](#前端组件开发)
6. [后端服务扩展](#后端服务扩展)
7. [部署配置](#部署配置)
8. [性能优化](#性能优化)

---

## 系统架构

NovaBlog 采用 **静态渲染 + 轻量级微服务** 的解耦架构：

```
┌─────────────────────────────────────────────────────────────┐
│                         用户浏览器                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Nginx 反向代理                           │
│            (静态文件服务 + API 请求转发)                      │
└─────────────────────────────────────────────────────────────┘
         │                                    │
         ▼                                    ▼
┌─────────────────────┐          ┌──────────────────────────┐
│   静态文件 (CDN)     │          │    API Server (Go)       │
│  HTML/CSS/JS/图片   │          │  - 用户认证              │
│                     │          │  - 评论管理              │
│  构建时生成         │          │  - 点赞系统              │
└─────────────────────┘          │  - SQLite 数据库         │
                                 └──────────────────────────┘
```

### 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 前端框架 | Astro 5.x | Islands Architecture，零 JS 默认输出 |
| UI 组件 | Vue 3 | 交互式岛屿组件 |
| 样式方案 | Tailwind CSS 4.x | 原子化 CSS |
| 后端框架 | Go + Gin | 极致轻量，内存占用 < 20MB |
| 数据库 | SQLite | 文件型数据库，无需额外服务 |
| 认证方案 | JWT | 无状态认证 |

### 架构优势

1. **极致性能**：静态页面零运行时，JS 按需加载
2. **低资源占用**：Go 服务 + SQLite 可在 512MB 内存环境运行
3. **SEO 友好**：纯静态 HTML 输出，搜索引擎完美抓取
4. **开发体验**：组件化开发，热重载支持

---

## 项目结构

```
NovaBlog/
├── src/                        # 前端源码
│   ├── components/             # Vue/Astro 组件
│   │   ├── CommentSection.vue  # 评论区组件
│   │   ├── LikeButton.vue      # 点赞按钮
│   │   ├── LoginForm.vue       # 登录表单
│   │   ├── UserStatus.vue      # 用户状态栏
│   │   ├── Counter.vue         # 计数器示例
│   │   ├── TableOfContents.astro # 目录组件
│   │   └── react/              # React 组件
│   │       ├── AnimatedCard.tsx
│   │       ├── FlipCard.tsx
│   │       ├── ParticleBackground.tsx
│   │       ├── TypewriterText.tsx
│   │       ├── Heatmap.tsx     # 热力图组件
│   │       ├── MicroList.tsx   # 微语列表
│   │       ├── MicroComposer.tsx # 发布微语
│   │       └── MicroPage.tsx   # 微语页面容器
│   ├── content/                # 内容集合
│   │   ├── config.ts           # 内容配置
│   │   └── blog/               # 博客文章
│   ├── layouts/                # 布局组件
│   │   ├── BaseLayout.astro    # 基础布局
│   │   └── PostLayout.astro    # 文章布局
│   ├── pages/                  # 页面路由
│   │   ├── index.astro         # 首页
│   │   ├── login.astro         # 登录页
│   │   ├── micro.astro         # 微语页
│   │   ├── blog/               # 博客相关页面
│   │   ├── tags/               # 标签页面
│   │   └── categories/         # 分类页面
│   ├── styles/                 # 全局样式
│   │   └── global.css          # CSS 变量和全局样式
│   ├── env.d.ts                # 类型声明
│   └── mdx-components.ts       # MDX 组件注册
├── public/                     # 静态资源
│   ├── images/                 # 图片资源
│   └── favicon.svg             # 网站图标
├── server/                     # 后端服务
│   ├── cmd/server/main.go      # 服务入口
│   ├── internal/               # 内部模块
│   │   ├── config/             # 配置管理
│   │   ├── database/           # 数据库连接
│   │   ├── handlers/           # HTTP 处理器
│   │   │   ├── auth.go         # 认证处理
│   │   │   ├── comment.go      # 评论处理
│   │   │   ├── like.go         # 点赞处理
│   │   │   └── micro.go        # 微语处理
│   │   ├── middleware/         # 中间件
│   │   ├── models/             # 数据模型
│   │   └── utils/              # 工具函数
│   ├── data/                   # 数据文件 (SQLite)
│   ├── migrations/             # 数据库迁移
│   ├── Dockerfile              # Docker 构建文件
│   └── go.mod                  # Go 依赖
├── docs/                       # 文档
├── astro.config.mjs            # Astro 配置
├── tailwind.config.mjs         # Tailwind 配置
├── docker-compose.yml          # Docker Compose 配置
├── nginx.conf                  # Nginx 配置
└── package.json                # NPM 依赖
```

---

## API 接口文档

### 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token
- **内容格式**: JSON

### 认证接口

#### 注册用户

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "string",    // 必填，3-50字符
  "email": "string",       // 必填，有效邮箱
  "password": "string",    // 必填，6-50字符
  "nickname": "string"     // 可选
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "role": "user",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

**错误码**:
- `400`: 请求参数无效
- `409`: 用户名或邮箱已存在

#### 用户登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "string",    // 用户名或邮箱
  "password": "string"
}
```

**响应**: 同注册接口

#### 获取当前用户

```http
GET /api/auth/me
Authorization: Bearer <token>
```

**响应**:
```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "https://...",
    "bio": "个人简介",
    "role": "user"
  }
}
```

#### 更新用户资料

```http
PUT /api/auth/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "string",
  "avatar": "string",
  "bio": "string"
}
```

### 评论接口

#### 获取文章评论

```http
GET /api/comments?post_id={post_id}&page=1&page_size=20
```

**参数**:
- `post_id` (必填): 文章 ID
- `page` (可选): 页码，默认 1
- `page_size` (可选): 每页数量，默认 20，最大 100

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "post_id": "hello-novablog",
      "user_id": 1,
      "content": "很棒的文章！",
      "status": "approved",
      "created_at": "2024-01-15T10:00:00Z",
      "user": {
        "id": 1,
        "username": "testuser",
        "nickname": "测试用户",
        "avatar": "https://..."
      },
      "replies": [
        {
          "id": 2,
          "parent_id": 1,
          "content": "感谢支持！",
          "user": {...}
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_page": 5
  }
}
```

#### 创建评论

```http
POST /api/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "post_id": "string",     // 必填
  "content": "string",     // 必填，1-2000字符
  "parent_id": 1           // 可选，回复的评论ID
}
```

#### 删除评论

```http
DELETE /api/comments/:id
Authorization: Bearer <token>
```

**权限**: 本人或管理员可删除

### 点赞接口

#### 切换点赞状态

```http
POST /api/likes/toggle
Content-Type: application/json
Authorization: Bearer <token>  // 可选

{
  "post_id": "string"      // 必填
}
```

**响应**:
```json
{
  "liked": true,           // 当前是否已点赞
  "like_count": 42         // 文章总点赞数
}
```

**说明**:
- 已登录用户：基于 user_id 判断
- 未登录用户：基于 IP Hash 判断（加盐防反向推导）

#### 获取点赞状态

```http
GET /api/likes/status?post_id={post_id}
Authorization: Bearer <token>  // 可选
```

**响应**: 同切换接口

### 微语接口

#### 获取微语列表

```http
GET /api/micros?page=1&page_size=20&user_id=1
Authorization: Bearer <token>  // 可选
```

**参数**:
- `page` (可选): 页码，默认 1
- `page_size` (可选): 每页数量，默认 20，最大 50
- `user_id` (可选): 指定用户的微语

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "content": "今天天气真好！",
      "images": "[]",
      "tags": "[\"生活\", \"日常\"]",
      "is_public": true,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z",
      "user": {
        "id": 1,
        "username": "testuser",
        "nickname": "测试用户",
        "avatar": "https://..."
      },
      "like_count": 5,
      "is_liked": false
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_page": 5
  }
}
```

#### 获取单条微语

```http
GET /api/micros/:id
Authorization: Bearer <token>  // 可选
```

#### 发布微语

```http
POST /api/micros
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "string",     // 必填，最多 2000 字
  "images": ["url1", "url2"],  // 可选，图片 URL 数组
  "tags": ["tag1", "tag2"],    // 可选，标签数组
  "is_public": true        // 可选，默认 true
}
```

**响应**:
```json
{
  "id": 1,
  "content": "今天天气真好！",
  "images": "[]",
  "tags": "[\"生活\"]",
  "is_public": true,
  "created_at": "2024-01-15T10:00:00Z",
  "user": {...},
  "like_count": 0,
  "is_liked": false
}
```

#### 更新微语

```http
PUT /api/micros/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "string",
  "images": ["url1"],
  "tags": ["tag1"],
  "is_public": true
}
```

**权限**: 仅作者可修改

#### 删除微语

```http
DELETE /api/micros/:id
Authorization: Bearer <token>
```

**权限**: 作者或管理员可删除

#### 点赞/取消点赞微语

```http
POST /api/micros/:id/like
Authorization: Bearer <token>
```

**响应**:
```json
{
  "liked": true,
  "message": "点赞成功"
}
```

#### 获取热力图数据

```http
GET /api/micros/heatmap?year=2024&user_id=1
```

**参数**:
- `year` (可选): 年份，默认当前年
- `user_id` (可选): 指定用户

**响应**:
```json
[
  { "date": "2024-01-15", "count": 3 },
  { "date": "2024-01-16", "count": 1 },
  { "date": "2024-01-20", "count": 5 }
]
```

#### 获取统计数据

```http
GET /api/micros/stats?user_id=1
```

**响应**:
```json
{
  "total_micros": 150,
  "total_users": 25,
  "top_users": [
    {
      "user_id": 1,
      "username": "admin",
      "nickname": "管理员",
      "avatar": "...",
      "post_count": 45
    }
  ]
}
```

### 错误响应格式

```json
{
  "error": "错误信息描述"
}
```

---

## 数据库结构

### users 表

```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,          -- 软删除
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL, -- bcrypt 哈希
  nickname VARCHAR(50),
  avatar VARCHAR(255),
  role VARCHAR(20) DEFAULT 'user', -- 'admin' | 'user'
  bio VARCHAR(500)
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

### comments 表

```sql
CREATE TABLE comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  post_id VARCHAR(100) NOT NULL,   -- 文章 slug
  user_id INTEGER NOT NULL,        -- 关联 users.id
  parent_id INTEGER,               -- 父评论ID，用于嵌套回复
  content TEXT NOT NULL,
  status VARCHAR(20) DEFAULT 'approved', -- 'pending' | 'approved' | 'spam'
  
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (parent_id) REFERENCES comments(id)
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at);
```

### likes 表

```sql
CREATE TABLE likes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  post_id VARCHAR(100) NOT NULL,
  user_id INTEGER,              -- 登录用户ID，可为空
  ip_hash VARCHAR(64),          -- 访客IP哈希
  
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 防止同一用户重复点赞
CREATE UNIQUE INDEX idx_post_user ON likes(post_id, user_id);
-- 防止同一IP重复点赞
CREATE UNIQUE INDEX idx_post_ip ON likes(post_id, ip_hash);
CREATE INDEX idx_likes_user_id ON likes(user_id);
```

### like_counts 表

```sql
CREATE TABLE like_counts (
  post_id VARCHAR(100) PRIMARY KEY,
  count INTEGER DEFAULT 0
);
```

### post_meta 表（预留）

```sql
CREATE TABLE post_meta (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_id VARCHAR(100) UNIQUE NOT NULL,
  view_count INTEGER DEFAULT 0,
  like_count INTEGER DEFAULT 0,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE INDEX idx_post_meta_post_id ON post_meta(post_id);
```

### micro_posts 表

```sql
CREATE TABLE micro_posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,          -- 软删除
  user_id INTEGER NOT NULL,     -- 关联 users.id
  content TEXT NOT NULL,        -- 微语内容，最多 2000 字
  images TEXT,                  -- JSON 数组存储图片 URL
  tags TEXT,                    -- JSON 数组存储标签
  is_public BOOLEAN DEFAULT 1,  -- 是否公开
  
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_micro_posts_user_id ON micro_posts(user_id);
CREATE INDEX idx_micro_posts_deleted_at ON micro_posts(deleted_at);
CREATE INDEX idx_micro_posts_created_at ON micro_posts(created_at);
```

### micro_post_likes 表

```sql
CREATE TABLE micro_post_likes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  micro_post_id INTEGER NOT NULL,  -- 关联 micro_posts.id
  user_id INTEGER NOT NULL,        -- 关联 users.id
  
  FOREIGN KEY (micro_post_id) REFERENCES micro_posts(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 防止同一用户重复点赞同一条微语
CREATE UNIQUE INDEX idx_micropost_user ON micro_post_likes(micro_post_id, user_id);
CREATE INDEX idx_micropost_likes_user_id ON micro_post_likes(user_id);
```

---

## 前端组件开发

### 创建 Vue 组件

在 `src/components/` 下创建 `.vue` 文件：

```vue
<!-- src/components/MyComponent.vue -->
<template>
  <div class="my-component">
    <h2>{{ title }}</h2>
    <button @click="handleClick">{{ count }}</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// Props 定义
const props = defineProps<{
  title: string;
  initialValue?: number;
}>();

// 响应式状态
const count = ref(props.initialValue || 0);

// 方法
const handleClick = () => {
  count.value++;
  emit('change', count.value);
};

// 事件定义
const emit = defineEmits<{
  change: [value: number];
}>();
</script>

<style scoped>
.my-component {
  padding: 1rem;
  background: var(--color-muted);
  border-radius: 0.5rem;
}
</style>
```

### 在 MDX 中使用

```mdx
import MyComponent from '../components/MyComponent.vue';

<MyComponent 
  title="我的组件" 
  initialValue={10}
  client:visible
/>
```

### 创建 Astro 组件

Astro 组件用于静态渲染，无客户端交互：

```astro
---
// src/components/StaticCard.astro
interface Props {
  title: string;
  content: string;
}

const { title, content } = Astro.props;
---

<div class="card">
  <h3>{title}</h3>
  <p>{content}</p>
</div>

<style>
  .card {
    padding: 1.5rem;
    background: var(--color-background);
    border: 1px solid var(--color-border);
    border-radius: 0.75rem;
  }
</style>
```

### 内容集合扩展

在 `src/content/config.ts` 中定义新的内容类型：

```typescript
import { defineCollection, z } from 'astro:content';

const docs = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    description: z.string().optional(),
    order: z.number().default(0),
  }),
});

export const collections = {
  blog: blogCollection,
  docs,  // 新增文档集合
};
```

### 自定义页面路由

在 `src/pages/` 下创建 `.astro` 文件：

```astro
---
// src/pages/about.astro
import BaseLayout from '../layouts/BaseLayout.astro';
---

<BaseLayout title="关于我">
  <section class="py-16">
    <h1 class="text-4xl font-bold mb-8">关于我</h1>
    <p>这是一个自定义页面。</p>
  </section>
</BaseLayout>
```

---

## 后端服务扩展

### 添加新的 API 路由

1. 创建处理器 (`server/internal/handlers/`):

```go
// server/internal/handlers/post.go
package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	// 业务逻辑
	c.JSON(http.StatusOK, gin.H{
		"data": []string{},
	})
}
```

2. 注册路由 (`server/cmd/server/main.go`):

```go
postHandler := handlers.NewPostHandler()
api.GET("/posts", postHandler.GetPosts)
```

### 添加中间件

```go
// server/internal/middleware/rateLimit.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

func RateLimit() gin.HandlerFunc {
	// 实现限流逻辑
	return func(c *gin.Context) {
		// 检查请求频率
		c.Next()
	}
}
```

使用：
```go
api.Use(middleware.RateLimit())
```

### 添加数据模型

```go
// server/internal/models/models.go
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	Name      string         `json:"name" gorm:"size:50;uniqueIndex"`
	Posts     []Post         `json:"posts" gorm:"many2many:post_tags;"`
}
```

自动迁移会在启动时执行。

---

## 部署配置

### Docker 部署

使用 `docker-compose.yml` 一键部署：

```yaml
version: '3.8'
services:
  frontend:
    build: .
    ports:
      - "4321:80"
    depends_on:
      - api

  api:
    build: ./server
    ports:
      - "8080:8080"
    volumes:
      - ./server/data:/app/data
    environment:
      - JWT_SECRET=your-secret-key
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=admin123
```

启动：
```bash
docker-compose up -d
```

### 手动部署

#### 前端构建

```bash
npm run build
```

产物在 `dist/` 目录，可部署到任何静态托管服务。

#### 后端构建

```bash
cd server
go build -o novablog-server cmd/server/main.go
```

运行：
```bash
./novablog-server
```

### Nginx 配置

```nginx
server {
    listen 80;
    server_name example.com;

    # 静态文件
    root /var/www/novablog/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PORT` | 服务端口 | `8080` |
| `JWT_SECRET` | JWT 密钥 | 随机生成 |
| `ADMIN_USERNAME` | 管理员用户名 | `admin` |
| `ADMIN_PASSWORD` | 管理员密码 | `admin123` |
| `DB_PATH` | 数据库路径 | `./data/novablog.db` |

---

## 性能优化

### 前端优化

1. **图片优化**
   - 使用 WebP 格式
   - 添加 `loading="lazy"` 懒加载
   - 使用 `srcset` 响应式图片

2. **代码分割**
   - Astro 默认零 JS
   - 交互组件使用 `client:visible` 懒加载

3. **缓存策略**
   - 静态资源设置长期缓存
   - HTML 设置短期缓存或 ETag

### 后端优化

1. **数据库索引**
   - 已在关键字段建立索引
   - 复杂查询使用 EXPLAIN 分析

2. **连接池**
   - SQLite 使用 WAL 模式
   - 设置合理的连接池大小

3. **响应压缩**
   - 启用 Gzip 压缩
   - API 响应体积减少 70%+

### 监控指标

```bash
# 查看 Go 服务内存占用
ps aux | grep novablog-server

# 查看数据库大小
ls -lh server/data/novablog.db
```

---

## 附录

### 常用命令

```bash
# 开发
npm run dev              # 启动前端开发服务器
cd server && go run .    # 启动后端服务

# 构建
npm run build            # 构建前端
cd server && go build .  # 构建后端

# 数据库
sqlite3 server/data/novablog.db  # 打开数据库 CLI

# Docker
docker-compose up -d     # 启动所有服务
docker-compose logs -f   # 查看日志
docker-compose down      # 停止服务
```

### 技术参考

- [Astro 文档](https://docs.astro.build)
- [Vue 3 文档](https://vuejs.org)
- [Tailwind CSS 文档](https://tailwindcss.com)
- [Gin 框架文档](https://gin-gonic.com)
- [GORM 文档](https://gorm.io)

---

如有问题或建议，欢迎提交 Issue 或 Pull Request！