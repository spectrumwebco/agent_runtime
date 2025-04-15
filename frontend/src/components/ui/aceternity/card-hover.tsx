"use client";

import React, { useState, useRef, useEffect } from "react";
import { cn } from "../../../utils/lib/utils";

interface CardHoverEffectProps {
  items: {
    title: string;
    description: string;
    icon?: React.ReactNode;
  }[];
  className?: string;
}

export const CardHoverEffect = ({
  items,
  className,
}: CardHoverEffectProps) => {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  return (
    <div
      className={cn(
        "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 py-10",
        className
      )}
    >
      {items.map((item, idx) => (
        <div
          key={idx}
          className="relative group block p-2 h-full w-full"
          onMouseEnter={() => setHoveredIndex(idx)}
          onMouseLeave={() => setHoveredIndex(null)}
        >
          <div
            className={cn(
              "absolute inset-0 rounded-lg bg-emerald-50 dark:bg-emerald-950/20 opacity-0 group-hover:opacity-100 transition-opacity",
              hoveredIndex === idx &&
                "opacity-100 border border-emerald-500/30"
            )}
          />
          <div
            className={cn(
              "relative z-10 p-5 rounded-lg bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 h-full",
              hoveredIndex === idx &&
                "border-emerald-500/50 shadow-xl shadow-emerald-500/10"
            )}
          >
            {item.icon && (
              <div className="flex justify-center items-center w-12 h-12 rounded-full bg-emerald-500/10 text-emerald-500 mb-4">
                {item.icon}
              </div>
            )}
            <h3 className="text-lg font-medium mb-2">{item.title}</h3>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              {item.description}
            </p>
          </div>
        </div>
      ))}
    </div>
  );
};
