import { useState, useEffect, useCallback } from 'react';

interface User {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
}

interface MicroPost {
  id: number;
  content: string;
  images: string;
  tags: string;
  is_public: boolean;
  created_at: string;
  updated_at: string;
  user: User;
  like_count: number;
  is_liked: boolean;
}

interface MicroListProps {
  userId?: string;
  onOpenComposer?: () => void;
  apiBaseUrl?: string;
}

const API_BASE = typeof window !== 'undefined' 
  ? (import.meta.env.VITE_API_BASE || 'http://localhost:8080/api')
  : 'http://localhost:8080/api';

export default function MicroList({ userId, onOpenComposer, apiBaseUrl }: MicroListProps) {
  const baseUrl = apiBaseUrl || API_BASE;
  const [micros, setMicros] = useState<MicroPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [total, setTotal] = useState(0);

  const fetchMicros = useCallback(async (pageNum: number, append = false) => {
    try {
      const params = new URLSearchParams({
        page: pageNum.toString(),
        page_size: '10',
      });
      if (userId) params.append('user_id', userId);

      const token = localStorage.getItem('token');
      const headers: HeadersInit = {};
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(`${baseUrl}/micros?${params}`, { headers });
      if (response.ok) {
        const result = await response.json();
        if (append) {
          setMicros(prev => [...prev, ...result.data]);
        } else {
          setMicros(result.data);
        }
        setTotal(result.pagination.total);
        setHasMore(result.pagination.page < result.pagination.total_page);
      }
    } catch (error) {
      console.error('Failed to fetch micros:', error);
    } finally {
      setLoading(false);
    }
  }, [userId, baseUrl]);

  useEffect(() => {
    fetchMicros(1);
  }, [fetchMicros]);

  const loadMore = () => {
    const nextPage = page + 1;
    setPage(nextPage);
    fetchMicros(nextPage, true);
  };

  const handleLike = async (microId: number) => {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('请登录后再点赞');
      return;
    }

    try {
      const response = await fetch(`${baseUrl}/micros/${microId}/like`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const result = await response.json();
        setMicros(prev => prev.map(m => {
          if (m.id === microId) {
            return {
              ...m,
              is_liked: result.liked,
              like_count: result.liked ? m.like_count + 1 : m.like_count - 1,
            };
          }
          return m;
        }));
      }
    } catch (error) {
      console.error('Failed to like:', error);
    }
  };

  const formatTime = (dateStr: string) => {
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
    return date.toLocaleDateString('zh-CN');
  };

  const parseJSON = (str: string) => {
    try {
      return JSON.parse(str || '[]');
    } catch {
      return [];
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-32">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-foreground">微语</h2>
        <div className="flex items-center gap-4">
          <span className="text-sm text-muted-foreground">共 {total} 条</span>
          {onOpenComposer && (
            <button
              onClick={onOpenComposer}
              className="btn-primary text-sm"
            >
              发布微语
            </button>
          )}
        </div>
      </div>

      {micros.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          暂无微语，来发布第一条吧！
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {micros.map(micro => {
              const images = parseJSON(micro.images);
              const tags = parseJSON(micro.tags);

              return (
                <div key={micro.id} className="card">
                  <div className="flex gap-3">
                    <div className="flex-shrink-0">
                      {micro.user.avatar ? (
                        <img
                          src={micro.user.avatar}
                          alt={micro.user.nickname || micro.user.username}
                          className="w-10 h-10 rounded-full object-cover"
                        />
                      ) : (
                        <div className="w-10 h-10 rounded-full bg-primary-500 flex items-center justify-center text-white font-semibold">
                          {(micro.user.nickname || micro.user.username).charAt(0).toUpperCase()}
                        </div>
                      )}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-semibold text-foreground">
                          {micro.user.nickname || micro.user.username}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          @{micro.user.username}
                        </span>
                        <span className="text-xs text-muted-foreground">·</span>
                        <span className="text-xs text-muted-foreground">
                          {formatTime(micro.created_at)}
                        </span>
                      </div>

                      <div className="text-foreground whitespace-pre-wrap break-words mb-3">
                        {micro.content}
                      </div>

                      {images.length > 0 && (
                        <div className={`grid gap-2 mb-3 ${
                          images.length === 1 ? 'grid-cols-1' :
                          images.length === 2 ? 'grid-cols-2' :
                          'grid-cols-3'
                        }`}>
                          {images.map((img: string, i: number) => (
                            <img
                              key={i}
                              src={img}
                              alt={`图片 ${i + 1}`}
                              className="rounded-lg object-cover w-full aspect-square"
                            />
                          ))}
                        </div>
                      )}

                      {tags.length > 0 && (
                        <div className="flex flex-wrap gap-2 mb-3">
                          {tags.map((tag: string, i: number) => (
                            <span key={i} className="tag text-xs">
                              #{tag}
                            </span>
                          ))}
                        </div>
                      )}

                      <div className="flex items-center gap-6 text-muted-foreground">
                        <button
                          onClick={() => handleLike(micro.id)}
                          className={`flex items-center gap-1 transition-colors ${
                            micro.is_liked ? 'text-red-500' : 'hover:text-red-500'
                          }`}
                        >
                          <svg
                            className="w-5 h-5"
                            fill={micro.is_liked ? 'currentColor' : 'none'}
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                            />
                          </svg>
                          <span className="text-sm">{micro.like_count || ''}</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

          {hasMore && (
            <div className="text-center pt-4">
              <button
                onClick={loadMore}
                className="btn-secondary"
              >
                加载更多
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
