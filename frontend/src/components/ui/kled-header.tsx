import React from "react";
import { cn } from "../../utils/cn";
import { TextGradient } from "./aceternity/text-gradient";

interface KledHeaderProps {
  className?: string;
}

export const KledHeader: React.FC<KledHeaderProps> = ({ className }) => {
  return (
    <header className={cn("border-b border-gray-200 dark:border-gray-700 p-4", className)}>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-full bg-emerald-500 flex items-center justify-center">
            <span className="text-white font-bold">K</span>
          </div>
          <TextGradient>
            <h1 className="text-xl font-bold">Kled.io</h1>
          </TextGradient>
          <span className="text-xs bg-emerald-100 dark:bg-emerald-900/30 text-emerald-800 dark:text-emerald-300 px-2 py-0.5 rounded-full">
            by Spectrum Web Co
          </span>
        </div>
        
        <div className="flex items-center gap-4">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            <span className="font-medium">Agent:</span> Kled
          </div>
          <div className="text-sm text-gray-500 dark:text-gray-400">
            <span className="font-medium">Role:</span> Senior Software Engineering Lead & Technical Authority for AI/ML
          </div>
        </div>
      </div>
    </header>
  );
};

export default KledHeader;
