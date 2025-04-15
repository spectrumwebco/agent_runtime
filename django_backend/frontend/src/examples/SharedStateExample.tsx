import React, { useState } from 'react';
import useSharedState from '../hooks/useSharedState';

/**
 * Example component demonstrating the use of shared state
 * between React frontend and Go/Python backend.
 */
const SharedStateExample: React.FC = () => {
  const [inputValue, setInputValue] = useState('');
  
  const { state, updateState, connected, error } = useSharedState({
    stateType: 'shared',
    stateId: 'example',
    initialState: {
      messages: [],
      count: 0
    }
  });
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (inputValue.trim()) {
      updateState({
        ...state,
        messages: [...state.messages, inputValue],
        count: state.count + 1
      });
      
      setInputValue('');
    }
  };
  
  return (
    <div className="shared-state-example">
      <h2>Shared State Example</h2>
      
      {/* Connection status */}
      <div className="connection-status">
        Status: {connected ? (
          <span className="connected">Connected</span>
        ) : (
          <span className="disconnected">Disconnected</span>
        )}
        {error && <div className="error">{error}</div>}
      </div>
      
      {/* Message form */}
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Type a message..."
          disabled={!connected}
        />
        <button type="submit" disabled={!connected}>
          Send
        </button>
      </form>
      
      {/* Message count */}
      <div className="message-count">
        Total messages: {state.count}
      </div>
      
      {/* Message list */}
      <ul className="message-list">
        {state.messages.map((message: string, index: number) => (
          <li key={index}>{message}</li>
        ))}
      </ul>
    </div>
  );
};

export default SharedStateExample;
