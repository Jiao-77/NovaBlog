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
      <!-- 工具栏 -->
      <div class="flex justify-between items-center mt-2">
        <div class="flex items-center gap-2">
          <!-- Emoji 按钮 -->
          <div class="relative emoji-picker-container">
            <button
              @click.stop="showEmojiPicker = !showEmojiPicker"
              class="flex items-center gap-1 px-3 py-1.5 text-sm text-foreground/60 hover:text-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded transition-colors"
              title="添加表情"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>表情</span>
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
            :disabled="uploadingImage"
            class="flex items-center gap-1 px-3 py-1.5 text-sm text-foreground/60 hover:text-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded transition-colors disabled:opacity-50"
            title="上传图片"
          >
            <svg v-if="!uploadingImage" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-5 h-5 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ uploadingImage ? '上传中...' : '图片' }}</span>
          </button>
          <input
            ref="imageInput"
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleImageUpload"
          />
        </div>
        <button
          @click="submitComment"
          :disabled="!newComment.trim() || submitting"
          class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ submitting ? '发布中...' : '发布评论' }}
        </button>
      </div>
      <!-- 上传错误提示 -->
      <p v-if="uploadError" class="text-sm text-red-500 mt-1">{{ uploadError }}</p>
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
            <div class="comment-content prose prose-sm dark:prose-invert max-w-none" v-html="renderMarkdown(comment.content)" @click="handleContentClick"></div>
            
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
          <div class="flex justify-between items-center mt-2">
            <div class="flex items-center gap-2">
              <!-- Emoji 按钮 -->
              <div class="relative emoji-picker-container">
                <button
                  @click.stop="showReplyEmojiPicker = showReplyEmojiPicker === comment.id ? null : comment.id"
                  class="flex items-center gap-1 px-2 py-1 text-xs text-foreground/60 hover:text-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded transition-colors"
                  title="添加表情"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span>表情</span>
                </button>
                <!-- Emoji 选择面板 -->
                <div
                  v-if="showReplyEmojiPicker === comment.id"
                  class="absolute left-0 top-full mt-1 bg-white dark:bg-gray-800 border border-border rounded-lg shadow-lg p-3 z-20 w-64"
                  @click.stop
                >
                  <div class="flex flex-wrap gap-1">
                    <button
                      v-for="emoji in commonEmojis"
                      :key="emoji"
                      @click="insertReplyEmoji(emoji, comment.id)"
                      class="w-7 h-7 flex items-center justify-center text-lg hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
                    >
                      {{ emoji }}
                    </button>
                  </div>
                </div>
              </div>
              <!-- 图片上传按钮 -->
              <button
                @click="triggerReplyImageUpload(comment.id)"
                :disabled="uploadingReplyImage"
                class="flex items-center gap-1 px-2 py-1 text-xs text-foreground/60 hover:text-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded transition-colors disabled:opacity-50"
                title="上传图片"
              >
                <svg v-if="!uploadingReplyImage" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <svg v-else class="w-4 h-4 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ uploadingReplyImage ? '上传中...' : '图片' }}</span>
              </button>
              <input
                :ref="el => replyImageInputs[comment.id] = el as HTMLInputElement"
                type="file"
                accept="image/*"
                class="hidden"
                @change="(e) => handleReplyImageUpload(e, comment.id)"
              />
            </div>
            <div class="flex gap-2">
              <button @click="replyTo = null; showReplyEmojiPicker = null" class="btn-secondary text-sm">取消</button>
              <button
                @click="submitReply(comment.id)"
                :disabled="!replyContent.trim() || submitting"
                class="btn-primary text-sm disabled:opacity-50"
              >
                回复
              </button>
            </div>
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
              <div class="comment-content prose prose-sm dark:prose-invert max-w-none text-sm" v-html="renderMarkdown(reply.content)" @click="handleContentClick"></div>
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

    <!-- 图片预览模态框（灯箱） -->
    <Teleport to="body">
      <Transition name="lightbox">
        <div
          v-if="previewImage"
          class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/90 backdrop-blur-sm select-none"
          @click="closePreview"
          @wheel.prevent="handleWheel"
        >
          <!-- 关闭按钮 -->
          <button
            class="absolute top-4 right-4 z-10 p-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all"
            @click="closePreview"
            aria-label="关闭预览"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          
          <!-- 缩放控制按钮 -->
          <div class="absolute top-4 left-4 z-10 flex gap-2">
            <button
              class="p-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all"
              @click.stop="zoomOut"
              :disabled="imageScale <= 0.5"
              :class="{ 'opacity-50 cursor-not-allowed': imageScale <= 0.5 }"
              aria-label="缩小"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM13 10H7" />
              </svg>
            </button>
            <button
              class="px-3 py-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all text-sm min-w-[60px]"
              @click.stop="resetZoom"
              aria-label="重置缩放"
            >
              {{ Math.round(imageScale * 100) }}%
            </button>
            <button
              class="p-2 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full transition-all"
              @click.stop="zoomIn"
              :disabled="imageScale >= 3"
              :class="{ 'opacity-50 cursor-not-allowed': imageScale >= 3 }"
              aria-label="放大"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v6m3-3H7" />
              </svg>
            </button>
          </div>
          
          <!-- 图片容器 -->
          <div
            class="relative overflow-hidden"
            @click.stop
            @mousedown="startDrag"
          >
            <img
              ref="previewImgRef"
              :src="previewImage"
              class="rounded-lg shadow-2xl transition-transform duration-100"
              :style="imageStyle"
              @click.stop
              @dragstart.prevent
              alt="图片预览"
            />
          </div>
          
          <!-- 提示文字 -->
          <p class="absolute bottom-4 left-1/2 -translate-x-1/2 text-white/60 text-sm">
            滚轮缩放 · 拖拽移动 · ESC 或点击背景关闭
          </p>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { marked } from 'marked';

// 配置 marked 选项 - 安全模式，禁止 HTML 标签
marked.setOptions({
  breaks: true, // 支持 GFM 换行
  gfm: true, // GitHub Flavored Markdown
});

const props = defineProps<{
  postId: string;
  apiBaseUrl?: string;
  imageUploadUrl?: string;
  imageUploadToken?: string;
}>();

const apiBaseUrl = props.apiBaseUrl || 'http://localhost:8080/api';
const imageUploadUrl = props.imageUploadUrl || 'https://picturebed.jiao77.cn/api/index.php';
const imageUploadToken = props.imageUploadToken || 'blog';

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

// 图片上传相关状态
const imageInput = ref<HTMLInputElement | null>(null);
const uploadingImage = ref(false);
const uploadError = ref('');
const previewImage = ref<string | null>(null);

// 回复图片上传相关状态
const replyImageInputs = ref<Record<number, HTMLInputElement>>({});
const uploadingReplyImage = ref(false);

// Emoji 选择器状态
const showEmojiPicker = ref(false);
const showReplyEmojiPicker = ref<number | null>(null);

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

// 灯箱缩放和拖拽状态
const previewImgRef = ref<HTMLImageElement | null>(null);
const imageScale = ref(1);
const imageTranslate = ref({ x: 0, y: 0 });
const isDragging = ref(false);
const dragStart = ref({ x: 0, y: 0 });
const dragOffset = ref({ x: 0, y: 0 });

// 计算属性 - 仅在浏览器环境中访问 localStorage
const isLoggedIn = computed(() => {
  if (typeof window === 'undefined') return false;
  return !!localStorage.getItem('token');
});

// 图片样式计算属性
const imageStyle = computed(() => ({
  transform: `translate(${imageTranslate.value.x}px, ${imageTranslate.value.y}px) scale(${imageScale.value})`,
  cursor: isDragging.value ? 'grabbing' : 'grab',
  maxWidth: '90vw',
  maxHeight: '90vh',
}));

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

// 触发图片上传
function triggerImageUpload() {
  if (imageInput.value) {
    imageInput.value.click();
  }
}

// 处理图片上传
async function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    uploadError.value = '请选择图片文件';
    return;
  }

  // 验证文件大小 (最大 5MB)
  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = '图片大小不能超过 5MB';
    return;
  }

  uploadError.value = '';
  uploadingImage.value = true;

  try {
    const formData = new FormData();
    formData.append('image', file);
    formData.append('token', imageUploadToken);

    const response = await fetch(imageUploadUrl, {
      method: 'POST',
      body: formData,
    });

    const data = await response.json();

    if (data.result === 'success' && data.url) {
      // 插入 Markdown 图片语法到评论内容
      const imageName = data.srcName || file.name.replace(/\.[^/.]+$/, '');
      const markdown = `![${imageName}](${data.url})`;

      // 在光标位置插入或追加到末尾
      const textarea = document.querySelector('textarea') as HTMLTextAreaElement;
      if (textarea) {
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const text = newComment.value;
        newComment.value = text.substring(0, start) + markdown + text.substring(end);
        // 恢复焦点
        textarea.focus();
        textarea.setSelectionRange(start + markdown.length, start + markdown.length);
      } else {
        newComment.value += markdown;
      }
    } else {
      uploadError.value = data.message || '上传失败，请稍后重试';
    }
  } catch (error) {
    console.error('Failed to upload image:', error);
    uploadError.value = '上传失败，请检查网络连接';
  } finally {
    uploadingImage.value = false;
    // 清空 input 以便重复选择同一文件
    if (target) {
      target.value = '';
    }
  }
}

// 插入 emoji 到评论
function insertEmoji(emoji: string) {
  const textarea = document.querySelector('.comment-section textarea') as HTMLTextAreaElement;
  if (textarea) {
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const text = newComment.value;
    newComment.value = text.substring(0, start) + emoji + text.substring(end);
    showEmojiPicker.value = false;
    // 恢复焦点
    textarea.focus();
    textarea.setSelectionRange(start + emoji.length, start + emoji.length);
  } else {
    newComment.value += emoji;
    showEmojiPicker.value = false;
  }
}

// 插入 emoji 到回复
function insertReplyEmoji(emoji: string, commentId: number) {
  replyContent.value += emoji;
  showReplyEmojiPicker.value = null;
}

// 触发回复图片上传
function triggerReplyImageUpload(commentId: number) {
  const input = replyImageInputs.value[commentId];
  if (input) {
    input.click();
  }
}

// 处理回复图片上传
async function handleReplyImageUpload(event: Event, commentId: number) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    alert('请选择图片文件');
    return;
  }

  // 验证文件大小 (最大 5MB)
  if (file.size > 5 * 1024 * 1024) {
    alert('图片大小不能超过 5MB');
    return;
  }

  uploadingReplyImage.value = true;

  try {
    const formData = new FormData();
    formData.append('image', file);
    formData.append('token', imageUploadToken);

    const response = await fetch(imageUploadUrl, {
      method: 'POST',
      body: formData,
    });

    const data = await response.json();

    if (data.result === 'success' && data.url) {
      const imageName = data.srcName || file.name.replace(/\.[^/.]+$/, '');
      const markdown = `![${imageName}](${data.url})`;
      replyContent.value += markdown;
    } else {
      alert(data.message || '上传失败，请稍后重试');
    }
  } catch (error) {
    console.error('Failed to upload image:', error);
    alert('上传失败，请检查网络连接');
  } finally {
    uploadingReplyImage.value = false;
    if (target) {
      target.value = '';
    }
  }
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

// 处理评论区域点击事件（事件委托）
function handleContentClick(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (target.tagName === 'IMG') {
    const img = target as HTMLImageElement;
    previewImage.value = img.src;
  }
}

// 关闭图片预览
function closePreview() {
  previewImage.value = null;
  // 重置缩放和位置
  imageScale.value = 1;
  imageTranslate.value = { x: 0, y: 0 };
}

// 放大
function zoomIn() {
  if (imageScale.value < 3) {
    imageScale.value = Math.min(3, imageScale.value + 0.25);
  }
}

// 缩小
function zoomOut() {
  if (imageScale.value > 0.5) {
    imageScale.value = Math.max(0.5, imageScale.value - 0.25);
    // 缩小时限制图片不超出边界
    limitTranslate();
  }
}

// 重置缩放
function resetZoom() {
  imageScale.value = 1;
  imageTranslate.value = { x: 0, y: 0 };
}

// 滚轮缩放
function handleWheel(event: WheelEvent) {
  if (!previewImage.value) return;
  
  const delta = event.deltaY > 0 ? -0.1 : 0.1;
  const newScale = Math.max(0.5, Math.min(3, imageScale.value + delta));
  
  if (newScale !== imageScale.value) {
    imageScale.value = newScale;
    if (newScale < 1) {
      limitTranslate();
    }
  }
}

// 开始拖拽
function startDrag(event: MouseEvent) {
  if (imageScale.value <= 1) return;
  
  isDragging.value = true;
  dragStart.value = { x: event.clientX, y: event.clientY };
  dragOffset.value = { ...imageTranslate.value };
  
  document.addEventListener('mousemove', handleDrag);
  document.addEventListener('mouseup', stopDrag);
}

// 拖拽中
function handleDrag(event: MouseEvent) {
  if (!isDragging.value) return;
  
  const dx = event.clientX - dragStart.value.x;
  const dy = event.clientY - dragStart.value.y;
  
  imageTranslate.value = {
    x: dragOffset.value.x + dx,
    y: dragOffset.value.y + dy,
  };
}

// 停止拖拽
function stopDrag() {
  isDragging.value = false;
  document.removeEventListener('mousemove', handleDrag);
  document.removeEventListener('mouseup', stopDrag);
  limitTranslate();
}

// 限制拖拽范围
function limitTranslate() {
  if (!previewImgRef.value) return;
  
  const img = previewImgRef.value;
  const maxOffsetX = (img.naturalWidth * imageScale.value - window.innerWidth * 0.9) / 2;
  const maxOffsetY = (img.naturalHeight * imageScale.value - window.innerHeight * 0.9) / 2;
  
  imageTranslate.value = {
    x: Math.max(-maxOffsetX, Math.min(maxOffsetX, imageTranslate.value.x)),
    y: Math.max(-maxOffsetY, Math.min(maxOffsetY, imageTranslate.value.y)),
  };
  
  // 如果图片比视口小，不允许移动
  if (img.naturalWidth * imageScale.value <= window.innerWidth * 0.9) {
    imageTranslate.value.x = 0;
  }
  if (img.naturalHeight * imageScale.value <= window.innerHeight * 0.9) {
    imageTranslate.value.y = 0;
  }
}

// 处理键盘事件
function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && previewImage.value) {
    previewImage.value = null;
  }
}

onMounted(() => {
  loadComments();
  // 监听键盘事件
  window.addEventListener('keydown', handleKeydown);
  // 点击外部关闭 emoji 选择器
  document.addEventListener('click', closeEmojiPickers);
});

// 关闭所有 emoji 选择器
function closeEmojiPickers(event: MouseEvent) {
  const target = event.target as HTMLElement;
  if (!target.closest('.emoji-picker-container')) {
    showEmojiPicker.value = false;
    showReplyEmojiPicker.value = null;
  }
}

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown);
  document.removeEventListener('click', closeEmojiPickers);
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
  @apply max-w-full rounded-lg my-2 cursor-pointer hover:opacity-90 transition-opacity;
  max-height: 400px;
  object-fit: contain;
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

/* 灯箱过渡动画 */
.lightbox-enter-active,
.lightbox-leave-active {
  transition: opacity 0.3s ease;
}

.lightbox-enter-active img,
.lightbox-leave-active img {
  transition: transform 0.3s ease;
}

.lightbox-enter-from,
.lightbox-leave-to {
  opacity: 0;
}

.lightbox-enter-from img,
.lightbox-leave-to img {
  transform: scale(0.9);
}
</style>
