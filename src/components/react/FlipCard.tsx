import { useState } from 'react';

interface FlipCardProps {
  frontTitle: string;
  frontDescription: string;
  backTitle: string;
  backDescription: string;
  frontColor?: string;
  backColor?: string;
}

export default function FlipCard({ 
  frontTitle,
  frontDescription,
  backTitle,
  backDescription,
  frontColor = '#3b82f6',
  backColor = '#10b981'
}: FlipCardProps) {
  const [isFlipped, setIsFlipped] = useState(false);

  const containerStyle: React.CSSProperties = {
    perspective: '1000px',
    width: '100%',
    height: '200px',
    cursor: 'pointer',
  };

  const innerStyle: React.CSSProperties = {
    position: 'relative',
    width: '100%',
    height: '100%',
    transformStyle: 'preserve-3d',
    transition: 'transform 0.6s cubic-bezier(0.4, 0, 0.2, 1)',
    transform: isFlipped ? 'rotateY(180deg)' : 'rotateY(0deg)',
  };

  const faceStyle = (color: string, isFront: boolean): React.CSSProperties => ({
    position: 'absolute',
    width: '100%',
    height: '100%',
    backfaceVisibility: 'hidden',
    borderRadius: '1rem',
    padding: '1.5rem',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    background: color,
    color: 'white',
    transform: isFront ? 'rotateY(0deg)' : 'rotateY(180deg)',
    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
  });

  return (
    <div 
      style={containerStyle}
      onClick={() => setIsFlipped(!isFlipped)}
      onMouseEnter={(e) => e.currentTarget.style.transform = 'scale(1.02)'}
      onMouseLeave={(e) => e.currentTarget.style.transform = 'scale(1)'}
    >
      <div style={innerStyle}>
        <div style={faceStyle(frontColor, true)}>
          <div style={{ textAlign: 'center' }}>
            <h3 style={{ margin: '0 0 0.5rem 0', fontSize: '1.25rem', fontWeight: 'bold' }}>{frontTitle}</h3>
            <p style={{ margin: 0, opacity: 0.9 }}>{frontDescription}</p>
          </div>
        </div>
        <div style={faceStyle(backColor, false)}>
          <div style={{ textAlign: 'center' }}>
            <h3 style={{ margin: '0 0 0.5rem 0', fontSize: '1.25rem', fontWeight: 'bold' }}>{backTitle}</h3>
            <p style={{ margin: 0, opacity: 0.9 }}>{backDescription}</p>
          </div>
        </div>
      </div>
    </div>
  );
}