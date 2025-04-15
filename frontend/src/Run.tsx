import React, { useState, useRef, useEffect } from "react";
import axios from "axios";
import io from "socket.io-client";
import "./static/run.css";
import AgentFeed from "./components/panels/AgentFeed";
import EnvFeed from "./components/panels/EnvFeed";
import LogPanel from "./components/panels/LogPanel";
import LRunControl from "./components/controls/LRunControl";
import { useImmer } from "use-immer";
import { Socket } from "socket.io-client";

interface RunConfig {
  agent: {
    model: {
      model_name: string;
    };
  };
  problem_statement: {
    type: string;
    input: string;
  };
  environment: {
    image_name: string;
    script: string;
    repo: {
      type: string;
      input: string;
    };
  };
  extra: {
    test_run: boolean;
  };
}

interface FeedItem {
  type: string;
  message: string;
  format?: string;
  step?: number;
}

interface SocketMessage {
  feed?: string;
  type?: string;
  message?: string;
  format?: string;
  thought_idx?: number;
}

const url = ""; // Will get this from .env
const socket: Socket = io(url);

const Run: React.FC = () => {
  const [isConnected, setIsConnected] = useState<boolean>(socket.connected);
  const [errorBanner, setErrorBanner] = useState<string>("");

  const runConfigDefault: RunConfig = {
    agent: {
      model: {
        model_name: "gpt4",
      },
    },
    problem_statement: {
      type: "",
      input: "",
    },
    environment: {
      image_name: "",
      script: "",
      repo: {
        type: "",
        input: "",
      },
    },
    extra: {
      test_run: false,
    },
  };
  const [runConfig, setRunConfig] = useImmer<RunConfig>(runConfigDefault);

  const [agentFeed, setAgentFeed] = useState<FeedItem[]>([]);
  const [envFeed, setEnvFeed] = useState<FeedItem[]>([]);
  const [highlightedStep, setHighlightedStep] = useState<number | null>(null);
  const [logs, setLogs] = useState<string>("");
  const [isComputing, setIsComputing] = useState<boolean>(false);

  const hoverTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const agentFeedRef = useRef<HTMLDivElement | null>(null);
  const envFeedRef = useRef<HTMLDivElement | null>(null);
  const logsRef = useRef<HTMLDivElement | null>(null);
  const isLogScrolled = useRef<boolean>(false);
  const isEnvScrolled = useRef<boolean>(false);
  const isAgentScrolled = useRef<boolean>(false);

  const [tabKey, setTabKey] = useState<string | null>("problem");

  const stillComputingTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  axios.defaults.baseURL = url;

  function scrollToHighlightedStep(highlightedStep: number | null, ref: React.RefObject<HTMLDivElement | null>): void {
    if (highlightedStep && ref.current) {
      console.log(
        "Scrolling to highlighted step",
        highlightedStep,
        ref.current,
      );
      const firstStepMessage = ref.current.querySelector(
        `.step${highlightedStep}`,
      );
      if (firstStepMessage && firstStepMessage instanceof HTMLElement) {
        window.requestAnimationFrame(() => {
          if (ref.current) {
            ref.current.scrollTo({
              top: firstStepMessage.offsetTop - ref.current.offsetTop,
              behavior: "smooth",
            });
          }
        });
      }
    }
  }

  function getOtherFeed(feedRef: React.RefObject<HTMLDivElement | null>): React.RefObject<HTMLDivElement | null> {
    return feedRef === agentFeedRef ? envFeedRef : agentFeedRef;
  }

  const handleMouseEnter = (item: FeedItem, feedRef: React.RefObject<HTMLDivElement | null>): void => {
    if (isComputing) {
      return;
    }

    const stepHighlight = item.step;

    if (hoverTimeoutRef.current) {
      clearTimeout(hoverTimeoutRef.current);
    }

    hoverTimeoutRef.current = setTimeout(() => {
      if (!isComputing) {
        setHighlightedStep(stepHighlight ?? null);
        scrollToHighlightedStep(stepHighlight ?? null, getOtherFeed(feedRef));
      }
    }, 250);
  };

  const handleMouseLeave = (): void => {
    console.log("Mouse left");
    if (hoverTimeoutRef.current) {
      clearTimeout(hoverTimeoutRef.current);
    }
    setHighlightedStep(null);
  };

  const requeueStopComputeTimeout = (): void => {
  };

  const handleSubmit = async (event: React.FormEvent): Promise<void> => {
    setTabKey(null);
    setIsComputing(true);
    event.preventDefault();
    setAgentFeed([]);
    setEnvFeed([]);
    setLogs("");
    setHighlightedStep(null);
    setErrorBanner("");
    try {
      await axios.get(`/run`, {
        params: { runConfig: JSON.stringify(runConfig) },
      });
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const handleStop = async (): Promise<void> => {
    setIsComputing(false);
    try {
      const response = await axios.get("/stop");
      console.log(response.data);
    } catch (error) {
      console.error("Error stopping:", error);
    }
  };

  const checkScrollPosition = (
    ref: React.RefObject<HTMLDivElement | null>,
    scrollStateRef: React.MutableRefObject<boolean>,
    offset = 0
  ): void => {
    if (ref.current) {
      scrollStateRef.current =
        ref.current.scrollTop + ref.current.clientHeight + offset <
        ref.current.scrollHeight;
    }
  };

  const scrollToBottom = (
    ref: React.RefObject<HTMLDivElement | null>,
    scrollStateRef: React.MutableRefObject<boolean>
  ): void => {
    if (ref.current && !scrollStateRef.current) {
      ref.current.scrollTop = ref.current.scrollHeight;
    }
  };

  const scrollDetectedLog = (): void =>
    checkScrollPosition(logsRef, isLogScrolled, 58);
  
  const scrollDetectedEnv = (): void =>
    checkScrollPosition(envFeedRef, isEnvScrolled);
  
  const scrollDetectedAgent = (): void =>
    checkScrollPosition(agentFeedRef, isAgentScrolled);
  
  const scrollLog = (): void => scrollToBottom(logsRef, isLogScrolled);
  const scrollEnv = (): void => scrollToBottom(envFeedRef, isEnvScrolled);
  const scrollAgent = (): void => scrollToBottom(agentFeedRef, isAgentScrolled);

  useEffect(() => {
    if (logsRef.current) {
      logsRef.current.addEventListener("scroll", scrollDetectedLog, {
        passive: true,
      });
    }
    
    if (envFeedRef.current) {
      envFeedRef.current.addEventListener("scroll", scrollDetectedEnv, {
        passive: true,
      });
    }
    
    if (agentFeedRef.current) {
      agentFeedRef.current.addEventListener("scroll", scrollDetectedAgent, {
        passive: true,
      });
    }

    const handleUpdate = (data: SocketMessage): void => {
      requeueStopComputeTimeout();
      if (data.feed === "agent") {
        setAgentFeed((prevMessages) => [
          ...prevMessages,
          {
            type: data.type || "",
            message: data.message || "",
            format: data.format,
            step: data.thought_idx,
          },
        ]);
        if (envFeedRef.current) {
          setTimeout(() => {
            scrollEnv();
          }, 100);
        }
      } else if (data.feed === "env") {
        setEnvFeed((prevMessages) => [
          ...prevMessages,
          {
            message: data.message || "",
            type: data.type || "",
            format: data.format,
            step: data.thought_idx,
          },
        ]);
        if (agentFeedRef.current) {
          setTimeout(() => {
            scrollAgent();
          }, 100);
        }
      }
    };

    const handleUpdateBanner = (data: SocketMessage): void => {
      if (data.message) {
        setErrorBanner(data.message);
      }
    };

    const handleLogMessage = (data: SocketMessage): void => {
      requeueStopComputeTimeout();
      if (data.message) {
        setLogs((prevLogs) => prevLogs + data.message);
        if (logsRef.current) {
          setTimeout(() => {
            scrollLog();
          }, 100);
        }
      }
    };

    const handleFinishedRun = (): void => {
      setIsComputing(false);
    };

    socket.on("update", handleUpdate);
    socket.on("log_message", handleLogMessage);
    socket.on("update_banner", handleUpdateBanner);
    socket.on("finish_run", handleFinishedRun);
    socket.on("connect", () => {
      console.log("Connected to server");
      setIsConnected(true);
      setErrorBanner("");
    });

    socket.on("disconnect", () => {
      console.log("Disconnected from server");
      setIsConnected(false);
      setErrorBanner("Connection to flask server lost, please restart it.");
      setIsComputing(false);
      scrollLog(); // reveal copy button
    });

    socket.on("connect_error", () => {
      setIsConnected(false);
      setErrorBanner(
        "Failed to connect to the flask server, please restart it.",
      );
      setIsComputing(false);
      scrollLog(); // reveal copy button
    });

    return () => {
      if (logsRef.current) {
        logsRef.current.removeEventListener("scroll", scrollDetectedLog);
      }
      if (envFeedRef.current) {
        envFeedRef.current.removeEventListener("scroll", scrollDetectedEnv);
      }
      if (agentFeedRef.current) {
        agentFeedRef.current.removeEventListener("scroll", scrollDetectedAgent);
      }
      
      socket.off("update", handleUpdate);
      socket.off("log_message", handleLogMessage);
      socket.off("finish_run", handleFinishedRun);
      socket.off("connect");
      socket.off("disconnect");
      socket.off("connect_error");
      socket.off("update_banner", handleUpdateBanner);
    };
  }, []);

  function renderErrorMessage(): JSX.Element | null {
    if (errorBanner) {
      return (
        <div className="alert alert-danger" role="alert">
          {errorBanner}
          <br />
          If you think this was a bug, please head over to{" "}
          <a
            href="https://github.com/SWE-agent/SWE-agent/issues"
            target="blank"
          >
            our GitHub issue tracker
          </a>
          , check if someone has already reported the issue, and if not, create
          a new issue. Please include the full log, all settings that you
          entered, and a screenshot of this page.
        </div>
      );
    }
    return null;
  }

  return (
    <div className="container-demo">
      {renderErrorMessage()}
      <LRunControl
        isComputing={isComputing}
        isConnected={isConnected}
        handleStop={handleStop}
        handleSubmit={handleSubmit}
        tabKey={tabKey}
        setTabKey={setTabKey}
        runConfig={runConfig}
        setRunConfig={setRunConfig}
        runConfigDefault={runConfigDefault}
      />
      <div id="demo">
        <hr />
        {/* Flipped Devin-inspired layout with IDE on left, menu on right */}
        <div className="panels devin-flipped-layout">
          {/* IDE Panel (EnvFeed) on the left */}
          <EnvFeed
            feed={envFeed}
            highlightedStep={highlightedStep}
            handleMouseEnter={handleMouseEnter}
            handleMouseLeave={handleMouseLeave}
            selfRef={envFeedRef}
          />
          {/* Agent Feed (Thoughts/Menu) on the right */}
          <AgentFeed
            feed={agentFeed}
            highlightedStep={highlightedStep}
            handleMouseEnter={handleMouseEnter}
            handleMouseLeave={handleMouseLeave}
            selfRef={agentFeedRef}
          />
          {/* Log Panel at the bottom */}
          <LogPanel logs={logs} logsRef={logsRef} isComputing={isComputing} />
        </div>
      </div>
      <hr />
    </div>
  );
};

export default Run;
