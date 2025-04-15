"use client";

import React, { useEffect, useRef, useState } from "react";
import { cn } from "../../../utils/lib/utils";

export const TextReveal = ({
  text,
  className,
}: {
  text: string;
  className?: string;
}) => {
  const [isIntersecting, setIsIntersecting] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsIntersecting(entry.isIntersecting);
      },
      {
        rootMargin: "0px",
        threshold: 0.1,
      }
    );

    if (ref.current) {
      observer.observe(ref.current);
    }

    return () => {
      if (ref.current) {
        observer.unobserve(ref.current);
      }
    };
  }, []);

  return (
    <div ref={ref} className={cn("relative", className)}>
      <p
        className={cn(
          "absolute inset-0 text-center text-emerald-500 transition-all duration-1000",
          isIntersecting
            ? "translate-y-0 opacity-100"
            : "translate-y-8 opacity-0"
        )}
      >
        {text}
      </p>
      <p className="text-center text-transparent">{text}</p>
    </div>
  );
};
