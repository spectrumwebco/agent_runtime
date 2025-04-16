import React, { useState } from "react";
import { cn } from "../../utils/cn";
import { SpotlightCard } from "./aceternity/spotlight-card";
import { GradientButton } from "./aceternity/gradient-button";

interface CoPilotIntegrationProps {
  className?: string;
}

export const CoPilotIntegration: React.FC<CoPilotIntegrationProps> = ({
  className,
}) => {
  const [isConnected, setIsConnected] = useState(false);
  const [agentMessages, setAgentMessages] = useState<{id: string, message: string, timestamp: Date}[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  
  const coAgentConfig = {
    backendUrl: "http://localhost:8000/api",
    wsUrl: "ws://localhost:8000/ws",
    agentId: "kled-coagent-01",
    modelId: "gemini-2.5-pro",
  };
  
  const handleConnect = () => {
    setIsProcessing(true);
    
    setTimeout(() => {
      setIsConnected(true);
      setIsProcessing(false);
      
      setAgentMessages(prev => [
        ...prev, 
        {
          id: Date.now().toString(),
          message: "CoPilotKit CoAgent connected successfully",
          timestamp: new Date()
        }
      ]);
    }, 1500);
  };
  
  const handleDisconnect = () => {
    setIsProcessing(true);
    
    setTimeout(() => {
      setIsConnected(false);
      setIsProcessing(false);
      
      setAgentMessages(prev => [
        ...prev, 
        {
          id: Date.now().toString(),
          message: "CoPilotKit CoAgent disconnected",
          timestamp: new Date()
        }
      ]);
    }, 1000);
  };
  
  const handleSendTestMessage = () => {
    if (!isConnected) return;
    
    setIsProcessing(true);
    
    setTimeout(() => {
      setAgentMessages(prev => [
        ...prev, 
        {
          id: Date.now().toString(),
          message: "Test message sent to CoAgent",
          timestamp: new Date()
        }
      ]);
      
      setTimeout(() => {
        setAgentMessages(prev => [
          ...prev, 
          {
            id: (Date.now() + 1).toString(),
            message: "Response received from CoAgent: Ready to assist with UI generation and agent coordination",
            timestamp: new Date()
          }
        ]);
        
        setIsProcessing(false);
      }, 1000);
    }, 500);
  };
  
  const formatTimestamp = (date: Date) => {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  };

  return (
    <SpotlightCard className={cn("flex flex-col", className)}>
      <div className="flex justify-between items-center p-4 border-b">
        <div>
          <h3 className="text-lg font-medium">CoPilotKit CoAgents</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Integration with AI-powered agents
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className={cn(
            "h-2 w-2 rounded-full",
            isConnected ? "bg-emerald-500" : "bg-gray-400"
          )}></span>
          <span className="text-sm">
            {isConnected ? "Connected" : "Disconnected"}
          </span>
        </div>
      </div>
      
      <div className="p-4 space-y-4">
        <div className="border rounded-md p-3 bg-white/5">
          <h4 className="text-sm font-medium mb-2">CoAgent Configuration</h4>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <div>Backend URL:</div>
            <div className="font-mono">{coAgentConfig.backendUrl}</div>
            
            <div>WebSocket URL:</div>
            <div className="font-mono">{coAgentConfig.wsUrl}</div>
            
            <div>Agent ID:</div>
            <div className="font-mono">{coAgentConfig.agentId}</div>
            
            <div>Model ID:</div>
            <div className="font-mono">{coAgentConfig.modelId}</div>
          </div>
        </div>
        
        <div className="flex gap-2">
          {!isConnected ? (
            <GradientButton
              onClick={handleConnect}
              disabled={isProcessing}
              className="flex-1"
            >
              Connect to CoAgent
            </GradientButton>
          ) : (
            <GradientButton
              variant="outline"
              onClick={handleDisconnect}
              disabled={isProcessing}
              className="flex-1"
            >
              Disconnect
            </GradientButton>
          )}
          
          <GradientButton
            variant="secondary"
            onClick={handleSendTestMessage}
            disabled={!isConnected || isProcessing}
          >
            Send Test Message
          </GradientButton>
        </div>
        
        {agentMessages.length > 0 && (
          <div className="border rounded-md p-3 bg-white/5 max-h-40 overflow-y-auto">
            <h4 className="text-sm font-medium mb-2">Agent Messages</h4>
            <div className="space-y-2">
              {agentMessages.map((msg) => (
                <div key={msg.id} className="text-sm border-l-2 border-emerald-500 pl-2">
                  <span className="text-gray-500 text-xs">[{formatTimestamp(msg.timestamp)}]</span>{" "}
                  {msg.message}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </SpotlightCard>
  );
};

export default CoPilotIntegration;
