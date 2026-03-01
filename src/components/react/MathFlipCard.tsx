import { useState, useEffect, useRef } from 'react';
import katex from 'katex';

interface MathFlipCardProps {
  latex: string;
  displayMode?: boolean;
  className?: string;
}

export default function MathFlipCard({ 
  latex, 
  displayMode = true,
  className = ''
}: MathFlipCardProps) {
  const [isFlipped, setIsFlipped] = useState(false);
  const mathContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (mathContainerRef.current && !isFlipped) {
      try {
        katex.render(latex, mathContainerRef.current, {
          displayMode: displayMode,
          throwOnError: false,
          trust: true,
        });
      } catch (error) {
        console.error('KaTeX rendering error:', error);
        if (mathContainerRef.current) {
          mathContainerRef.current.textContent = latex;
        }
      }
    }
  }, [latex, displayMode, isFlipped]);

  const containerStyle: React.CSSProperties = {
    perspective: '1000px',
    width: '100%',
    minHeight: displayMode ? '120px' : '60px',
    position: 'relative',
  };

  const innerStyle: React.CSSProperties = {
    position: 'relative',
    width: '100%',
    minHeight: 'inherit',
    transformStyle: 'preserve-3d',
    transition: 'transform 0.6s cubic-bezier(0.4, 0, 0.2, 1)',
    transform: isFlipped ? 'rotateY(180deg)' : 'rotateY(0deg)',
  };

  const faceStyle = (isFront: boolean): React.CSSProperties => ({
    position: 'absolute',
    width: '100%',
    minHeight: 'inherit',
    backfaceVisibility: 'hidden',
    borderRadius: '0.75rem',
    padding: displayMode ? '1.5rem' : '0.75rem 1rem',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    background: 'var(--color-muted, #f1f5f9)',
    border: '1px solid var(--color-border, #e2e8f0)',
    transform: isFront ? 'rotateY(0deg)' : 'rotateY(180deg)',
    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
  });

  const codeStyle: React.CSSProperties = {
    fontFamily: 'var(--font-mono, "JetBrains Mono", monospace)',
    fontSize: displayMode ? '0.9rem' : '0.85rem',
    color: 'var(--color-foreground, #1e293b)',
    background: 'transparent',
    wordBreak: 'break-word',
    textAlign: 'center',
    width: '100%',
  };

  const buttonStyle: React.CSSProperties = {
    position: 'absolute',
    bottom: '0.5rem',
    right: '0.5rem',
    padding: '0.25rem 0.75rem',
    fontSize: '0.75rem',
    fontWeight: 500,
    color: 'var(--primary-500, #0ea5e9)',
    background: 'transparent',
    border: '1px solid var(--primary-500, #0ea5e9)',
    borderRadius: '0.375rem',
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    zIndex: 10,
  };

  const handleButtonClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsFlipped(!isFlipped);
  };

  return (
    <div style={containerStyle} className={className}>
      <div style={innerStyle}>
        <div style={faceStyle(true)}>
          <div 
            ref={mathContainerRef}
            style={{
              color: 'var(--color-foreground, #1e293b)',
              width: '100%',
              display: 'flex',
              justifyContent: 'center',
            }}
          />
          <button
            style={buttonStyle}
            onClick={handleButtonClick}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'var(--primary-500, #0ea5e9)';
              e.currentTarget.style.color = 'white';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'transparent';
              e.currentTarget.style.color = 'var(--primary-500, #0ea5e9)';
            }}
          >
            显示 LaTeX
          </button>
        </div>
        <div style={faceStyle(false)}>
          <pre style={codeStyle}>{latex}</pre>
          <button
            style={buttonStyle}
            onClick={handleButtonClick}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'var(--primary-500, #0ea5e9)';
              e.currentTarget.style.color = 'white';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'transparent';
              e.currentTarget.style.color = 'var(--primary-500, #0ea5e9)';
            }}
          >
            显示公式
          </button>
        </div>
      </div>
    </div>
  );
}
