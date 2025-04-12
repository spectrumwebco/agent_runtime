/**
 * SharedStateClient - WebSocket client for shared state between backend and frontend
 * Provides real-time communication with the backend shared state system
 */
class SharedStateClient {
  /**
   * Create a new SharedStateClient
   * @param {string} serverUrl - The URL of the WebSocket server
   */
  constructor(serverUrl = 'ws://localhost:8080/ws') {
    this.serverUrl = serverUrl;
    this.socket = null;
    this.connected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10;
    this.reconnectDelay = 1000; // Start with 1 second delay
    this.subscriptions = new Map(); // topic -> Set of callbacks
    this.messageQueue = []; // Queue of messages to send when reconnected
    this.clientId = `client-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  /**
   * Connect to the WebSocket server
   * @returns {Promise} - Resolves when connected, rejects on error
   */
  connect() {
    return new Promise((resolve, reject) => {
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
          this.socket.send(JSON.stringify(message));
        }

        resolve();
      };

      this.socket.onclose = (event) => {
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

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        reject(error);
      };

      this.socket.onmessage = (event) => {
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
   * @param {string} type - The message type
   * @param {any} data - The message data
   */
  sendMessage(type, data) {
    const message = {
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
   * @param {string} topic - The topic to subscribe to
   * @param {function} callback - The callback to call when a message is received
   */
  subscribe(topic, callback) {
    if (!this.subscriptions.has(topic)) {
      this.subscriptions.set(topic, new Set());
      if (this.connected) {
        this.sendMessage('subscribe', topic);
      }
    }

    const callbacks = this.subscriptions.get(topic);
    callbacks.add(callback);
  }

  /**
   * Unsubscribe from a topic
   * @param {string} topic - The topic to unsubscribe from
   * @param {function} callback - The callback to remove (optional)
   */
  unsubscribe(topic, callback = null) {
    if (this.subscriptions.has(topic)) {
      if (callback) {
        const callbacks = this.subscriptions.get(topic);
        callbacks.delete(callback);
        if (callbacks.size === 0) {
          this.subscriptions.delete(topic);
          if (this.connected) {
            this.sendMessage('unsubscribe', topic);
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
   * @param {string} stateType - The type of state
   * @param {string} stateId - The ID of the state
   * @param {function} callback - The callback to call when state is updated
   */
  subscribeToState(stateType, stateId, callback) {
    const topic = `state:${stateType}:${stateId}`;
    this.subscribe(topic, callback);
  }

  /**
   * Unsubscribe from state updates
   * @param {string} stateType - The type of state
   * @param {string} stateId - The ID of the state
   * @param {function} callback - The callback to remove (optional)
   */
  unsubscribeFromState(stateType, stateId, callback = null) {
    const topic = `state:${stateType}:${stateId}`;
    this.unsubscribe(topic, callback);
  }

  /**
   * Send an event to the server
   * @param {object} eventData - The event data
   */
  sendEvent(eventData) {
    this.sendMessage('event', eventData);
  }

  /**
   * Update state on the server
   * @param {string} stateType - The type of state
   * @param {string} stateId - The ID of the state
   * @param {object} data - The state data
   */
  updateState(stateType, stateId, data) {
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
  close() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      this.connected = false;
    }
  }
}

export default SharedStateClient;
