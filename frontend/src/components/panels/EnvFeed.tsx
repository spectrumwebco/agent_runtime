import React, { RefObject } from "react";
import EnvMessage from "../EnvMessage";

import MacBar from "../MacBar";
import terminalLogo from "../../assets/panel_icons/terminal.png";
import "../../static/envFeed.css";

interface FeedItem {
  type: string;
  message: string;
  format?: string;
  step?: number | null;
}

interface EnvFeedProps {
  feed: FeedItem[];
  highlightedStep: number | null;
  handleMouseEnter: (item: FeedItem, feedRef: RefObject<HTMLDivElement | null>) => void;
  handleMouseLeave: () => void;
  selfRef: RefObject<HTMLDivElement | null>;
}

const EnvFeed: React.FC<EnvFeedProps> = ({
  feed,
  highlightedStep,
  handleMouseEnter,
  handleMouseLeave,
  selfRef,
}) => {
  return (
    <div id="envFeed" className="envFeed">
      <MacBar title="IDE / Terminal" logo={terminalLogo} dark={false} />
      <div className="scrollableDiv" ref={selfRef}>
        <div className="innerDiv">
          {feed.map((item, index) => (
            <EnvMessage
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

export default EnvFeed;
