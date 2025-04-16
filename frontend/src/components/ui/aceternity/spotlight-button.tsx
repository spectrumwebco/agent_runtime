import React, { useState, useRef, useEffect } from "react";
import { cn } from "../../../utils/cn";

interface SpotlightButtonProps {
  className?: string;
  children: React.ReactNode;
  spotlightColor?: string;
  variant?: "default" | "outline" | "ghost";
  size?: "sm" | "md" | "lg";
  onClick?: () => void;
}

export const SpotlightButton: React.FC<SpotlightButtonProps> = ({
  className,
  children,
  spotlightColor = "rgba(80, 230, 180, 0.2)",
  variant = "default",
  size = "md",
  onClick,
}) => {
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isHovering, setIsHovering] = useState(false);

  const variantClasses = {
    default: "bg-black text-white border border-emerald-500/20",
    outline: "bg-transparent border border-emerald-500 text-emerald-500",
    ghost: "bg-transparent text-emerald-500 hover:bg-emerald-500/10",
  };

  const sizeClasses = {
    sm: "text-sm px-3 py-1",
    md: "px-4 py-2",
    lg: "text-lg px-6 py-3",
  };

  const handleMouseMove = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (!buttonRef.current) return;
    
    const rect = buttonRef.current.getBoundingClientRect();
    setPosition({
      x: e.clientX - rect.left,
      y: e.clientY - rect.top,
    });
  };

  return (
    <button
      ref={buttonRef}
      className={cn(
        "relative overflow-hidden rounded-lg font-medium transition-all duration-200",
        "hover:shadow-md active:scale-95",
        variantClasses[variant],
        sizeClasses[size],
        className
      )}
      onMouseEnter={() => setIsHovering(true)}
      onMouseLeave={() => setIsHovering(false)}
      onMouseMove={handleMouseMove}
      onClick={onClick}
    >
      {isHovering && (
        <div
          className="absolute pointer-events-none -inset-px opacity-0 transition duration-300 group-hover:opacity-100"
          style={{
            background: `radial-gradient(600px circle at ${position.x}px ${position.y}px, ${spotlightColor}, transparent 40%)`,
          }}
        />
      )}
      <span className="relative z-10">{children}</span>
    </button>
  );
};

export default SpotlightButton;
