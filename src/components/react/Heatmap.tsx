import { useState, useEffect } from 'react';

interface HeatmapData {
  date: string;
  count: number;
}

interface HeatmapProps {
  userId?: string;
  year?: number;
  apiBaseUrl?: string;
}

const API_BASE = typeof window !== 'undefined' 
  ? (import.meta.env.VITE_API_BASE || 'http://localhost:8080/api')
  : 'http://localhost:8080/api';

const COLORS = [
  'bg-primary-100 dark:bg-primary-900/30',
  'bg-primary-200 dark:bg-primary-800/40',
  'bg-primary-300 dark:bg-primary-700/50',
  'bg-primary-400 dark:bg-primary-600/60',
  'bg-primary-500 dark:bg-primary-500/70',
];

const MONTHS = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'];
const DAYS = ['日', '一', '二', '三', '四', '五', '六'];

export default function Heatmap({ userId, year = new Date().getFullYear(), apiBaseUrl }: HeatmapProps) {
  const baseUrl = apiBaseUrl || API_BASE;
  const [data, setData] = useState<HeatmapData[]>([]);
  const [loading, setLoading] = useState(true);
  const [hoveredCell, setHoveredCell] = useState<{ date: string; count: number } | null>(null);

  useEffect(() => {
    const fetchHeatmap = async () => {
      try {
        const params = new URLSearchParams({ year: year.toString() });
        if (userId) params.append('user_id', userId);
        
        const response = await fetch(`${baseUrl}/micros/heatmap?${params}`);
        if (response.ok) {
          const result = await response.json();
          setData(result);
        }
      } catch (error) {
        console.error('Failed to fetch heatmap:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchHeatmap();
  }, [userId, year, baseUrl]);

  const getDaysInYear = (year: number) => {
    const days: Date[] = [];
    const startDate = new Date(year, 0, 1);
    const endDate = new Date(year, 11, 31);
    
    for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
      days.push(new Date(d));
    }
    return days;
  };

  const getCountForDate = (date: Date): number => {
    const dateStr = date.toISOString().split('T')[0];
    const item = data.find(d => d.date === dateStr);
    return item ? item.count : 0;
  };

  const getColorClass = (count: number): string => {
    if (count === 0) return 'bg-muted dark:bg-muted/50';
    if (count <= 2) return COLORS[0];
    if (count <= 4) return COLORS[1];
    if (count <= 6) return COLORS[2];
    if (count <= 8) return COLORS[3];
    return COLORS[4];
  };

  const formatDisplayDate = (date: Date): string => {
    return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`;
  };

  const getWeeksInYear = (year: number) => {
    const days = getDaysInYear(year);
    const weeks: Date[][] = [];
    let currentWeek: Date[] = [];

    const firstDay = days[0];
    const firstDayOfWeek = firstDay.getDay();
    for (let i = 0; i < firstDayOfWeek; i++) {
      currentWeek.push(new Date(year, 0, 1 - firstDayOfWeek + i));
    }

    days.forEach(day => {
      if (day.getDay() === 0 && currentWeek.length > 0) {
        weeks.push(currentWeek);
        currentWeek = [];
      }
      currentWeek.push(day);
    });

    if (currentWeek.length > 0) {
      weeks.push(currentWeek);
    }

    return weeks;
  };

  const weeks = getWeeksInYear(year);
  const totalCount = data.reduce((sum, item) => sum + item.count, 0);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <h3 className="text-lg font-semibold text-foreground">{year} 年微语热力图</h3>
          <span className="text-sm text-muted-foreground">共 {totalCount} 条</span>
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <span>少</span>
          {COLORS.map((color, i) => (
            <div key={i} className={`w-3 h-3 rounded-sm ${color}`}></div>
          ))}
          <span>多</span>
        </div>
      </div>

      <div className="overflow-x-auto">
        <div className="inline-flex gap-1">
          <div className="flex flex-col gap-1 mr-2 text-xs text-muted-foreground">
            {DAYS.map((day, i) => (
              <div key={i} className="h-3 flex items-center">{i % 2 === 1 ? day : ''}</div>
            ))}
          </div>
          
          {weeks.map((week, weekIndex) => (
            <div key={weekIndex} className="flex flex-col gap-1">
              {DAYS.map((_, dayIndex) => {
                const day = week[dayIndex];
                if (!day || day.getFullYear() !== year) {
                  return <div key={dayIndex} className="w-3 h-3"></div>;
                }

                const count = getCountForDate(day);
                const colorClass = getColorClass(count);

                return (
                  <div
                    key={dayIndex}
                    className={`w-3 h-3 rounded-sm ${colorClass} cursor-pointer transition-transform hover:scale-125`}
                    onMouseEnter={() => setHoveredCell({ date: formatDisplayDate(day), count })}
                    onMouseLeave={() => setHoveredCell(null)}
                  />
                );
              })}
            </div>
          ))}
        </div>
      </div>

      {hoveredCell && (
        <div className="text-sm text-muted-foreground">
          {hoveredCell.date}：{hoveredCell.count} 条微语
        </div>
      )}

      <div className="flex gap-4 text-xs text-muted-foreground">
        {MONTHS.map((month, i) => (
          <span key={i} className="flex-1 text-center">{month}</span>
        ))}
      </div>
    </div>
  );
}
