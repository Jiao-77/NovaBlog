import { useState, useEffect } from 'react';

interface TypewriterTextProps {
  text: string;
  speed?: number;
  loop?: boolean;
  style?: React.CSSProperties;
}

export default function TypewriterText({ 
  text, 
  speed = 100, 
  loop = false,
  style = {}
}: TypewriterTextProps) {
  const [displayedText, setDisplayedText] = useState('');
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    const timeout = setTimeout(() => {
      if (!isDeleting) {
        // 打字
        if (currentIndex < text.length) {
          setDisplayedText(text.slice(0, currentIndex + 1));
          setCurrentIndex(currentIndex + 1);
        } else if (loop) {
          // 打完后等待，然后开始删除
          setTimeout(() => setIsDeleting(true), 1500);
        }
      } else {
        // 删除
        if (currentIndex > 0) {
          setDisplayedText(text.slice(0, currentIndex - 1));
          setCurrentIndex(currentIndex - 1);
        } else {
          setIsDeleting(false);
        }
      }
    }, isDeleting ? speed / 2 : speed);

    return () => clearTimeout(timeout);
  }, [currentIndex, isDeleting, text, speed, loop]);

  const containerStyle: React.CSSProperties = {
    fontFamily: 'monospace',
    fontSize: '1.5rem',
    color: '#3b82f6',
    ...style
  };

  const cursorStyle: React.CSSProperties = {
    display: 'inline-block',
    width: '3px',
    height: '1.5rem',
    background: '#3b82f6',
    marginLeft: '2px',
    animation: 'blink 1s infinite',
    verticalAlign: 'middle',
  };

  return (
    <>
      <style>{`
        @keyframes blink {
          0%, 50% { opacity: 1; }
          51%, 100% { opacity: 0; }
        }
      `}</style>
      <span style={containerStyle}>
        {displayedText}
        <span style={cursorStyle} />
      </span>
    </>
  );
}