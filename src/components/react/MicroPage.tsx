import { useState, useCallback } from 'react';
import MicroList from './MicroList';
import MicroComposer from './MicroComposer';
import Heatmap from './Heatmap';

interface MicroPageProps {
  apiBaseUrl?: string;
}

const API_BASE = typeof window !== 'undefined' 
  ? (import.meta.env.VITE_API_BASE || 'http://localhost:8080/api')
  : 'http://localhost:8080/api';

export default function MicroPage({ apiBaseUrl }: MicroPageProps) {
  const baseUrl = apiBaseUrl || API_BASE;
  const [refreshKey, setRefreshKey] = useState(0);
  const [heatmapKey, setHeatmapKey] = useState(0);

  const handlePostSuccess = useCallback(() => {
    setRefreshKey(prev => prev + 1);
    setHeatmapKey(prev => prev + 1);
  }, []);

  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-8">
      <div className="xl:col-span-2">
        <div className="mb-8">
          <MicroComposer 
            apiBaseUrl={baseUrl} 
            onSuccess={handlePostSuccess}
          />
        </div>
        <MicroList 
          key={refreshKey}
          apiBaseUrl={baseUrl} 
        />
      </div>
      
      <div className="xl:col-span-1">
        <div className="sticky top-24 space-y-6">
          <div className="card">
            <Heatmap 
              key={heatmapKey}
              apiBaseUrl={baseUrl} 
            />
          </div>
          
          <div className="card">
            <h3 className="text-lg font-semibold text-foreground mb-4">关于微语</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              微语是一个轻量级的分享空间，让你可以随时记录生活中的点滴、灵感与感悟。
              支持文字、标签，所有注册用户都可以发布。
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
