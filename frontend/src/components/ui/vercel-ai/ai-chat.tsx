import React, { useState } from "react";
import { useChat } from "@vercel/ai";
import { cn } from "../../../utils/cn";
import { Button } from "../shadcn/button";

interface AIChatProps {
  className?: string;
  initialMessages?: { role: "user" | "assistant"; content: string }[];
  modelId?: string;
}

export const AIChat: React.FC<AIChatProps> = ({
  className,
  initialMessages = [],
  modelId = "gemini-2.5-pro",
}) => {
  const [input, setInput] = useState("");
  
  const { messages, handleSubmit, handleInputChange, isLoading } = useChat({
    api: "/api/chat",
    initialMessages,
    body: {
      model: modelId,
    },
  });

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleSubmit(e);
    setInput("");
  };

  return (
    <div className={cn("flex flex-col h-full", className)}>
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message, i) => (
          <div
            key={i}
            className={cn(
              "flex flex-col max-w-[80%] rounded-lg p-4",
              message.role === "user"
                ? "bg-emerald-100 dark:bg-emerald-900/20 ml-auto"
                : "bg-gray-100 dark:bg-gray-800 mr-auto"
            )}
          >
            <div className="text-xs font-medium mb-1">
              {message.role === "user" ? "You" : "Kled"}
            </div>
            <div className="whitespace-pre-wrap">{message.content}</div>
          </div>
        ))}
        {isLoading && (
          <div className="flex items-center justify-center py-4">
            <div className="animate-pulse text-emerald-500">
              Kled is thinking...
            </div>
          </div>
        )}
      </div>
      
      <div className="border-t border-gray-200 dark:border-gray-700 p-4">
        <form onSubmit={onSubmit} className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => {
              setInput(e.target.value);
              handleInputChange(e);
            }}
            placeholder="Ask Kled anything..."
            className="flex-1 rounded-md border border-gray-300 dark:border-gray-600 bg-transparent px-4 py-2 focus:outline-none focus:ring-2 focus:ring-emerald-500"
          />
          <Button type="submit" variant="emerald" disabled={isLoading || !input.trim()}>
            Send
          </Button>
        </form>
      </div>
    </div>
  );
};

export default AIChat;
