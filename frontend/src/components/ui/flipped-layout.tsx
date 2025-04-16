import React from "react";
import { cn } from "../../utils/cn";

interface FlippedLayoutProps {
  editor: React.ReactNode;
  controls: React.ReactNode;
  className?: string;
}

export const FlippedLayout: React.FC<FlippedLayoutProps> = ({
  editor,
  controls,
  className,
}) => {
  return (
    <div className={cn("flex h-full w-full", className)}>
      {/* Left side: IDE/Code Editor (70% width) */}
      <div className="w-[70%] h-full overflow-hidden">
        {editor}
      </div>
      
      {/* Right side: Controls/Menu (30% width) */}
      <div className="w-[30%] min-w-[300px] h-full border-l border-gray-200 dark:border-gray-700 overflow-y-auto">
        {controls}
      </div>
    </div>
  );
};

export default FlippedLayout;
