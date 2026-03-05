<template>
  <div class="micro-page">
    <!-- 发布框 -->
    <div class="card mb-6">
      <textarea
        ref="textareaRef"
        v-model="newContent"
        placeholder="分享你的想法..."
        class="w-full p-4 bg-muted/50 border border-border rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-primary-500"
        rows="4"
      ></textarea>
      <div class="flex items-center justify-between mt-3">
        <div class="flex items-center gap-2">
          <!-- Emoji 按钮 -->
          <div class="relative">
            <button
              @click="showEmojiPicker = !showEmojiPicker"
              class="btn-ghost p-2 rounded-lg"
              title="添加表情"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </button>
            <!-- Emoji 选择面板 -->
            <div
              v-if="showEmojiPicker"
              class="absolute left-0 top-full mt-1 bg-white dark:bg-gray-800 border border-border rounded-lg shadow-lg p-3 z-20 w-72"
              @click.stop
            >
              <div class="flex flex-wrap gap-1">
                <button
                  v-for="emoji in commonEmojis"
                  :key="emoji"
                  @click="insertEmoji(emoji)"
                  class="w-8 h-8 flex items-center justify-center text-xl hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
                >
                  {{ emoji }}
                </button>
              </div>
            </div>
          </div>
          <!-- 图片上传按钮 -->
          <button
            @click="triggerImageUpload"
            class="btn-ghost p-2 rounded-lg"
            title="上传图片"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </button>
          <input
            ref="imageInput"
            type="file"
            accept="image/*"
            multiple
            class="hidden"
            @change="handleImageSelect"
          />
          <!-- 标签按钮 -->
          <div class="relative" ref="tagInputRef">
            <button
              @click="showTagInput = !showTagInput"
              class="btn-ghost p-2 rounded-lg"
              title="添加标签"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
              </svg>
            </button>
            <!-- 标签输入面板 -->
            <div
              v-if="showTagInput"
              class="absolute left-0 top-full mt-1 bg-white dark:bg-gray-800 border border-border rounded-lg shadow-lg p-3 z-20 w-64"
              @click.stop
            >
              <div class="flex gap-2 mb-2">
                <input
                  v-model="tagInput"
                  type="text"
                  placeholder="输入标签（回车添加）"
                  class="flex-1 px-2 py-1 text-sm border border-border rounded bg-background focus:outline-none focus:ring-1 focus:ring-primary-500"
                  @keydown.enter="addTag"
                />
                <button
                  @click="addTag"
                  class="px-2 py-1 text-sm bg-primary-500 text-white rounded hover:bg-primary-600"
                >
                  添加
                </button>
              </div>
              <div v-if="newTags.length > 0" class="flex flex-wrap gap-1">
                <span
                  v-for="(tag, index) in newTags"
                  :key="index"
                  class="inline-flex items-center gap-1 px-2 py-0.5 text-xs bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded-full"
                >
                  #{{ tag }}
                  <button @click="removeTag(index)" class="hover:text-red-500">×</button>
                </span>
              </div>
              <div class="mt-2 pt-2 border-t border-border text-xs text-foreground/40">
                最多添加 5 个标签
              </div>
            </div>
          </div>
        </div>
        <button
          @click="publishMicro"
          :disabled="!newContent.trim() || publishing"
          class="btn-primary disabled:opacity-50"
        >
          {{ publishing ? '发布中...' : '发布' }}
        </button>
      </div>
      <!-- 已选标签显示 -->
      <div v-if="newTags.length > 0 && !showTagInput" class="mt-2 flex flex-wrap gap-1">
        <span
          v-for="(tag, index) in newTags"
          :key="index"
          class="inline-flex items-center gap-1 px-2 py-0.5 text-xs bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded-full"
        >
          #{{ tag }}
          <button @click="removeTag(index)" class="hover:text-red-500">×</button>
        </span>
      </div>
      <!-- 图片预览 -->
      <div v-if="newImages.length > 0 || uploadingImages" class="mt-3 grid grid-cols-3 gap-2">
        <!-- 上传进度 -->
        <div v-if="uploadingImages" class="relative aspect-square rounded-lg overflow-hidden bg-muted flex flex-col items-center justify-center p-2">
          <div class="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden mb-2">
            <div 
              class="h-full bg-primary-500 transition-all duration-100 rounded-full"
              :style="{ width: uploadProgress + '%' }"
            ></div>
          </div>
          <span class="text-xs text-foreground/60">{{ uploadProgress }}%</span>
        </div>
        <div
          v-for="(img, index) in newImages"
          :key="index"
          class="relative aspect-square rounded-lg overflow-hidden bg-muted"
        >
          <img :src="img" class="w-full h-full object-cover" />
          <button
            @click="removeImage(index)"
            class="absolute top-1 right-1 p-1 bg-black/50 rounded-full text-white hover:bg-black/70"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 微语列表 -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin w-8 h-8 border-2 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
      <p class="mt-2 text-foreground/40">加载中...</p>
    </div>

    <div v-else-if="micros.length === 0" class="text-center py-16 text-foreground/40">
      <svg class="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
      </svg>
      <p>暂无微语，来发布第一条吧！</p>
    </div>

    <div v-else class="space-y-4">
      <MicroPost
        v-for="micro in micros"
        :key="micro.id"
        :id="micro.id"
        :author="micro.user?.nickname || micro.user?.username || '博主'"
        :content="micro.content"
        :images="parseImages(micro.images)"
        :tags="parseTags(micro.tags)"
        :created-at="micro.created_at"
        :like-count="micro.like_count"
        :comment-count="micro.comment_count"
        :is-liked="micro.is_liked"
        :api-base-url="apiBaseUrl"
        @like="handleLikeUpdate"
        @tag-click="handleTagClick"
      />
    </div>

    <!-- 加载更多 -->
    <div v-if="pagination.totalPage > 1" class="text-center py-4">
      <button
        v-if="pagination.page < pagination.totalPage"
        @click="loadMore"
        :disabled="loadingMore"
        class="btn-secondary"
      >
        {{ loadingMore ? '加载中...' : '加载更多' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import MicroPost from './MicroPost.vue';

interface Props {
  apiBaseUrl?: string;
  imageUploadUrl?: string;
  imageUploadToken?: string;
}

const props = withDefaults(defineProps<Props>(), {
  apiBaseUrl: 'http://localhost:8080/api',
  imageUploadUrl: 'https://picturebed.jiao77.cn/api/index.php',
  imageUploadToken: 'blog',
});

// 定义事件
const emit = defineEmits<{
  (e: 'published'): void;
}>();

// 状态
const micros = ref<any[]>([]);
const loading = ref(true);
const loadingMore = ref(false);
const publishing = ref(false);
const newContent = ref('');
const newImages = ref<string[]>([]);
const newTags = ref<string[]>([]);
const tagInput = ref('');
const imageInput = ref<HTMLInputElement | null>(null);
const textareaRef = ref<HTMLTextAreaElement | null>(null);
const tagInputRef = ref<HTMLElement | null>(null);
const uploadingImages = ref(false);
const uploadProgress = ref(0);
const showEmojiPicker = ref(false);
const showTagInput = ref(false);
const currentTag = ref('');

// 常用 emoji 列表
const commonEmojis = [
  '😀', '😃', '😄', '😁', '😅', '😂', '🤣', '😊', '😇', '🙂',
  '😉', '😌', '😍', '🥰', '😘', '😋', '😛', '😜', '🤪', '😝',
  '🤗', '🤔', '🤭', '🤫', '🤥', '😶', '😐', '😑', '😏', '😒',
  '🙄', '😬', '😮', '🥱', '😴', '🤤', '😷', '🤒', '🤕', '🤢',
  '👍', '👎', '👏', '🙌', '🤝', '🙏', '💪', '🎉', '🎊', '💯',
  '❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '💔', '❣️', '💕',
  '🔥', '⭐', '🌟', '✨', '💫', '🎯', '🏆', '🚀', '💡', '📌',
];

const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
  totalPage: 0,
});

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

// 加载微语列表
async function loadMicros(page = 1) {
  try {
    let url = `${props.apiBaseUrl}/micro?page=${page}&page_size=${pagination.value.pageSize}`;
    if (currentTag.value) {
      url += `&tag=${encodeURIComponent(currentTag.value)}`;
    }
    
    const response = await fetch(url, { headers: getAuthHeaders() });

    if (response.ok) {
      const data = await response.json();
      if (page === 1) {
        micros.value = data.data || [];
      } else {
        micros.value.push(...(data.data || []));
      }
      pagination.value = data.pagination || pagination.value;
    }
  } catch (error) {
    console.error('Failed to load micros:', error);
  } finally {
    loading.value = false;
    loadingMore.value = false;
  }
}

// 加载更多
async function loadMore() {
  loadingMore.value = true;
  await loadMicros(pagination.value.page + 1);
}

// 发布微语
async function publishMicro() {
  if (!newContent.value.trim() || publishing.value) return;

  // 检查登录状态
  if (!isLoggedIn.value) {
    alert('请先登录');
    return;
  }

  publishing.value = true;

  try {
    const response = await fetch(`${props.apiBaseUrl}/micro`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        content: newContent.value,
        images: newImages.value,
        tags: newTags.value,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      // 添加到列表开头
      micros.value.unshift(data);
      // 清空输入
      newContent.value = '';
      newImages.value = [];
      newTags.value = [];
      // 触发事件通知侧边栏刷新
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('micro-published'));
      }
    } else {
      const error = await response.json();
      alert(error.error || '发布失败');
    }
  } catch (error) {
    console.error('Failed to publish:', error);
    alert('发布失败，请稍后重试');
  } finally {
    publishing.value = false;
  }
}

// 触发图片上传
function triggerImageUpload() {
  imageInput.value?.click();
}

// 处理图片选择
async function handleImageSelect(event: Event) {
  const target = event.target as HTMLInputElement;
  const files = target.files;
  if (!files || files.length === 0) return;

  uploadingImages.value = true;
  uploadProgress.value = 0;

  // 模拟进度条：使用 1-e^(-x) 函数，约5秒完成
  const totalDuration = 5000; // 5秒
  const interval = 100; // 每100ms更新一次
  let elapsed = 0;
  
  const progressTimer = setInterval(() => {
    elapsed += interval;
    // 1 - e^(-t/1000) 函数，t单位ms，除以1000让曲线更平滑
    const progress = (1 - Math.exp(-elapsed / 1500)) * 100;
    uploadProgress.value = Math.min(Math.round(progress), 95); // 最多到95%，真正完成后到100%
  }, interval);

  try {
    for (const file of Array.from(files)) {
      if (!file.type.startsWith('image/')) continue;
      if (file.size > 5 * 1024 * 1024) {
        alert('图片大小不能超过 5MB');
        continue;
      }

      const formData = new FormData();
      formData.append('image', file);
      formData.append('token', props.imageUploadToken);

      const response = await fetch(props.imageUploadUrl, {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();
      if (data.result === 'success' && data.url) {
        newImages.value.push(data.url);
      }
    }
  } catch (error) {
    console.error('Failed to upload image:', error);
  } finally {
    clearInterval(progressTimer);
    uploadProgress.value = 100;
    // 短暂显示100%后隐藏
    setTimeout(() => {
      uploadingImages.value = false;
      uploadProgress.value = 0;
    }, 300);
  }

  // 清空 input
  target.value = '';
}

// 移除图片
function removeImage(index: number) {
  newImages.value.splice(index, 1);
}

// 添加标签
function addTag() {
  const tag = tagInput.value.trim();
  if (tag && !newTags.value.includes(tag) && newTags.value.length < 5) {
    newTags.value.push(tag);
    tagInput.value = '';
  }
}

// 移除标签
function removeTag(index: number) {
  newTags.value.splice(index, 1);
}

// 插入 emoji
function insertEmoji(emoji: string) {
  if (textareaRef.value) {
    const start = textareaRef.value.selectionStart;
    const end = textareaRef.value.selectionEnd;
    const text = newContent.value;
    newContent.value = text.substring(0, start) + emoji + text.substring(end);
    showEmojiPicker.value = false;
    // 恢复焦点
    textareaRef.value.focus();
    textareaRef.value.setSelectionRange(start + emoji.length, start + emoji.length);
  } else {
    newContent.value += emoji;
    showEmojiPicker.value = false;
  }
}

// 解析图片 JSON
function parseImages(imagesJson: string): string[] {
  if (!imagesJson) return [];
  try {
    return JSON.parse(imagesJson);
  } catch {
    return [];
  }
}

// 解析标签 JSON
function parseTags(tagsJson: string): string[] {
  if (!tagsJson) return [];
  try {
    return JSON.parse(tagsJson);
  } catch {
    return [];
  }
}

// 处理标签点击
function handleTagClick(tag: string) {
  // 触发事件通知侧边栏和父组件
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('tag-filter', { detail: tag }));
  }
}

// 处理点赞更新
function handleLikeUpdate(id: number, liked: boolean) {
  const micro = micros.value.find(m => m.id === id);
  if (micro) {
    micro.is_liked = liked;
    micro.like_count += liked ? 1 : -1;
  }
}

onMounted(() => {
  loadMicros();
  
  // 点击外部关闭 emoji 选择器和标签输入
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement;
    
    // 关闭 emoji 选择器
    if (!target.closest('.relative') || !target.closest('button')) {
      showEmojiPicker.value = false;
    }
    
    // 关闭标签输入面板
    if (tagInputRef.value && !tagInputRef.value.contains(target)) {
      showTagInput.value = false;
    }
  });
  
  // 监听标签筛选事件
  window.addEventListener('tag-filter', ((e: CustomEvent) => {
    currentTag.value = e.detail || '';
    loadMicros(1);
  }) as EventListener);
});
</script>