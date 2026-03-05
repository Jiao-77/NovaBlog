<template>
  <div class="micro-comment-section">
    <h4 class="text-base font-semibold mb-4">评论</h4>

    <!-- 评论输入框 -->
    <div v-if="isLoggedIn" class="mb-6">
      <textarea
        v-model="newComment"
        placeholder="写下你的评论..."
        class="w-full p-3 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none text-sm"
        rows="3"
      ></textarea>
      <div class="flex justify-end mt-2">
        <button
          @click="submitComment"
          :disabled="!newComment.trim() || submitting"
          class="btn-primary text-sm disabled:opacity-50"
        >
          {{ submitting ? '发布中...' : '发布评论' }}
        </button>
      </div>
    </div>

    <!-- 未登录提示 -->
    <div v-else class="mb-6 p-3 bg-muted rounded-lg text-center text-sm">
      <p class="text-foreground/60">
        <a href="/login" class="text-primary-500 hover:underline">登录</a> 后参与评论
      </p>
    </div>

    <!-- 评论列表 -->
    <div v-if="loading" class="text-center py-4">
      <div class="animate-spin w-6 h-6 border-2 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
    </div>

    <div v-else-if="comments.length === 0" class="text-center py-4 text-foreground/40 text-sm">
      <p>暂无评论</p>
    </div>

    <div v-else class="space-y-4">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="flex gap-3">
          <!-- 头像 -->
          <div class="flex-shrink-0">
            <div class="w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <span class="text-primary-600 dark:text-primary-400 text-sm font-medium">
                {{ getInitial(comment.user) }}
              </span>
            </div>
          </div>

          <!-- 内容 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="font-medium text-sm">{{ getDisplayName(comment.user) }}</span>
              <span class="text-xs text-foreground/40">{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="text-sm text-foreground/80">{{ comment.content }}</div>

            <!-- 回复按钮 -->
            <button
              @click="replyTo = comment.id"
              class="text-xs text-primary-500 hover:underline mt-1"
            >
              回复
            </button>
          </div>
        </div>

        <!-- 回复输入框 -->
        <div v-if="replyTo === comment.id" class="mt-3 ml-11">
          <textarea
            v-model="replyContent"
            placeholder="写下你的回复..."
            class="w-full p-2 border border-border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none text-sm"
            rows="2"
          ></textarea>
          <div class="flex justify-end gap-2 mt-2">
            <button @click="replyTo = null" class="btn-secondary text-xs">取消</button>
            <button
              @click="submitReply(comment.id)"
              :disabled="!replyContent.trim() || submitting"
              class="btn-primary text-xs disabled:opacity-50"
            >
              回复
            </button>
          </div>
        </div>

        <!-- 子评论 -->
        <div v-if="comment.replies && comment.replies.length > 0" class="mt-3 ml-11 space-y-3">
          <div v-for="reply in comment.replies" :key="reply.id" class="flex gap-2">
            <div class="flex-shrink-0">
              <div class="w-6 h-6 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                <span class="text-primary-600 dark:text-primary-400 text-xs font-medium">
                  {{ getInitial(reply.user) }}
                </span>
              </div>
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-0.5">
                <span class="font-medium text-xs">{{ getDisplayName(reply.user) }}</span>
                <span class="text-xs text-foreground/40">{{ formatDate(reply.created_at) }}</span>
              </div>
              <div class="text-xs text-foreground/80">{{ reply.content }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';

interface Props {
  microId: number;
  apiBaseUrl?: string;
}

const props = withDefaults(defineProps<Props>(), {
  apiBaseUrl: 'http://localhost:8080/api',
});

// 状态
const comments = ref<any[]>([]);
const loading = ref(true);
const submitting = ref(false);
const newComment = ref('');
const replyTo = ref<number | null>(null);
const replyContent = ref('');

// 计算属性
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
      `${props.apiBaseUrl}/micro-comments?micro_id=${props.microId}`
    );
    if (response.ok) {
      const data = await response.json();
      comments.value = data.data || [];
    }
  } catch (error) {
    console.error('Failed to load comments:', error);
  } finally {
    loading.value = false;
  }
}

// 提交评论
async function submitComment() {
  if (!newComment.value.trim() || submitting.value) return;

  submitting.value = true;
  try {
    const response = await fetch(`${props.apiBaseUrl}/micro-comments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        micro_id: props.microId,
        content: newComment.value,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      comments.value.unshift(data);
      newComment.value = '';
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
  if (!replyContent.value.trim() || submitting.value) return;

  submitting.value = true;
  try {
    const response = await fetch(`${props.apiBaseUrl}/micro-comments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        micro_id: props.microId,
        parent_id: parentId,
        content: replyContent.value,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      // 找到父评论并添加回复
      const parentComment = comments.value.find(c => c.id === parentId);
      if (parentComment) {
        if (!parentComment.replies) {
          parentComment.replies = [];
        }
        parentComment.replies.push(data);
      }
      replyContent.value = '';
      replyTo.value = null;
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

// 获取用户首字母
function getInitial(user: any): string {
  if (!user) return '?';
  const name = user.nickname || user.username || '匿';
  return name[0].toUpperCase();
}

// 获取用户显示名称
function getDisplayName(user: any): string {
  if (!user) return '匿名用户';
  return user.nickname || user.username || '匿名用户';
}

// 格式化日期
function formatDate(dateString: string): string {
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
    month: 'short',
    day: 'numeric',
  });
}

onMounted(() => {
  loadComments();
});
</script>

<style scoped>
.micro-comment-section {
  @apply text-sm;
}
</style>