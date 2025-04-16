import React, { useState, useRef, useEffect } from 'react';
import { cn } from '../../../utils/cn';

interface ThreeDCardProps {
  className?: string;
  children: React.ReactNode;
  glareEnabled?: boolean;
  rotationIntensity?: number;
  glareIntensity?: number;
  borderRadius?: number;
}

export const ThreeDCard: React.FC<ThreeDCardProps> = ({
  className,
  children,
  glareEnabled = true,
  rotationIntensity = 20,
  glareIntensity = 0.5,
  borderRadius = 20,
}) => {
  const cardRef = useRef<HTMLDivElement>(null);
  const [rotation, setRotation] = useState({ x: 0, y: 0 });
  const [glarePosition, setGlarePosition] = useState({ x: 0, y: 0 });
  const [isHovered, setIsHovered] = useState(false);

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current) return;
    
    const card = cardRef.current;
    const rect = card.getBoundingClientRect();
    
    const x = ((e.clientX - rect.left) / rect.width - 0.5) * 2;
    const y = ((e.clientY - rect.top) / rect.height - 0.5) * 2;
    
    setRotation({
      x: -y * rotationIntensity,
      y: x * rotationIntensity,
    });
    
    setGlarePosition({
      x: (e.clientX - rect.left) / rect.width * 100,
      y: (e.clientY - rect.top) / rect.height * 100,
    });
  };

  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
    setRotation({ x: 0, y: 0 });
    setGlarePosition({ x: 50, y: 50 });
  };

  return (
    <div
      ref={cardRef}
      className={cn(
        'relative overflow-hidden transition-all duration-200',
        className
      )}
      style={{
        borderRadius: `${borderRadius}px`,
        transform: isHovered
          ? `perspective(1000px) rotateX(${rotation.x}deg) rotateY(${rotation.y}deg) scale3d(1.05, 1.05, 1.05)`
          : 'perspective(1000px) rotateX(0) rotateY(0) scale3d(1, 1, 1)',
        transition: 'transform 0.2s ease',
      }}
      onMouseMove={handleMouseMove}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      {glareEnabled && isHovered && (
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            background: `radial-gradient(circle at ${glarePosition.x}% ${glarePosition.y}%, rgba(255, 255, 255, ${glareIntensity}), transparent)`,
            mixBlendMode: 'overlay',
          }}
        />
      )}
      <div
        className="relative z-10 transform-style-3d"
        style={{
          transform: isHovered ? 'translateZ(50px)' : 'translateZ(0)',
          transition: 'transform 0.2s ease',
        }}
      >
        {children}
      </div>
    </div>
  );
};

export default ThreeDCard;
