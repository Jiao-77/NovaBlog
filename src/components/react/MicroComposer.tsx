import { useState } from 'react';

interface MicroComposerProps {
  onClose?: () => void;
  onSuccess?: () => void;
  apiBaseUrl?: string;
}

const API_BASE = typeof window !== 'undefined' 
  ? (import.meta.env.VITE_API_BASE || 'http://localhost:8080/api')
  : 'http://localhost:8080/api';

export default function MicroComposer({ onClose, onSuccess, apiBaseUrl }: MicroComposerProps) {
  const baseUrl = apiBaseUrl || API_BASE;
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');
  const [isPublic, setIsPublic] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!content.trim()) {
      alert('请输入内容');
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      alert('请登录后再发布');
      return;
    }

    setSubmitting(true);
    try {
      const tagList = tags
        .split(/[,，\s]+/)
        .map(t => t.trim())
        .filter(t => t.length > 0);

      const response = await fetch(`${baseUrl}/micros`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          content: content.trim(),
          tags: tagList,
          is_public: isPublic,
          images: [],
        }),
      });

      if (response.ok) {
        setContent('');
        setTags('');
        setIsPublic(true);
        onSuccess?.();
        onClose?.();
      } else {
        const error = await response.json();
        alert(error.error || '发布失败');
      }
    } catch (error) {
      console.error('Failed to post:', error);
      alert('发布失败，请重试');
    } finally {
      setSubmitting(false);
    }
  };

  const remainingChars = 2000 - content.length;

  return (
    <div className="card">
      <div className="mb-4">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="分享你的想法..."
          className="input min-h-[120px] resize-none"
          maxLength={2000}
        />
        <div className="flex justify-end mt-1">
          <span className={`text-xs ${remainingChars < 100 ? 'text-red-500' : 'text-muted-foreground'}`}>
            {remainingChars} 字
          </span>
        </div>
      </div>

      <div className="mb-4">
        <input
          type="text"
          value={tags}
          onChange={(e) => setTags(e.target.value)}
          placeholder="标签（用逗号或空格分隔）"
          className="input"
        />
      </div>

      <div className="flex items-center justify-between">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={isPublic}
            onChange={(e) => setIsPublic(e.target.checked)}
            className="w-4 h-4 rounded border-border text-primary-500 focus:ring-primary-500"
          />
          <span className="text-sm text-muted-foreground">公开可见</span>
        </label>

        <div className="flex gap-2">
          {onClose && (
            <button
              onClick={onClose}
              className="btn-secondary"
              disabled={submitting}
            >
              取消
            </button>
          )}
          <button
            onClick={handleSubmit}
            className="btn-primary"
            disabled={submitting || !content.trim()}
          >
            {submitting ? '发布中...' : '发布'}
          </button>
        </div>
      </div>
    </div>
  );
}
