<template>
  <article class="micro-post group">
    <!-- 用户头像 -->
    <div class="flex-shrink-0">
      <div class="w-12 h-12 rounded-full bg-gradient-to-br from-primary-400 to-purple-500 flex items-center justify-center text-white font-bold text-lg">
        {{ getInitial(author) }}
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="flex-1 min-w-0">
      <!-- 作者信息 -->
      <div class="flex items-center gap-2 mb-2">
        <span class="font-semibold text-foreground">{{ author }}</span>
        <span class="text-sm text-foreground/40">{{ formatTime(createdAt) }}</span>
      </div>

      <!-- 内容 -->
      <div class="micro-content prose prose-sm dark:prose-invert max-w-none mb-3" v-html="renderedContent"></div>

      <!-- 标签 -->
      <div v-if="tags && tags.length > 0" class="flex flex-wrap gap-1 mb-3">
        <span
          v-for="tag in tags"
          :key="tag"
          class="inline-block px-2 py-0.5 text-xs bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded-full cursor-pointer hover:bg-primary-200 dark:hover:bg-primary-800 transition-colors"
          @click="$emit('tag-click', tag)"
        >
          #{{ tag }}
        </span>
      </div>

      <!-- 图片网格 -->
      <div v-if="images && images.length > 0" class="image-grid mb-3">
        <div
          :class="[
            'grid gap-2',
            images.length === 1 ? 'grid-cols-1 max-w-md' : '',
            images.length === 2 ? 'grid-cols-2 max-w-lg' : '',
            images.length >= 3 ? 'grid-cols-3 max-w-xl' : ''
          ]"
        >
          <div
            v-for="(image, index) in images.slice(0, 9)"
            :key="index"
            class="relative overflow-hidden rounded-lg cursor-pointer group/img"
            :class="images.length === 1 ? 'aspect-video' : 'aspect-square'"
            @click="previewImageAt(index)"
          >
            <img
              :src="image"
              :alt="`图片 ${index + 1}`"
              class="w-full h-full object-cover transition-transform duration-300 group-hover/img:scale-105"
              loading="lazy"
            />
            <div v-if="index === 8 && images.length > 9" class="absolute inset-0 bg-black/50 flex items-center justify-center">
              <span class="text-white text-xl font-bold">+{{ images.length - 9 }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作栏 -->
      <div class="flex items-center gap-6 text-foreground/50">
        <!-- 点赞 -->
        <button
          @click="toggleLike"
          class="flex items-center gap-1.5 hover:text-red-500 transition-colors"
          :class="{ 'text-red-500': liked }"
        >
          <svg class="w-5 h-5" :fill="liked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          <span class="text-sm">{{ likeCount || '' }}</span>
        </button>

        <!-- 评论 -->
        <button
          @click="showComments = !showComments"
          class="flex items-center gap-1.5 hover:text-primary-500 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <span class="text-sm">{{ commentCount || '' }}</span>
        </button>

        <!-- 分享 -->
        <button
          @click="sharePost"
          class="flex items-center gap-1.5 hover:text-primary-500 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
          <span class="text-sm">分享</span>
        </button>
      </div>

      <!-- 评论区 -->
      <Transition name="slide">
        <div v-if="showComments" class="mt-4 pt-4 border-t border-border">
          <MicroCommentSection
            :micro-id="id"
            :api-base-url="apiBaseUrl"
          />
        </div>
      </Transition>
    </div>

    <!-- 图片预览 -->
    <Teleport to="body">
      <Transition name="lightbox">
        <div
          v-if="previewIndex !== null"
          class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/90 backdrop-blur-sm"
          @click="closePreview"
        >
          <button
            class="absolute top-4 right-4 p-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all"
            @click="closePreview"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <button
            v-if="images && images.length > 1"
            class="absolute left-4 p-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all"
            @click.stop="prevImage"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>

          <img
            v-if="images && previewIndex !== null"
            :src="images[previewIndex]"
            class="max-w-[90vw] max-h-[90vh] rounded-lg shadow-2xl"
            @click.stop
          />

          <button
            v-if="images && images.length > 1"
            class="absolute right-4 p-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all"
            @click.stop="nextImage"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <div v-if="images && images.length > 1" class="absolute bottom-4 text-white/60 text-sm">
            {{ previewIndex + 1 }} / {{ images.length }}
          </div>
        </div>
      </Transition>
    </Teleport>
  </article>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { marked } from 'marked';
import MicroCommentSection from './MicroCommentSection.vue';

interface Props {
  id: number;
  author?: string;
  content: string;
  images?: string[];
  tags?: string[];
  createdAt: string;
  likeCount?: number;
  commentCount?: number;
  isLiked?: boolean;
  apiBaseUrl?: string;
}

const props = withDefaults(defineProps<Props>(), {
  author: '博主',
  likeCount: 0,
  commentCount: 0,
  isLiked: false,
  tags: () => [],
  apiBaseUrl: 'http://localhost:8080/api',
});

const emit = defineEmits<{
  (e: 'like', id: number, liked: boolean): void;
  (e: 'share', id: number): void;
  (e: 'tag-click', tag: string): void;
}>();

// 状态
const liked = ref(props.isLiked);
const likeCount = ref(props.likeCount);
const showComments = ref(false);

// 图片预览
const previewIndex = ref<number | null>(null);

// 渲染 Markdown 内容
const renderedContent = computed(() => {
  try {
    return marked.parse(props.content, { breaks: true, gfm: true }) as string;
  } catch {
    return props.content;
  }
});

// 获取头像首字母
function getInitial(name: string): string {
  return name.charAt(0).toUpperCase();
}

// 格式化时间
function formatTime(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  if (hours < 24) return `${hours} 小时前`;
  if (days < 7) return `${days} 天前`;

  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  });
}

// 切换点赞
async function toggleLike() {
  try {
    const response = await fetch(`${props.apiBaseUrl}/micro/${props.id}/like`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
    });

    if (response.ok) {
      const data = await response.json();
      liked.value = data.liked;
      likeCount.value = data.like_count;
      emit('like', props.id, data.liked);
    }
  } catch (error) {
    console.error('Failed to toggle like:', error);
  }
}

// 分享
function sharePost() {
  if (navigator.share) {
    navigator.share({
      title: '分享微语',
      text: props.content.slice(0, 100),
      url: window.location.href,
    });
  } else {
    // 复制链接
    navigator.clipboard.writeText(window.location.href);
    alert('链接已复制到剪贴板');
  }
  emit('share', props.id);
}

// 获取认证头
function getAuthHeaders(): Record<string, string> {
  if (typeof window === 'undefined') return {};
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

// 图片预览
function previewImageAt(index: number) {
  previewIndex.value = index;
}

function closePreview() {
  previewIndex.value = null;
}

function prevImage() {
  if (props.images && previewIndex.value !== null) {
    previewIndex.value = (previewIndex.value - 1 + props.images.length) % props.images.length;
  }
}

function nextImage() {
  if (props.images && previewIndex.value !== null) {
    previewIndex.value = (previewIndex.value + 1) % props.images.length;
  }
}
</script>

<style scoped>
.micro-post {
  @apply flex gap-4 p-4 bg-background border border-border rounded-xl transition-shadow duration-200;
}

.micro-post:hover {
  @apply shadow-md;
}

.micro-content :deep(p) {
  @apply mb-2 last:mb-0;
}

.micro-content :deep(a) {
  @apply text-primary-500 hover:underline;
}

.micro-content :deep(code) {
  @apply px-1.5 py-0.5 bg-muted rounded text-sm font-mono;
}

.image-grid img {
  @apply transition-transform duration-300;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.lightbox-enter-active,
.lightbox-leave-active {
  transition: opacity 0.2s ease;
}

.lightbox-enter-from,
.lightbox-leave-to {
  opacity: 0;
}
</style>