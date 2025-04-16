import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { cn } from "../../../lib/utils";

export const FloatingNav = ({
  navItems,
  className,
}: {
  navItems: {
    name: string;
    link: string;
    icon?: JSX.Element;
  }[];
  className?: string;
}) => {
  const [isVisible, setIsVisible] = useState(true);
  const [activeIndex, setActiveIndex] = useState<number | null>(null);
  const [lastScrollY, setLastScrollY] = useState(0);

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      if (currentScrollY > lastScrollY && isVisible) {
        setIsVisible(false);
      } else if (currentScrollY < lastScrollY && !isVisible) {
        setIsVisible(true);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, [isVisible, lastScrollY]);

  return (
    <AnimatePresence mode="wait">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{
          opacity: isVisible ? 1 : 0,
          y: isVisible ? 0 : -20,
        }}
        transition={{ duration: 0.2 }}
        className={cn(
          "fixed top-4 inset-x-0 max-w-fit mx-auto z-50 flex items-center justify-center space-x-4 px-4 py-2 rounded-full border border-transparent dark:border-white/[0.2] bg-white dark:bg-black shadow-[0px_2px_3px_-1px_rgba(0,0,0,0.1),0px_1px_0px_0px_rgba(25,28,33,0.02),0px_0px_0px_1px_rgba(25,28,33,0.08)] dark:shadow-[0px_2px_3px_-1px_rgba(0,0,0,0.1),0px_1px_0px_0px_rgba(25,28,33,0.02),0px_0px_0px_1px_rgba(255,255,255,0.08)]",
          className
        )}
      >
        {navItems.map((navItem, index) => (
          <a
            key={`nav-item-${index}`}
            href={navItem.link}
            onMouseEnter={() => setActiveIndex(index)}
            onMouseLeave={() => setActiveIndex(null)}
            className={cn(
              "relative px-4 py-2 rounded-full text-sm font-medium transition-colors",
              activeIndex === index
                ? "text-white"
                : "text-neutral-700 dark:text-neutral-300 hover:text-emerald-500 dark:hover:text-emerald-400"
            )}
          >
            <span className="relative z-10 flex items-center gap-2">
              {navItem.icon && navItem.icon}
              <span>{navItem.name}</span>
            </span>
            {activeIndex === index && (
              <motion.div
                layoutId="pill-indicator"
                transition={{ type: "spring", duration: 0.5 }}
                className="absolute inset-0 rounded-full bg-emerald-500"
              ></motion.div>
            )}
          </a>
        ))}
      </motion.div>
    </AnimatePresence>
  );
};
