"use client";

import React from "react";
import { cn } from "../../../utils/lib/utils";

interface MeteorsProps {
  number?: number;
  className?: string;
}

export const Meteors = ({ number = 20, className }: MeteorsProps) => {
  const meteors = [...Array(number)].map((_, i) => (
    <span
      key={i}
      className={cn(
        "absolute left-1/2 top-1/2 h-0.5 w-0.5 rotate-[215deg] animate-meteor rounded-[9999px] bg-emerald-500 shadow-[0_0_0_1px_#ffffff10]",
        "before:absolute before:top-1/2 before:h-[1px] before:w-[50px] before:-translate-y-[50%] before:transform before:bg-gradient-to-r before:from-[#10b981] before:to-transparent before:content-['']"
      )}
      style={{
        top: `${Math.floor(Math.random() * 100)}%`,
        left: `${Math.floor(Math.random() * 100)}%`,
        animationDelay: `${Math.random() * 1}s`,
        animationDuration: `${Math.random() * 3 + 2}s`,
      }}
    />
  ));

  return (
    <div
      className={cn(
        "absolute inset-0 overflow-hidden",
        className
      )}
    >
      {meteors}
    </div>
  );
};
