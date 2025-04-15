import React, { useState, useEffect } from 'react';
import { useSharedState } from '../hooks/useSharedState';
import { SharedStateProvider } from '../contexts/SharedStateContext';

const SharedStateComponent = () => {
  const [message, setMessage] = useState('');
  const [agentState, setAgentState] = useState({});
  
  const { 
    state, 
    loading, 
    error, 
    updateState, 
    sendEvent 
  } = useSharedState('agent', 'agent-1', { status: 'idle', progress: 0 });

  useEffect(() => {
    if (state) {
      setAgentState(state);
    }
  }, [state]);

  const handleMessageChange = (e) => {
    setMessage(e.target.value);
  };

  const handleSendMessage = () => {
    if (!message.trim()) return;
    
    sendEvent({
      type: 'user_message',
      content: message,
      timestamp: new Date().toISOString(),
    });
    
    setMessage('');
  };

  const handleStartAgent = () => {
    updateState({
      status: 'running',
      progress: 0,
      startedAt: new Date().toISOString(),
    });
  };

  const handleStopAgent = () => {
    updateState({
      status: 'stopped',
      stoppedAt: new Date().toISOString(),
    });
  };

  const handleUpdateProgress = () => {
    const newProgress = Math.min(100, (agentState.progress || 0) + 10);
    updateState({
      progress: newProgress,
      ...(newProgress === 100 ? { status: 'completed', completedAt: new Date().toISOString() } : {}),
    });
  };

  if (loading) {
    return <div className="loading">Loading agent state...</div>;
  }

  if (error) {
    return <div className="error">Error: {error.message}</div>;
  }

  return (
    <div className="shared-state-example">
      <h2>Agent State</h2>
      <div className="state-display">
        <pre>{JSON.stringify(agentState, null, 2)}</pre>
      </div>

      <div className="controls">
        <h3>Controls</h3>
        <div className="button-group">
          <button 
            onClick={handleStartAgent}
            disabled={agentState.status === 'running'}
          >
            Start Agent
          </button>
          <button 
            onClick={handleStopAgent}
            disabled={agentState.status !== 'running'}
          >
            Stop Agent
          </button>
          <button 
            onClick={handleUpdateProgress}
            disabled={agentState.status !== 'running' || agentState.progress === 100}
          >
            Update Progress (+10%)
          </button>
        </div>

        <h3>Send Message</h3>
        <div className="message-input">
          <input
            type="text"
            value={message}
            onChange={handleMessageChange}
            placeholder="Type a message..."
          />
          <button onClick={handleSendMessage}>Send</button>
        </div>
      </div>
    </div>
  );
};

const SharedStateExample = () => {
  return (
    <SharedStateProvider serverUrl="ws://localhost:8080/ws">
      <div className="example-container">
        <h1>Shared State Example</h1>
        <SharedStateComponent />
      </div>
    </SharedStateProvider>
  );
};

export default SharedStateExample;
