import React from "react";
import { cn } from "../../../utils/cn";

interface GradientButtonProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
  variant?: "default" | "outline" | "secondary";
}

export const GradientButton: React.FC<GradientButtonProps> = ({
  children,
  className,
  onClick,
  disabled = false,
  variant = "default",
}) => {
  const baseClasses =
    "relative inline-flex items-center justify-center px-4 py-2 rounded-md font-medium text-sm transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500";

  const variantClasses = {
    default: "text-white border border-emerald-500 shadow-md",
    outline:
      "bg-transparent border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800",
    secondary:
      "bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 border border-gray-200 dark:border-gray-700 hover:bg-gray-200 dark:hover:bg-gray-700",
  };

  const disabledClasses = "opacity-50 cursor-not-allowed";

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={cn(
        baseClasses,
        variantClasses[variant],
        disabled && disabledClasses,
        className,
      )}
    >
      {variant === "default" && !disabled && (
        <div className="absolute inset-0 rounded-md overflow-hidden">
          <div
            className="absolute inset-0 bg-gradient-to-r from-emerald-400 to-emerald-600 animate-gradient"
            style={{
              backgroundSize: "200% 200%",
              animation: "gradientAnimation 3s ease infinite",
            }}
          />
        </div>
      )}
      <span className="relative z-10">{children}</span>
    </button>
  );
};

export default GradientButton;
