"use client";

import React, { useState, useEffect } from "react";
import { cn } from "../../../utils/lib/utils";

interface FloatingNavbarProps {
  navItems: {
    name: string;
    link: string;
    icon?: React.ReactNode;
  }[];
  className?: string;
}

export const FloatingNavbar = ({
  navItems,
  className,
}: FloatingNavbarProps) => {
  const [activeIndex, setActiveIndex] = useState<number | null>(0);
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  return (
    <div
      className={cn(
        "fixed bottom-4 left-1/2 transform -translate-x-1/2 z-50",
        className
      )}
    >
      <div className="flex items-center justify-center space-x-4 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md rounded-full p-2 px-4 shadow-lg border border-white/20 dark:border-gray-800/20">
        {navItems.map((item, index) => (
          <button
            key={index}
            onClick={() => setActiveIndex(index)}
            className={cn(
              "relative px-4 py-2 rounded-full transition-all duration-300",
              activeIndex === index
                ? "text-white"
                : "text-gray-500 dark:text-gray-400 hover:text-emerald-500 dark:hover:text-emerald-400"
            )}
          >
            {activeIndex === index && (
              <span className="absolute inset-0 rounded-full bg-emerald-500 dark:bg-emerald-600" />
            )}
            <span className="relative flex items-center space-x-2">
              {item.icon && <span>{item.icon}</span>}
              <span>{item.name}</span>
            </span>
          </button>
        ))}
      </div>
    </div>
  );
};
