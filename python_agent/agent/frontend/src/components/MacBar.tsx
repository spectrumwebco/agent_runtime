import React from "react";

import "../static/macbar.css";

interface MacBarProps {
  title: string;
  logo: string;
  dark?: boolean;
  className?: string;
}

const MacBar: React.FC<MacBarProps> = ({ title, logo, dark = false, className = "" }) => {
  const darkClass = dark ? "dark" : "";
  return (
    <div className={`mac-window-top-bar ${darkClass} ${className}`}>
      <div className="label">
        <img src={logo} alt={title} />
        <span>{title}</span>
      </div>
    </div>
  );
};

export default MacBar;
