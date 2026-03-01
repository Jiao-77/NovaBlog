<template>
  <div class="comment-section">
    <h3 class="text-xl font-bold mb-6">评论</h3>
    
    <!-- 评论输入框 -->
    <div v-if="isLoggedIn" class="mb-8">
      <textarea
        v-model="newComment"
        placeholder="写下你的评论..."
        class="w-full p-4 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
        rows="4"
      ></textarea>
      <div class="flex justify-end mt-2">
        <button
          @click="submitComment"
          :disabled="!newComment.trim() || submitting"
          class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ submitting ? '发布中...' : '发布评论' }}
        </button>
      </div>
    </div>
    
    <!-- 未登录提示 -->
    <div v-else class="mb-8 p-4 bg-muted rounded-lg text-center">
      <p class="text-foreground/60">
        <a href="/login" class="text-primary-500 hover:underline">登录</a> 后参与评论
      </p>
    </div>
    
    <!-- 评论列表 -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
      <p class="mt-2 text-foreground/40">加载评论中...</p>
    </div>
    
    <div v-else-if="comments.length === 0" class="text-center py-8 text-foreground/40">
      <p>暂无评论，来抢沙发吧！</p>
    </div>
    
    <div v-else class="space-y-6">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="flex gap-4">
          <!-- 头像 -->
          <div class="flex-shrink-0">
            <div class="w-10 h-10 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <span class="text-primary-600 dark:text-primary-400 font-medium">
                {{ getInitial(comment.user) }}
              </span>
            </div>
          </div>
          
          <!-- 内容 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="font-medium">{{ getDisplayName(comment.user) }}</span>
              <span class="text-xs text-foreground/40">{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="comment-content prose prose-sm dark:prose-invert max-w-none" v-html="renderMarkdown(comment.content)"></div>
            
            <!-- 回复按钮 -->
            <button
              @click="replyTo = comment.id"
              class="text-sm text-primary-500 hover:underline mt-2"
            >
              回复
            </button>
          </div>
        </div>
        
        <!-- 回复输入框 -->
        <div v-if="replyTo === comment.id" class="mt-4 ml-14">
          <textarea
            v-model="replyContent"
            placeholder="写下你的回复..."
            class="w-full p-3 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none text-sm"
            rows="3"
          ></textarea>
          <div class="flex justify-end gap-2 mt-2">
            <button @click="replyTo = null" class="btn-secondary text-sm">取消</button>
            <button
              @click="submitReply(comment.id)"
              :disabled="!replyContent.trim() || submitting"
              class="btn-primary text-sm disabled:opacity-50"
            >
              回复
            </button>
          </div>
        </div>
        
        <!-- 子评论 -->
        <div v-if="comment.replies && comment.replies.length > 0" class="mt-4 ml-14 space-y-4">
          <div v-for="reply in comment.replies" :key="reply.id" class="flex gap-3">
            <div class="flex-shrink-0">
              <div class="w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                <span class="text-primary-600 dark:text-primary-400 text-sm font-medium">
                  {{ getInitial(reply.user) }}
                </span>
              </div>
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="font-medium text-sm">{{ getDisplayName(reply.user) }}</span>
                <span class="text-xs text-foreground/40">{{ formatDate(reply.created_at) }}</span>
              </div>
              <div class="comment-content prose prose-sm dark:prose-invert max-w-none text-sm" v-html="renderMarkdown(reply.content)"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 分页 -->
    <div v-if="pagination.totalPage > 1" class="flex justify-center gap-2 mt-8">
      <button
        @click="loadPage(pagination.page - 1)"
        :disabled="pagination.page <= 1"
        class="btn-secondary disabled:opacity-50"
      >
        上一页
      </button>
      <span class="px-4 py-2 text-foreground/60">
        {{ pagination.page }} / {{ pagination.totalPage }}
      </span>
      <button
        @click="loadPage(pagination.page + 1)"
        :disabled="pagination.page >= pagination.totalPage"
        class="btn-secondary disabled:opacity-50"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { marked } from 'marked';

// 配置 marked 选项 - 安全模式，禁止 HTML 标签
marked.setOptions({
  breaks: true, // 支持 GFM 换行
  gfm: true, // GitHub Flavored Markdown
});

const props = defineProps<{
  postId: string;
  apiBaseUrl?: string;
}>();

const apiBaseUrl = props.apiBaseUrl || 'http://localhost:8080/api';

// 状态
const comments = ref<any[]>([]);
const loading = ref(true);
const submitting = ref(false);
const newComment = ref('');
const replyTo = ref<number | null>(null);
const replyContent = ref('');
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
  totalPage: 0,
});

// 计算属性 - 仅在浏览器环境中访问 localStorage
const isLoggedIn = computed(() => {
  if (typeof window === 'undefined') return false;
  return !!localStorage.getItem('token');
});

// 获取认证头
function getAuthHeaders(): Record<string, string> {
  if (typeof window === 'undefined') return {};
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

// 加载评论
async function loadComments() {
  loading.value = true;
  try {
    const response = await fetch(
      `${apiBaseUrl}/comments?post_id=${props.postId}&page=${pagination.value.page}&page_size=${pagination.value.pageSize}`
    );
    const data = await response.json();
    if (response.ok) {
      comments.value = Array.isArray(data.data) ? data.data : [];
      pagination.value = data.pagination || { page: 1, pageSize: 20, total: 0, totalPage: 0 };
    }
  } catch (error) {
    console.error('Failed to load comments:', error);
    comments.value = [];
  } finally {
    loading.value = false;
  }
}

// 提交评论
async function submitComment() {
  if (!newComment.value.trim()) return;
  
  submitting.value = true;
  try {
    const response = await fetch(`${apiBaseUrl}/comments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        post_id: props.postId,
        content: newComment.value,
      }),
    });
    
    if (response.ok) {
      newComment.value = '';
      await loadComments();
    } else {
      const error = await response.json();
      alert(error.error || '发布失败');
    }
  } catch (error) {
    console.error('Failed to submit comment:', error);
    alert('发布失败，请稍后重试');
  } finally {
    submitting.value = false;
  }
}

// 提交回复
async function submitReply(parentId: number) {
  if (!replyContent.value.trim()) return;
  
  submitting.value = true;
  try {
    const response = await fetch(`${apiBaseUrl}/comments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        post_id: props.postId,
        parent_id: parentId,
        content: replyContent.value,
      }),
    });
    
    if (response.ok) {
      replyContent.value = '';
      replyTo.value = null;
      await loadComments();
    } else {
      const error = await response.json();
      alert(error.error || '回复失败');
    }
  } catch (error) {
    console.error('Failed to submit reply:', error);
    alert('回复失败，请稍后重试');
  } finally {
    submitting.value = false;
  }
}

// 加载指定页
function loadPage(page: number) {
  pagination.value.page = page;
  loadComments();
}

// 获取用户首字母
function getInitial(user: any): string {
  if (!user) return '?';
  const name = user.nickname || user.username || '匿名';
  return name[0].toUpperCase();
}

// 获取用户显示名称
function getDisplayName(user: any): string {
  if (!user) return '匿名用户';
  return user.nickname || user.username || '匿名用户';
}

// 格式化日期
function formatDate(dateString: string) {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);
  
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  if (hours < 24) return `${hours} 小时前`;
  if (days < 30) return `${days} 天前`;
  
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

// 渲染 Markdown 内容（安全模式）
function renderMarkdown(content: string): string {
  if (!content) return '';
  try {
    // 使用 marked 解析 Markdown，返回 HTML 字符串
    return marked.parse(content) as string;
  } catch (error) {
    console.error('Failed to parse markdown:', error);
    return content;
  }
}

onMounted(() => {
  loadComments();
});
</script>

<style scoped>
.comment-section {
  @apply mt-12 pt-8 border-t border-border;
}

/* 评论内容 Markdown 样式 */
.comment-content :deep(p) {
  @apply my-2 leading-relaxed text-gray-700 dark:text-gray-300;
}

.comment-content :deep(p:first-child) {
  @apply mt-0;
}

.comment-content :deep(p:last-child) {
  @apply mb-0;
}

.comment-content :deep(a) {
  @apply text-primary-500 hover:underline;
}

.comment-content :deep(strong) {
  @apply font-bold text-gray-900 dark:text-gray-100;
}

.comment-content :deep(em) {
  @apply italic;
}

.comment-content :deep(code) {
  @apply bg-gray-100 dark:bg-gray-800 px-1.5 py-0.5 rounded text-sm font-mono text-primary-600 dark:text-primary-400;
}

.comment-content :deep(pre) {
  @apply bg-gray-100 dark:bg-gray-800 p-4 rounded-lg overflow-x-auto my-3;
}

.comment-content :deep(pre code) {
  @apply bg-transparent p-0 text-sm;
}

.comment-content :deep(blockquote) {
  @apply border-l-4 border-primary-500 pl-4 my-3 italic text-gray-500 dark:text-gray-400;
}

.comment-content :deep(ul),
.comment-content :deep(ol) {
  @apply my-2 pl-6;
}

.comment-content :deep(ul) {
  @apply list-disc;
}

.comment-content :deep(ol) {
  @apply list-decimal;
}

.comment-content :deep(li) {
  @apply my-1 text-gray-700 dark:text-gray-300;
}

.comment-content :deep(h1),
.comment-content :deep(h2),
.comment-content :deep(h3),
.comment-content :deep(h4) {
  @apply font-bold text-gray-900 dark:text-gray-100 mt-4 mb-2;
}

.comment-content :deep(h1) {
  @apply text-xl;
}

.comment-content :deep(h2) {
  @apply text-lg;
}

.comment-content :deep(h3) {
  @apply text-base;
}

.comment-content :deep(hr) {
  @apply border-gray-200 dark:border-gray-700 my-4;
}

.comment-content :deep(img) {
  @apply max-w-full rounded-lg my-2;
}

.comment-content :deep(table) {
  @apply w-full border-collapse my-4;
}

.comment-content :deep(th),
.comment-content :deep(td) {
  @apply border border-gray-200 dark:border-gray-700 px-3 py-2 text-left;
}

.comment-content :deep(th) {
  @apply bg-gray-100 dark:bg-gray-800 font-bold;
}
</style>
