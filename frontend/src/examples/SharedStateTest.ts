/**
 * Test script for the shared state system
 * This script can be run in a Node.js environment to test the WebSocket connection
 * and shared state functionality without a full React application.
 */

import WebSocket from 'ws';

const WS_URL = 'ws://localhost:8080/ws';
const CLIENT_ID = `test-client-${Date.now()}`;

interface StateMessage {
  type: string;
  data: any;
  timestamp: string;
}

interface StateUpdateEvent {
  type: string;
  state_type: string;
  state_id: string;
  data: Record<string, any>;
  timestamp?: string;
}

const connectWebSocket = (): WebSocket => {
  const url = new URL(WS_URL);
  url.searchParams.append('client_id', CLIENT_ID);
  
  const ws = new WebSocket(url.toString());
  
  ws.on('open', () => {
    console.log('Connected to WebSocket server');
    
    subscribeToState(ws, 'agent', 'agent-1');
    
    startTestSequence(ws);
  });
  
  ws.on('message', (data: WebSocket.Data) => {
    try {
      const message = JSON.parse(data.toString());
      console.log('Received message:', message);
      
      if (message.type === 'message' && message.topic) {
        console.log(`Received update for topic ${message.topic}:`, message.data);
      }
    } catch (err) {
      console.error('Error parsing message:', err);
    }
  });
  
  ws.on('close', () => {
    console.log('Disconnected from WebSocket server');
  });
  
  ws.on('error', (error: Error) => {
    console.error('WebSocket error:', error);
  });
  
  return ws;
};

const subscribeToState = (ws: WebSocket, stateType: string, stateId: string): void => {
  const topic = `state:${stateType}:${stateId}`;
  const message: StateMessage = {
    type: 'subscribe',
    data: topic,
    timestamp: new Date().toISOString()
  };
  
  ws.send(JSON.stringify(message));
  console.log(`Subscribed to ${topic}`);
};

const sendEvent = (ws: WebSocket, eventType: string, eventData: Record<string, any>): void => {
  const message: StateMessage = {
    type: 'event',
    data: {
      type: eventType,
      ...eventData,
      timestamp: new Date().toISOString()
    },
    timestamp: new Date().toISOString()
  };
  
  ws.send(JSON.stringify(message));
  console.log(`Sent ${eventType} event:`, eventData);
};

const updateState = (ws: WebSocket, stateType: string, stateId: string, data: Record<string, any>): void => {
  const message: StateMessage = {
    type: 'event',
    data: {
      type: 'state_update',
      state_type: stateType,
      state_id: stateId,
      data: data
    },
    timestamp: new Date().toISOString()
  };
  
  ws.send(JSON.stringify(message));
  console.log(`Updated ${stateType} state for ${stateId}:`, data);
};

const startTestSequence = (ws: WebSocket): void => {
  console.log('Starting test sequence...');
  
  setTimeout(() => {
    updateState(ws, 'agent', 'agent-1', {
      status: 'initializing',
      progress: 0,
      message: 'Agent initializing...'
    });
  }, 1000);
  
  setTimeout(() => {
    updateState(ws, 'agent', 'agent-1', {
      status: 'running',
      progress: 10,
      message: 'Agent started'
    });
  }, 3000);
  
  setTimeout(() => {
    updateState(ws, 'agent', 'agent-1', {
      status: 'running',
      progress: 30,
      message: 'Processing task...'
    });
  }, 5000);
  
  setTimeout(() => {
    sendEvent(ws, 'custom_event', {
      name: 'test_event',
      value: 'This is a test event'
    });
  }, 7000);
  
  setTimeout(() => {
    updateState(ws, 'agent', 'agent-1', {
      status: 'running',
      progress: 70,
      message: 'Almost done...'
    });
  }, 9000);
  
  setTimeout(() => {
    updateState(ws, 'agent', 'agent-1', {
      status: 'completed',
      progress: 100,
      message: 'Task completed successfully',
      result: {
        id: 'task-123',
        output: 'Task output data'
      }
    });
  }, 11000);
  
  setTimeout(() => {
    console.log('Test sequence completed');
    process.exit(0);
  }, 13000);
};

const ws = connectWebSocket();

process.on('SIGINT', () => {
  console.log('Closing connection...');
  ws.close();
  process.exit(0);
});
