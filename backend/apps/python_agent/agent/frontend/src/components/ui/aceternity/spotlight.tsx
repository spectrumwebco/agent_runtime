"use client";

import React, { useRef, useState, useEffect } from "react";
import { cn } from "../../../utils/lib/utils";

interface SpotlightProps {
  className?: string;
  children: React.ReactNode;
}

export const Spotlight = ({
  children,
  className = "",
}: SpotlightProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const mouseX = useRef(0);
  const mouseY = useRef(0);

  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (!isMounted || !containerRef.current) return;
    
    const handleMouseMove = (e: MouseEvent) => {
      const rect = containerRef.current!.getBoundingClientRect();
      mouseX.current = e.clientX - rect.left;
      mouseY.current = e.clientY - rect.top;
      
      if (containerRef.current) {
        containerRef.current.style.setProperty("--mouse-x", `${mouseX.current}px`);
        containerRef.current.style.setProperty("--mouse-y", `${mouseY.current}px`);
      }
    };
    
    if (containerRef.current) {
      containerRef.current.addEventListener("mousemove", handleMouseMove);
    }
    
    return () => {
      if (containerRef.current) {
        containerRef.current.removeEventListener("mousemove", handleMouseMove);
      }
    };
  }, [isMounted]);

  return (
    <div
      ref={containerRef}
      className={cn(
        "relative h-full w-full overflow-hidden rounded-md bg-background",
        className
      )}
    >
      <div
        className="pointer-events-none absolute inset-0 z-30 h-full w-full bg-[radial-gradient(circle_at_var(--mouse-x)_var(--mouse-y),rgba(16,185,129,0.15),transparent_40%)]"
        style={{
          "--mouse-x": "0px",
          "--mouse-y": "0px",
        } as React.CSSProperties}
      />
      <div className="relative z-10 h-full w-full">{children}</div>
    </div>
  );
};
