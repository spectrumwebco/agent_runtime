/**
 * SharedStateClient - WebSocket client for shared state between backend and frontend
 * Provides real-time communication with the backend shared state system
 */

interface WebSocketMessage {
  type: string;
  data: any;
  timestamp: string;
}

type StateCallback = (data: any) => void;

class SharedStateClient {
  private serverUrl: string;
  private socket: WebSocket | null;
  private connected: boolean;
  private reconnectAttempts: number;
  private maxReconnectAttempts: number;
  private reconnectDelay: number;
  private subscriptions: Map<string, Set<StateCallback>>;
  private messageQueue: WebSocketMessage[];
  private clientId: string;

  /**
   * Create a new SharedStateClient
   * @param serverUrl - The URL of the WebSocket server
   */
  constructor(serverUrl: string = 'ws://localhost:8080/ws') {
    this.serverUrl = serverUrl;
    this.socket = null;
    this.connected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10;
    this.reconnectDelay = 1000; // Start with 1 second delay
    this.subscriptions = new Map<string, Set<StateCallback>>(); // topic -> Set of callbacks
    this.messageQueue = []; // Queue of messages to send when reconnected
    this.clientId = `client-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  /**
   * Connect to the WebSocket server
   * @returns Promise that resolves when connected, rejects on error
   */
  connect(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      if (this.connected && this.socket) {
        resolve();
        return;
      }

      const url = new URL(this.serverUrl);
      url.searchParams.append('client_id', this.clientId);

      this.socket = new WebSocket(url.toString());

      this.socket.onopen = () => {
        console.log('WebSocket connected');
        this.connected = true;
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;

        for (const [topic, callbacks] of this.subscriptions.entries()) {
          if (callbacks.size > 0) {
            this.sendMessage('subscribe', topic);
          }
        }

        while (this.messageQueue.length > 0) {
          const message = this.messageQueue.shift();
          if (message && this.socket) {
            this.socket.send(JSON.stringify(message));
          }
        }

        resolve();
      };

      this.socket.onclose = (event: CloseEvent) => {
        console.log(`WebSocket closed: ${event.code} ${event.reason}`);
        this.connected = false;
        this.socket = null;

        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          const delay = Math.min(30000, this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts));
          console.log(`Reconnecting in ${delay}ms...`);
          setTimeout(() => {
            this.reconnectAttempts++;
            this.connect().catch(err => {
              console.error('Reconnection failed:', err);
            });
          }, delay);
        } else {
          console.error('Max reconnection attempts reached');
          reject(new Error('Max reconnection attempts reached'));
        }
      };

      this.socket.onerror = (error: Event) => {
        console.error('WebSocket error:', error);
        reject(error);
      };

      this.socket.onmessage = (event: MessageEvent) => {
        try {
          const message = JSON.parse(event.data);
          if (message.type === 'message' && message.topic) {
            const callbacks = this.subscriptions.get(message.topic);
            if (callbacks) {
              callbacks.forEach(callback => {
                try {
                  callback(message.data);
                } catch (err) {
                  console.error(`Error in subscription callback for topic ${message.topic}:`, err);
                }
              });
            }
          }
        } catch (err) {
          console.error('Error processing message:', err);
        }
      };
    });
  }

  /**
   * Send a message to the WebSocket server
   * @param type - The message type
   * @param data - The message data
   */
  sendMessage(type: string, data: any): void {
    const message: WebSocketMessage = {
      type,
      data,
      timestamp: new Date().toISOString()
    };

    if (this.connected && this.socket) {
      this.socket.send(JSON.stringify(message));
    } else {
      this.messageQueue.push(message);
      if (!this.socket) {
        this.connect().catch(err => {
          console.error('Connection failed:', err);
        });
      }
    }
  }

  /**
   * Subscribe to a topic
   * @param topic - The topic to subscribe to
   * @param callback - The callback to call when a message is received
   */
  subscribe(topic: string, callback: StateCallback): void {
    if (!this.subscriptions.has(topic)) {
      this.subscriptions.set(topic, new Set<StateCallback>());
      if (this.connected) {
        this.sendMessage('subscribe', topic);
      }
    }

    const callbacks = this.subscriptions.get(topic);
    if (callbacks) {
      callbacks.add(callback);
    }
  }

  /**
   * Unsubscribe from a topic
   * @param topic - The topic to unsubscribe from
   * @param callback - The callback to remove (optional)
   */
  unsubscribe(topic: string, callback: StateCallback | null = null): void {
    if (this.subscriptions.has(topic)) {
      if (callback) {
        const callbacks = this.subscriptions.get(topic);
        if (callbacks) {
          callbacks.delete(callback);
          if (callbacks.size === 0) {
            this.subscriptions.delete(topic);
            if (this.connected) {
              this.sendMessage('unsubscribe', topic);
            }
          }
        }
      } else {
        this.subscriptions.delete(topic);
        if (this.connected) {
          this.sendMessage('unsubscribe', topic);
        }
      }
    }
  }

  /**
   * Subscribe to state updates
   * @param stateType - The type of state
   * @param stateId - The ID of the state
   * @param callback - The callback to call when state is updated
   */
  subscribeToState(stateType: string, stateId: string, callback: StateCallback): void {
    const topic = `state:${stateType}:${stateId}`;
    this.subscribe(topic, callback);
  }

  /**
   * Unsubscribe from state updates
   * @param stateType - The type of state
   * @param stateId - The ID of the state
   * @param callback - The callback to remove (optional)
   */
  unsubscribeFromState(stateType: string, stateId: string, callback: StateCallback | null = null): void {
    const topic = `state:${stateType}:${stateId}`;
    this.unsubscribe(topic, callback);
  }

  /**
   * Send an event to the server
   * @param eventData - The event data
   */
  sendEvent(eventData: Record<string, any>): void {
    this.sendMessage('event', eventData);
  }

  /**
   * Update state on the server
   * @param stateType - The type of state
   * @param stateId - The ID of the state
   * @param data - The state data
   */
  updateState(stateType: string, stateId: string, data: Record<string, any>): void {
    this.sendEvent({
      type: 'state_update',
      state_type: stateType,
      state_id: stateId,
      data
    });
  }

  /**
   * Close the connection
   */
  close(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      this.connected = false;
    }
  }
}

export default SharedStateClient;
