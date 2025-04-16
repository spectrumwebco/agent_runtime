import React, { useState } from "react";
import { motion } from "framer-motion";
import { cn } from "../../../lib/utils";

interface Tab {
  id: string;
  label: string;
  content: React.ReactNode;
}

interface AnimatedTabsProps {
  tabs: Tab[];
  defaultTabId?: string;
  className?: string;
  tabClassName?: string;
  contentClassName?: string;
}

export const AnimatedTabs: React.FC<AnimatedTabsProps> = ({
  tabs,
  defaultTabId,
  className,
  tabClassName,
  contentClassName,
}) => {
  const [activeTab, setActiveTab] = useState<string>(defaultTabId || tabs[0]?.id || "");

  return (
    <div className={cn("w-full", className)}>
      <div className="flex border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={cn(
              "px-4 py-2 text-sm font-medium relative transition-colors",
              activeTab === tab.id
                ? "text-emerald-500"
                : "text-gray-600 dark:text-gray-400 hover:text-emerald-500 dark:hover:text-emerald-400",
              tabClassName
            )}
          >
            {tab.label}
            {activeTab === tab.id && (
              <motion.div
                layoutId="active-tab-indicator"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-emerald-500"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3 }}
              />
            )}
          </button>
        ))}
      </div>
      <div className={cn("py-4", contentClassName)}>
        {tabs.map((tab) => (
          <div
            key={tab.id}
            className={cn(
              "transition-opacity duration-300",
              activeTab === tab.id ? "block" : "hidden"
            )}
          >
            {tab.content}
          </div>
        ))}
      </div>
    </div>
  );
};
