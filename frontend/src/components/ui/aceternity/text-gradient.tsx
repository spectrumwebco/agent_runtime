import React from "react";
import { cn } from "../../../utils/cn";

interface TextGradientProps {
  className?: string;
  children: React.ReactNode;
}

export const TextGradient: React.FC<TextGradientProps> = ({
  className,
  children,
}) => {
  return (
    <span className={cn("aceternity-text-gradient", className)}>
      {children}
    </span>
  );
};

export default TextGradient;
