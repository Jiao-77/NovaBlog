import { useEffect, useRef } from 'react';

interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
  size: number;
  opacity: number;
}

interface ParticleBackgroundProps {
  particleCount?: number;
  color?: string;
  speed?: number;
  children?: React.ReactNode;
}

export default function ParticleBackground({ 
  particleCount = 50,
  color = '#3b82f6',
  speed = 1,
  children 
}: ParticleBackgroundProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const particlesRef = useRef<Particle[]>([]);
  const animationRef = useRef<number | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // 设置画布大小
    const resizeCanvas = () => {
      const rect = canvas.parentElement?.getBoundingClientRect();
      if (rect) {
        canvas.width = rect.width;
        canvas.height = rect.height;
      }
    };

    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);

    // 初始化粒子
    const initParticles = () => {
      const newParticles: Particle[] = [];
      for (let i = 0; i < particleCount; i++) {
        newParticles.push({
          x: Math.random() * canvas.width,
          y: Math.random() * canvas.height,
          vx: (Math.random() - 0.5) * speed,
          vy: (Math.random() - 0.5) * speed,
          size: Math.random() * 3 + 1,
          opacity: Math.random() * 0.5 + 0.2,
        });
      }
      particlesRef.current = newParticles;
    };

    initParticles();

    // 动画循环
    const animate = () => {
      if (!ctx || !canvas) return;
      
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      const particles = particlesRef.current;

      particles.forEach((particle: Particle, i: number) => {
        // 更新位置
        particle.x += particle.vx;
        particle.y += particle.vy;

        // 边界检测
        if (particle.x < 0 || particle.x > canvas.width) particle.vx *= -1;
        if (particle.y < 0 || particle.y > canvas.height) particle.vy *= -1;

        // 绘制粒子
        ctx.beginPath();
        ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
        ctx.fillStyle = `${color}${Math.floor(particle.opacity * 255).toString(16).padStart(2, '0')}`;
        ctx.fill();

        // 绘制连线
        particles.slice(i + 1).forEach((other: Particle) => {
          const dx = particle.x - other.x;
          const dy = particle.y - other.y;
          const distance = Math.sqrt(dx * dx + dy * dy);

          if (distance < 150) {
            ctx.beginPath();
            ctx.moveTo(particle.x, particle.y);
            ctx.lineTo(other.x, other.y);
            const opacity = (1 - distance / 150) * 0.2;
            ctx.strokeStyle = `${color}${Math.floor(opacity * 255).toString(16).padStart(2, '0')}`;
            ctx.stroke();
          }
        });
      });

      animationRef.current = requestAnimationFrame(animate);
    };

    animate();

    return () => {
      window.removeEventListener('resize', resizeCanvas);
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [color, speed, particleCount]);

  const containerStyle: React.CSSProperties = {
    position: 'relative',
    width: '100%',
    height: '300px',
    background: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)',
    borderRadius: '1rem',
    overflow: 'hidden',
  };

  const canvasStyle: React.CSSProperties = {
    position: 'absolute',
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
  };

  const contentStyle: React.CSSProperties = {
    position: 'relative',
    zIndex: 1,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    height: '100%',
    color: 'white',
  };

  return (
    <div style={containerStyle}>
      <canvas ref={canvasRef} style={canvasStyle} />
      <div style={contentStyle}>
        {children || (
          <div style={{ textAlign: 'center' }}>
            <h3 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>✨ 粒子动效</h3>
            <p style={{ opacity: 0.8 }}>鼠标悬停查看效果</p>
          </div>
        )}
      </div>
    </div>
  );
}