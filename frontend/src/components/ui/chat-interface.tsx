import React, { useState } from "react";
import { cn } from "../../utils/cn";
import { useSharedState } from "./shared-state-provider";

interface ChatInterfaceProps {
  className?: string;
}

type Message = {
  id: string;
  content: string;
  sender: "user" | "agent";
  timestamp: Date;
};

export const ChatInterface: React.FC<ChatInterfaceProps> = ({ className }) => {
  const { state, updateState } = useSharedState();
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSendMessage = () => {
    if (!inputValue.trim() || isProcessing) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      content: inputValue,
      sender: "user",
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputValue("");
    setIsProcessing(true);

    setTimeout(() => {
      const agentMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: `I'm processing your request: "${inputValue}"`,
        sender: "agent",
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, agentMessage]);
      setIsProcessing(false);
    }, 1000);
  };

  const formatTimestamp = (date: Date) => {
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 shadow-md flex flex-col",
        className,
      )}
    >
      <div className="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
        <div>
          <h3 className="text-lg font-medium">Chat Interface</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Communicate with the Kled.io agent
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span
            className={cn(
              "h-2 w-2 rounded-full",
              state.isGenerating ? "bg-emerald-500" : "bg-gray-400",
            )}
          ></span>
          <span className="text-sm">
            {state.isGenerating ? "Active" : "Inactive"}
          </span>
        </div>
      </div>

      <div className="flex-1 p-4 overflow-y-auto max-h-80">
        {messages.length === 0 ? (
          <div className="text-center text-gray-500 dark:text-gray-400 py-8">
            No messages yet. Start a conversation with the agent.
          </div>
        ) : (
          <div className="space-y-4">
            {messages.map((message) => (
              <div
                key={message.id}
                className={cn(
                  "p-3 rounded-lg max-w-[80%]",
                  message.sender === "user"
                    ? "bg-emerald-500 text-white ml-auto"
                    : "bg-gray-700 dark:bg-gray-800 text-white mr-auto",
                )}
              >
                <div className="text-sm">{message.content}</div>
                <div className="text-xs mt-1 opacity-70">
                  {formatTimestamp(message.timestamp)}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="p-4 border-t border-gray-200 dark:border-gray-700">
        <div className="flex gap-2">
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
            placeholder="Type a message..."
            className="flex-1 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 text-gray-900 dark:text-white"
            disabled={state.isGenerating || isProcessing}
          />
          <button
            onClick={handleSendMessage}
            disabled={!inputValue.trim() || state.isGenerating || isProcessing}
            className={cn(
              "px-4 py-2 rounded-md font-medium text-sm transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500",
              inputValue.trim() && !state.isGenerating && !isProcessing
                ? "bg-gradient-to-r from-emerald-400 to-emerald-600 text-white"
                : "bg-gray-300 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed",
            )}
          >
            {isProcessing ? "Sending..." : "Send"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ChatInterface;
