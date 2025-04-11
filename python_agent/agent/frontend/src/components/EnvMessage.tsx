import React, { RefObject } from "react";

import "../static/message.css";
import "../static/envMessage.css";

import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import bash from "react-syntax-highlighter/dist/esm/languages/prism/bash";
import { prism } from "react-syntax-highlighter/dist/esm/styles/prism";


interface FeedItem {
  type: string;
  message: string;
  format?: string;
  step?: number | null;
}

interface EnvMessageProps {
  item: FeedItem;
  handleMouseEnter: (item: FeedItem, feedRef: RefObject<HTMLDivElement | null>) => void;
  handleMouseLeave: () => void;
  isHighlighted: boolean;
  feedRef: RefObject<HTMLDivElement | null>;
}

function capitalizeFirstLetter(str: string): string {
  return str[0].toUpperCase() + str.slice(1);
}

const EnvMessage: React.FC<EnvMessageProps> = ({
  item,
  handleMouseEnter,
  handleMouseLeave,
  isHighlighted,
  feedRef,
}) => {
  const stepClass = item.step !== null ? `step${item.step}` : "";
  const highlightClass = isHighlighted ? "highlight" : "";
  const messageTypeClass = "envMessage" + capitalizeFirstLetter(item.type);

  const paddingBottom = item.type === "command" ? "0" : "0.5em";
  const paddingTop = ["output", "diff"].includes(item.type) ? "0" : "0.5em";

  const customStyle = {
    margin: 0,
    padding: `${paddingTop} 0.5em ${paddingBottom} 0.5em`,
    overflowX: "hidden",
    overflowY: "hidden",
    lineHeight: "100%",
    backgroundColor: "transparent",
    fontSize: "93%",
  };

  const codeTagProps = {
    style: {
      boxShadow: "none",
      margin: "0",
      overflowY: "hidden",
      overflowX: "hidden",
      padding: "0",
      lineHeight: "inherit",
      fontSize: "93%",
    },
  };

  const typeToLanguage: Record<string, string> = {
    command: "bash",
    output: "markdown",
    diff: "diff",
  };

  if (item.format !== "text") {
    return (
      <div
        className={`message envMessage ${stepClass} ${highlightClass}  ${messageTypeClass}`}
        onMouseEnter={() => handleMouseEnter(item, feedRef)}
        onMouseLeave={handleMouseLeave}
      >
        <SyntaxHighlighter
          codeTagProps={codeTagProps}
          customStyle={customStyle}
          language={typeToLanguage[item.type]}
          style={{ backgroundColor: "transparent", ...prism }}
          wrapLines={true}
          showLineNumbers={false}
        >
          {item.message}
        </SyntaxHighlighter>
      </div>
    );
  } else {
    return (
      <div
        className={`message ${stepClass} ${highlightClass} ${messageTypeClass}`}
        onMouseEnter={() => handleMouseEnter(item, feedRef)}
        onMouseLeave={handleMouseLeave}
      >
        <span>{item.message}</span>
      </div>
    );
  }
};

export default EnvMessage;
