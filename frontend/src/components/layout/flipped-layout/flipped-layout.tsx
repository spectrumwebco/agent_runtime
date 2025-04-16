import React, { useState } from 'react';
import { cn } from '../../../utils/cn';

interface FlippedLayoutProps {
  className?: string;
  children?: React.ReactNode;
  leftPanel?: React.ReactNode;
  rightPanel?: React.ReactNode;
  initialLeftWidth?: number;
}

export const FlippedLayout: React.FC<FlippedLayoutProps> = ({
  className,
  children,
  leftPanel,
  rightPanel,
  initialLeftWidth = 60,
}) => {
  const [leftWidth, setLeftWidth] = useState(initialLeftWidth);
  const [isDragging, setIsDragging] = useState(false);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging) return;
    
    const container = e.currentTarget.getBoundingClientRect();
    const newLeftWidth = ((e.clientX - container.left) / container.width) * 100;
    
    if (newLeftWidth >= 20 && newLeftWidth <= 80) {
      setLeftWidth(newLeftWidth);
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  return (
    <div 
      className={cn('flex h-full w-full overflow-hidden', className)}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
    >
      {/* Left Panel (IDE/Code Editor) */}
      <div 
        className="h-full overflow-auto"
        style={{ width: `${leftWidth}%` }}
      >
        {leftPanel}
      </div>
      
      {/* Resizable Divider */}
      <div 
        className={cn(
          'w-1 h-full bg-gray-300 dark:bg-gray-700 hover:bg-emerald-500 dark:hover:bg-emerald-500 cursor-col-resize transition-colors',
          isDragging && 'bg-emerald-500 dark:bg-emerald-500'
        )}
        onMouseDown={handleMouseDown}
      />
      
      {/* Right Panel (Menu/Controls) */}
      <div 
        className="h-full overflow-auto"
        style={{ width: `${100 - leftWidth - 0.25}%` }}
      >
        {rightPanel}
      </div>
      
      {children}
    </div>
  );
};

export default FlippedLayout;
