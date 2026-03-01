<template>
  <button
    @click="toggleLike"
    :disabled="loading"
    class="like-button flex items-center gap-2 px-4 py-2 rounded-full transition-all duration-300"
    :class="[
      liked 
        ? 'bg-red-100 dark:bg-red-900/30 text-red-500' 
        : 'bg-muted hover:bg-red-50 dark:hover:bg-red-900/20 text-foreground/60 hover:text-red-500'
    ]"
  >
    <!-- 心形图标 -->
    <svg
      class="w-5 h-5 transition-transform duration-300"
      :class="{ 'scale-125': animating }"
      viewBox="0 0 24 24"
      fill="currentColor"
    >
      <path
        d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
      />
    </svg>
    
    <!-- 点赞数 -->
    <span class="font-medium">{{ likeCount }}</span>
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

const props = defineProps<{
  postId: string;
  apiBaseUrl?: string;
}>();

const emit = defineEmits<{
  (e: 'update', liked: boolean, count: number): void;
}>();

const apiBaseUrl = props.apiBaseUrl || 'http://localhost:8080/api';

// 状态
const liked = ref(false);
const likeCount = ref(0);
const loading = ref(false);
const animating = ref(false);

// 获取认证头
function getAuthHeaders(): Record<string, string> {
  // 仅在浏览器环境中访问 localStorage
  if (typeof window === 'undefined') return {};
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

// 加载点赞状态
async function loadLikeStatus() {
  loading.value = true;
  try {
    const response = await fetch(
      `${apiBaseUrl}/likes?post_id=${props.postId}`,
      {
        headers: getAuthHeaders(),
      }
    );
    
    if (response.ok) {
      const data = await response.json();
      liked.value = data.liked;
      likeCount.value = data.like_count;
    }
  } catch (error) {
    console.error('Failed to load like status:', error);
  } finally {
    loading.value = false;
  }
}

// 切换点赞
async function toggleLike() {
  loading.value = true;
  
  try {
    const response = await fetch(`${apiBaseUrl}/likes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        post_id: props.postId,
      }),
    });
    
    if (response.ok) {
      const data = await response.json();
      liked.value = data.liked;
      likeCount.value = data.like_count;
      
      // 触发动画
      if (data.liked) {
        animating.value = true;
        setTimeout(() => {
          animating.value = false;
        }, 300);
      }
      
      emit('update', data.liked, data.like_count);
    } else {
      const error = await response.json();
      console.error('Failed to toggle like:', error);
    }
  } catch (error) {
    console.error('Failed to toggle like:', error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadLikeStatus();
});
</script>

<style scoped>
.like-button {
  @apply cursor-pointer select-none;
}

.like-button:disabled {
  @apply cursor-wait;
}
</style>