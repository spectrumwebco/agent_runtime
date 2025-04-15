import React, { useState, useEffect } from 'react';
import { useSharedState } from '../hooks/useSharedState';
import { SharedStateProvider } from '../contexts/SharedStateContext';

interface TestState {
  count: number;
  message: string;
  lastUpdate: string;
}

const SharedStateTest: React.FC = () => {
  const { state, loading, error, updateState } = useSharedState<TestState>(
    'test',
    'counter',
    {
      count: 0,
      message: 'Initial message',
      lastUpdate: new Date().toISOString()
    }
  );

  const [localCount, setLocalCount] = useState<number>(0);

  useEffect(() => {
    if (state?.count !== undefined) {
      setLocalCount(state.count);
    }
  }, [state?.count]);

  const handleIncrement = () => {
    updateState({
      count: (state?.count || 0) + 1,
      lastUpdate: new Date().toISOString()
    });
  };

  const handleUpdateMessage = () => {
    updateState({
      message: `Updated at ${new Date().toLocaleTimeString()}`,
      lastUpdate: new Date().toISOString()
    });
  };

  if (loading) {
    return <div className="p-4 bg-gray-100 rounded-lg">Loading shared state...</div>;
  }

  if (error) {
    return <div className="p-4 bg-red-100 text-red-700 rounded-lg">Error: {error.message}</div>;
  }

  return (
    <div className="p-6 max-w-md mx-auto bg-white rounded-xl shadow-md">
      <h2 className="text-xl font-bold mb-4">Shared State Test</h2>
      
      <div className="mb-4 p-4 bg-gray-50 rounded-lg">
        <h3 className="font-semibold mb-2">Current State:</h3>
        <pre className="bg-gray-100 p-2 rounded text-sm">
          {JSON.stringify(state, null, 2)}
        </pre>
      </div>
      
      <div className="flex space-x-4 mb-4">
        <button 
          onClick={handleIncrement}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Increment Count
        </button>
        
        <button 
          onClick={handleUpdateMessage}
          className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        >
          Update Message
        </button>
      </div>
      
      <div className="mt-4 text-sm text-gray-500">
        Last updated: {state?.lastUpdate ? new Date(state.lastUpdate).toLocaleString() : 'Never'}
      </div>
    </div>
  );
};

const SharedStateTestWrapper: React.FC = () => {
  return (
    <SharedStateProvider serverUrl="ws://localhost:8080/ws">
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold mb-6">Shared State System Demo</h1>
        <SharedStateTest />
      </div>
    </SharedStateProvider>
  );
};

export default SharedStateTestWrapper;
