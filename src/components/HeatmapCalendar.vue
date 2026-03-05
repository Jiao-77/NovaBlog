<template>
  <div class="heatmap-calendar">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold">{{ title }}</h3>
      <div class="flex items-center gap-2 text-sm text-foreground/60">
        <span>{{ totalContributions }} 条微语</span>
        <span class="text-foreground/30">|</span>
        <span>{{ activeDays }} 天活跃</span>
      </div>
    </div>

    <!-- 热力图网格 -->
    <div class="overflow-x-auto hide-scrollbar">
      <div class="flex gap-1" style="min-width: fit-content;">
        <div v-for="(week, weekIndex) in weeks" :key="weekIndex" class="flex flex-col gap-1">
          <div
            v-for="(day, dayIndex) in week"
            :key="dayIndex"
            class="w-3 h-3 rounded-sm cursor-pointer transition-all duration-200 hover:ring-2 hover:ring-primary-400"
            :class="getLevelClass(day.level)"
            :style="{ backgroundColor: getLevelColor(day.level) }"
            @mouseenter="showTooltip($event, day)"
            @mouseleave="hideTooltip"
          ></div>
        </div>
      </div>
    </div>

    <!-- 图例 -->
    <div class="flex items-center justify-end gap-2 mt-3 text-xs text-foreground/40">
      <span>少</span>
      <div class="flex gap-1">
        <div class="w-3 h-3 rounded-sm" :style="{ backgroundColor: getLevelColor(0) }"></div>
        <div class="w-3 h-3 rounded-sm" :style="{ backgroundColor: getLevelColor(1) }"></div>
        <div class="w-3 h-3 rounded-sm" :style="{ backgroundColor: getLevelColor(2) }"></div>
        <div class="w-3 h-3 rounded-sm" :style="{ backgroundColor: getLevelColor(3) }"></div>
        <div class="w-3 h-3 rounded-sm" :style="{ backgroundColor: getLevelColor(4) }"></div>
      </div>
      <span>多</span>
    </div>

    <!-- Tooltip -->
    <Teleport to="body">
      <Transition name="tooltip">
        <div
          v-if="tooltipVisible"
          class="fixed z-[9999] px-3 py-2 text-sm bg-gray-900 dark:bg-gray-700 text-white rounded-lg shadow-lg pointer-events-none"
          :style="{ left: tooltipX + 'px', top: tooltipY + 'px' }"
        >
          <div class="font-medium">{{ tooltipDate }}</div>
          <div class="text-gray-300">{{ tooltipCount }} 条微语</div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

interface DayData {
  date: Date;
  count: number;
  level: number;
}

interface Props {
  title?: string;
  data?: Record<string, number>;
  year?: number;
  colorScheme?: 'green' | 'blue' | 'purple' | 'orange';
}

const props = withDefaults(defineProps<Props>(), {
  title: '活动热力图',
  year: () => new Date().getFullYear(),
  colorScheme: 'green',
});

// 颜色方案
const colorSchemes = {
  green: ['#ebedf0', '#9be9a8', '#40c463', '#30a14e', '#216e39'],
  blue: ['#ebedf0', '#a5d6ff', '#79c0ff', '#58a6ff', '#1f6feb'],
  purple: ['#ebedf0', '#d2b4fe', '#c084fc', '#a855f7', '#7c3aed'],
  orange: ['#ebedf0', '#fed7aa', '#fdba74', '#fb923c', '#ea580c'],
};

// 暗黑模式颜色方案
const darkColorSchemes = {
  green: ['#161b22', '#0e4429', '#006d32', '#26a641', '#39d353'],
  blue: ['#161b22', '#0a3069', '#0550ae', '#0969da', '#1f6feb'],
  purple: ['#161b22', '#3b0764', '#6b21a8', '#9333ea', '#a855f7'],
  orange: ['#161b22', '#431407', '#7c2d12', '#c2410c', '#ea580c'],
};

const isDark = ref(false);

// 检查暗黑模式
onMounted(() => {
  isDark.value = document.documentElement.classList.contains('dark');
  
  // 监听主题变化
  const observer = new MutationObserver(() => {
    isDark.value = document.documentElement.classList.contains('dark');
  });
  
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class'],
  });
});

// 获取颜色
function getLevelColor(level: number): string {
  const schemes = isDark.value ? darkColorSchemes : colorSchemes;
  return schemes[props.colorScheme][level];
}

// 获取级别样式类
function getLevelClass(level: number): string {
  return '';
}

// 生成一年的日期数据
const weeks = computed<DayData[][]>(() => {
  const year = props.year;
  const startDate = new Date(year, 0, 1);
  const endDate = new Date(year, 11, 31);
  
  // 调整到周日开始
  const firstDay = startDate.getDay();
  const firstSunday = new Date(startDate);
  firstSunday.setDate(startDate.getDate() - firstDay);
  
  const weeks: DayData[][] = [];
  let currentWeek: DayData[] = [];
  let currentDate = new Date(firstSunday);
  
  while (currentDate <= endDate || currentWeek.length > 0) {
    const dateStr = formatDateKey(currentDate);
    const count = props.data?.[dateStr] || 0;
    
    currentWeek.push({
      date: new Date(currentDate),
      count,
      level: getLevel(count),
    });
    
    if (currentWeek.length === 7) {
      weeks.push(currentWeek);
      currentWeek = [];
    }
    
    currentDate.setDate(currentDate.getDate() + 1);
    
    // 防止无限循环
    if (weeks.length > 54) break;
  }
  
  return weeks;
});

// 总贡献数
const totalContributions = computed(() => {
  return Object.values(props.data || {}).reduce((sum, count) => sum + count, 0);
});

// 活跃天数
const activeDays = computed(() => {
  return Object.values(props.data || {}).filter(count => count > 0).length;
});

// 根据数量获取级别
function getLevel(count: number): number {
  if (count === 0) return 0;
  if (count <= 2) return 1;
  if (count <= 4) return 2;
  if (count <= 6) return 3;
  return 4;
}

// 格式化日期为 key
function formatDateKey(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

// Tooltip 相关
const tooltipVisible = ref(false);
const tooltipX = ref(0);
const tooltipY = ref(0);
const tooltipDate = ref('');
const tooltipCount = ref(0);

function showTooltip(event: MouseEvent, day: DayData) {
  const rect = (event.target as HTMLElement).getBoundingClientRect();
  tooltipX.value = rect.left + rect.width / 2;
  tooltipY.value = rect.top - 50;
  
  tooltipDate.value = day.date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
  tooltipCount.value = day.count;
  tooltipVisible.value = true;
}

function hideTooltip() {
  tooltipVisible.value = false;
}
</script>

<style scoped>
.heatmap-calendar {
  @apply p-4 bg-background border border-border rounded-xl;
}

.tooltip-enter-active,
.tooltip-leave-active {
  transition: opacity 0.15s ease;
}

.tooltip-enter-from,
.tooltip-leave-to {
  opacity: 0;
}
</style>