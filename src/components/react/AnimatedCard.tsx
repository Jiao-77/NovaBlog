import { useState } from 'react';

interface AnimatedCardProps {
  title: string;
  description: string;
  color?: string;
}

export default function AnimatedCard({ 
  title, 
  description, 
  color = '#3b82f6' 
}: AnimatedCardProps) {
  const [isHovered, setIsHovered] = useState(false);

  const cardStyle: React.CSSProperties = {
    padding: '1.5rem',
    background: isHovered ? color : `${color}dd`,
    borderRadius: '1rem',
    color: 'white',
    cursor: 'pointer',
    transform: isHovered ? 'translateY(-8px) scale(1.02)' : 'translateY(0) scale(1)',
    boxShadow: isHovered 
      ? '0 20px 40px rgba(0, 0, 0, 0.3)' 
      : '0 4px 6px rgba(0, 0, 0, 0.1)',
    transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
  };

  return (
    <div 
      style={cardStyle}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <h3 style={{ 
        margin: '0 0 0.5rem 0', 
        fontSize: '1.25rem',
        fontWeight: 'bold' 
      }}>
        {title}
      </h3>
      <p style={{ margin: 0, opacity: 0.9 }}>
        {description}
      </p>
    </div>
  );
}