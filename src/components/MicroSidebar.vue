<template>
  <div class="space-y-6">
    <!-- 热力图 -->
    <HeatmapCalendar
      title="发布活动"
      :data="heatmapData"
      :year="year"
      color-scheme="green"
    />

    <!-- 统计卡片 -->
    <div class="card">
      <h3 class="text-lg font-semibold mb-4">统计</h3>
      <div v-if="loading" class="text-center py-4">
        <div class="animate-spin w-6 h-6 border-2 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
      </div>
      <div v-else class="grid grid-cols-2 gap-4">
        <div class="text-center p-4 bg-muted/50 rounded-lg">
          <div class="text-2xl font-bold text-primary-500">{{ stats.total_micros || 0 }}</div>
          <div class="text-sm text-foreground/60">总微语</div>
        </div>
        <div class="text-center p-4 bg-muted/50 rounded-lg">
          <div class="text-2xl font-bold text-primary-500">{{ stats.month_micros || 0 }}</div>
          <div class="text-sm text-foreground/60">本月发布</div>
        </div>
        <div class="text-center p-4 bg-muted/50 rounded-lg">
          <div class="text-2xl font-bold text-primary-500">{{ formatNumber(stats.total_likes) }}</div>
          <div class="text-sm text-foreground/60">总点赞</div>
        </div>
        <div class="text-center p-4 bg-muted/50 rounded-lg">
          <div class="text-2xl font-bold text-primary-500">{{ stats.total_comments || 0 }}</div>
          <div class="text-sm text-foreground/60">总评论</div>
        </div>
      </div>
    </div>

    <!-- 标签云 -->
    <div class="card">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold">热门标签</h3>
        <button
          v-if="currentTag"
          @click="clearTagFilter"
          class="text-xs text-primary-500 hover:text-primary-600"
        >
          清除筛选
        </button>
      </div>
      <div v-if="loadingTags" class="text-center py-2">
        <div class="animate-spin w-5 h-5 border-2 border-primary-500 border-t-transparent rounded-full mx-auto"></div>
      </div>
      <div v-else-if="tags.length === 0" class="text-center py-2 text-foreground/40 text-sm">
        暂无标签
      </div>
      <div v-else class="flex flex-wrap gap-2">
        <span
          v-for="tag in tags"
          :key="tag.name"
          class="inline-block px-2 py-1 text-xs rounded-full cursor-pointer transition-colors"
          :class="currentTag === tag.name 
            ? 'bg-primary-500 text-white' 
            : 'bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 hover:bg-primary-200 dark:hover:bg-primary-800'"
          @click="filterByTag(tag.name)"
        >
          #{{ tag.name }} ({{ tag.count }})
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import HeatmapCalendar from './HeatmapCalendar.vue';

interface Props {
  apiBaseUrl?: string;
}

const props = withDefaults(defineProps<Props>(), {
  apiBaseUrl: 'http://localhost:8080/api',
});

const year = new Date().getFullYear();
const loading = ref(true);
const loadingTags = ref(true);
const heatmapData = ref<Record<string, number>>({});
const stats = ref({
  total_micros: 0,
  month_micros: 0,
  total_likes: 0,
  total_comments: 0,
});
const tags = ref<{ name: string; count: number }[]>([]);
const currentTag = ref('');

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k';
  }
  return String(num);
}

// 加载热力图数据
async function loadHeatmap() {
  try {
    const response = await fetch(`${props.apiBaseUrl}/micro/heatmap?year=${year}`);
    if (response.ok) {
      const data = await response.json();
      heatmapData.value = data.data || {};
    }
  } catch (error) {
    console.error('Failed to load heatmap:', error);
  }
}

// 加载统计数据
async function loadStats() {
  try {
    const response = await fetch(`${props.apiBaseUrl}/micro/stats`);
    if (response.ok) {
      const data = await response.json();
      stats.value = data;
    }
  } catch (error) {
    console.error('Failed to load stats:', error);
  } finally {
    loading.value = false;
  }
}

// 加载标签
async function loadTags() {
  try {
    const response = await fetch(`${props.apiBaseUrl}/micro/tags`);
    if (response.ok) {
      const data = await response.json();
      tags.value = data.tags || [];
    }
  } catch (error) {
    console.error('Failed to load tags:', error);
  } finally {
    loadingTags.value = false;
  }
}

// 按标签筛选
function filterByTag(tag: string) {
  currentTag.value = tag;
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('tag-filter', { detail: tag }));
  }
}

// 清除筛选
function clearTagFilter() {
  currentTag.value = '';
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent('tag-filter', { detail: '' }));
  }
}

// 刷新所有数据
function refresh() {
  loadHeatmap();
  loadStats();
  loadTags();
}

// 暴露方法给父组件
defineExpose({ refresh });

onMounted(() => {
  loadHeatmap();
  loadStats();
  loadTags();

  // 监听刷新事件
  window.addEventListener('refresh-sidebar', refresh);
  
  // 监听标签筛选事件（同步当前选中状态）
  window.addEventListener('tag-filter', ((e: CustomEvent) => {
    currentTag.value = e.detail || '';
  }) as EventListener);
});

onUnmounted(() => {
  window.removeEventListener('refresh-sidebar', refresh);
});
</script>