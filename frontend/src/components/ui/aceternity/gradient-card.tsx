import React from "react";
import { cn } from "../../../utils/cn";

interface GradientCardProps {
  className?: string;
  children: React.ReactNode;
  hoverEffect?: boolean;
}

export const GradientCard: React.FC<GradientCardProps> = ({
  className,
  children,
  hoverEffect = true,
}) => {
  return (
    <div
      className={cn(
        "aceternity-card relative overflow-hidden rounded-xl border p-6",
        hoverEffect &&
          "transition-all duration-300 hover:shadow-lg hover:-translate-y-1",
        className,
      )}
    >
      <div className="aceternity-glow" />
      <div className="relative z-10">{children}</div>
    </div>
  );
};

export default GradientCard;
