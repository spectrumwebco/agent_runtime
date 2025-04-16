import React from "react";
import { cn } from "../../../utils/cn";

interface AnimatedBackgroundProps {
  className?: string;
  children?: React.ReactNode;
  variant?: "dots" | "grid" | "waves";
  color?: "emerald" | "blue" | "purple" | "amber";
  animate?: boolean;
}

export const AnimatedBackground: React.FC<AnimatedBackgroundProps> = ({
  className,
  children,
  variant = "dots",
  color = "emerald",
  animate = true,
}) => {
  const colorClasses = {
    emerald: "from-emerald-500/20 via-transparent to-emerald-500/20",
    blue: "from-blue-500/20 via-transparent to-blue-500/20",
    purple: "from-purple-500/20 via-transparent to-purple-500/20",
    amber: "from-amber-500/20 via-transparent to-amber-500/20",
  };

  const variantClasses = {
    dots: "bg-dots",
    grid: "bg-grid",
    waves: "bg-waves",
  };

  return (
    <div className={cn("relative w-full h-full overflow-hidden", className)}>
      <div
        className={cn(
          "absolute inset-0 z-0",
          variantClasses[variant],
          animate && "animate-pulse",
          colorClasses[color],
        )}
      />
      <div className="relative z-10 h-full">{children}</div>
    </div>
  );
};

export default AnimatedBackground;
