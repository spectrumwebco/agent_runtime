import React, { RefObject } from "react";
import Message from "../AgentMessage";
import MacBar from "../MacBar";
import editorLogo from "../../assets/panel_icons/editor.png";
import "../../static/agentFeed.css";

interface FeedItem {
  type: string;
  message: string;
  format?: string;
  step?: number;
}

interface AgentFeedProps {
  feed: FeedItem[];
  highlightedStep: number | null;
  handleMouseEnter: (item: FeedItem, feedRef: RefObject<HTMLDivElement | null>) => void;
  handleMouseLeave: () => void;
  selfRef: RefObject<HTMLDivElement | null>;
}

const AgentFeed: React.FC<AgentFeedProps> = ({
  feed,
  highlightedStep,
  handleMouseEnter,
  handleMouseLeave,
  selfRef,
}) => {
  return (
    <div id="agentFeed" className="agentFeed">
      <MacBar title="Agent Thoughts" logo={editorLogo} dark={true} />
      <div className="scrollableDiv" ref={selfRef}>
        <div className="innerDiv">
          {feed.map((item, index) => (
            <Message
              key={index}
              item={item}
              handleMouseEnter={handleMouseEnter}
              handleMouseLeave={handleMouseLeave}
              isHighlighted={
                item.step !== null && highlightedStep === item.step
              }
              feedRef={selfRef}
            />
          ))}
          <div style={{ clear: "both", marginTop: "1em" }} />
        </div>
      </div>
    </div>
  );
};

export default AgentFeed;
