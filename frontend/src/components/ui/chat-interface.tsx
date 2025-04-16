import React, { useState, useRef, useEffect } from "react";
import { cn } from "../../utils/cn";
import { SpotlightCard } from "./aceternity/spotlight-card";
import { GradientButton } from "./aceternity/gradient-button";

interface Message {
  id: string;
  role: "user" | "assistant" | "system";
  content: string;
  timestamp: Date;
}

interface ChatInterfaceProps {
  className?: string;
}

export const ChatInterface: React.FC<ChatInterfaceProps> = ({
  className,
}) => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      role: "system",
      content: "Welcome to Kled.io Agent Runtime. How can I assist you today?",
      timestamp: new Date(),
    },
  ]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = () => {
    if (!input.trim()) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: "user",
      content: input,
      timestamp: new Date(),
    };
    
    setMessages((prev) => [...prev, userMessage]);
    setInput("");
    setIsLoading(true);

    setTimeout(() => {
      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: "I'm the Kled.io Agent Runtime assistant. I'm processing your request and will help you with your task.",
        timestamp: new Date(),
      };
      
      setMessages((prev) => [...prev, assistantMessage]);
      setIsLoading(false);
    }, 1500);
  };

  const formatTimestamp = (date: Date) => {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <SpotlightCard className={cn("flex flex-col h-full", className)}>
      <div className="flex justify-between items-center p-4 border-b">
        <h3 className="text-lg font-medium">Kled.io Agent Chat</h3>
        <div className="flex items-center gap-2">
          <span className={cn(
            "h-2 w-2 rounded-full",
            isLoading ? "bg-yellow-500" : "bg-emerald-500"
          )}></span>
          <span className="text-sm">
            {isLoading ? "Processing..." : "Ready"}
          </span>
        </div>
      </div>
      
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message) => (
          <div
            key={message.id}
            className={cn(
              "flex flex-col max-w-[80%] rounded-lg p-3",
              message.role === "user"
                ? "ml-auto bg-emerald-500 text-white"
                : message.role === "system"
                ? "mx-auto bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200"
                : "mr-auto bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200"
            )}
          >
            <div className="text-sm">{message.content}</div>
            <div className="text-xs mt-1 opacity-70 self-end">
              {formatTimestamp(message.timestamp)}
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>
      
      <div className="p-4 border-t">
        <div className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                handleSendMessage();
              }
            }}
            placeholder="Type your message..."
            className="flex-1 px-4 py-2 rounded-md border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-emerald-500"
            disabled={isLoading}
          />
          <GradientButton
            onClick={handleSendMessage}
            disabled={isLoading || !input.trim()}
          >
            Send
          </GradientButton>
        </div>
      </div>
    </SpotlightCard>
  );
};
