import React, { RefObject } from "react";
import MacBar from "../MacBar";
import workspaceLogo from "../../assets/panel_icons/workspace.png";
import "../../static/logPanel.css";
import { Button } from "react-bootstrap";
import { Clipboard } from "react-bootstrap-icons";

interface LogPanelProps {
  logs: string;
  logsRef: RefObject<HTMLDivElement | null>;
  isComputing: boolean;
}

const LogPanel: React.FC<LogPanelProps> = ({ logs, logsRef, isComputing }) => {
  const copyToClipboard = (text: string): void => {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    document.body.appendChild(textarea);

    textarea.select();
    document.execCommand("copy");

    document.body.removeChild(textarea);
  };

  const handleCopy = (): void => {
    const contentElement = document.getElementById("logContent");
    if (contentElement) {
      const contentToCopy = contentElement.innerText;
      copyToClipboard(contentToCopy);
    }
  };

  return (
    <div id="logPanel" className="logPanel">
      <MacBar logo={workspaceLogo} title="Log file" dark={true} />
      <div className="scrollableDiv" ref={logsRef}>
        <div className="innerDiv">
          <pre id="logContent">{logs}</pre>
          <div style={{ clear: "both", marginTop: "1em" }} />
        </div>
        {!isComputing && logs && (
          <div
            style={{ display: "flex", justifyContent: "center", width: "100%" }}
          >
            <Button
              variant="light"
              onClick={handleCopy}
              style={{ marginBottom: 20, marginRight: 20 }}
            >
              <Clipboard /> Copy to clipboard
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

export default LogPanel;
