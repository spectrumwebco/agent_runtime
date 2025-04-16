/**
 * v0 Integration Script for Kled.io
 * 
 * This script integrates v0 with the Devin API for frontend development.
 * It provides hooks for v0 to communicate with the Devin API during
 * frontend development.
 */

import * as fs from 'fs';
import * as path from 'path';
import * as https from 'https';
import { URL } from 'url';

interface Config {
  devinApiUrl: string;
  devinApiKey: string | undefined;
  v0Email: string;
  sessionId: string | null;
}

interface ApiResponse {
  status: string;
  message?: string;
}

const config: Config = {
  devinApiUrl: process.env.DEVIN_API_URL || 'http://185.192.220.224:8000',
  devinApiKey: process.env.DEVIN_API_KEY,
  v0Email: process.env.V0_EMAIL || 'ove.govender@gmail.com',
  sessionId: null
};

try {
  const homeDir = process.env.HOME || process.env.USERPROFILE;
  if (homeDir) {
    const sessionIdPath = path.join(homeDir, '.devin_session_id');
    if (fs.existsSync(sessionIdPath)) {
      config.sessionId = fs.readFileSync(sessionIdPath, 'utf8').trim();
    }
  }
} catch (error) {
  console.error('Error loading session ID:', error);
}

/**
 * Call the Devin API
 * @param endpoint - API endpoint
 * @param method - HTTP method
 * @param data - Request data
 * @returns Promise<ApiResponse> - Response data
 */
async function callDevinApi(endpoint: string, method: string, data?: any): Promise<ApiResponse> {
  return new Promise((resolve, reject) => {
    const url = new URL(endpoint, config.devinApiUrl);
    const options = {
      method: method,
      headers: {
        'Content-Type': 'application/json',
        'x-api-key': config.devinApiKey || ''
      }
    };
    
    const req = https.request(url, options, (res) => {
      let responseData = '';
      
      res.on('data', (chunk) => {
        responseData += chunk;
      });
      
      res.on('end', () => {
        try {
          const parsedData = JSON.parse(responseData) as ApiResponse;
          resolve(parsedData);
        } catch (error) {
          reject(new Error(`Failed to parse response: ${responseData}`));
        }
      });
    });
    
    req.on('error', (error) => {
      reject(error);
    });
    
    if (data) {
      req.write(JSON.stringify(data));
    }
    
    req.end();
  });
}

/**
 * Initialize a session with the Devin API
 * @returns Promise<string> - Session ID
 */
async function initializeSession(): Promise<string> {
  if (config.sessionId) {
    console.log(`Using existing session: ${config.sessionId}`);
    return config.sessionId;
  }
  
  const timestamp = Math.floor(Date.now() / 1000);
  const data = {
    session_id: generateUUID(),
    repo_name: 'agent_runtime',
    timestamp: timestamp,
    event: 'v0_frontend_start',
    status: 'running',
    metadata: {
      integration_version: '1.0',
      environment: 'development',
      user: config.v0Email
    }
  };
  
  try {
    const response = await callDevinApi('/api/v1/sessions', 'POST', data);
    
    if (response.status === 'success') {
      config.sessionId = data.session_id;
      
      const homeDir = process.env.HOME || process.env.USERPROFILE;
      if (homeDir) {
        const sessionIdPath = path.join(homeDir, '.devin_session_id');
        fs.writeFileSync(sessionIdPath, config.sessionId);
      }
      
      console.log(`Session initialized: ${config.sessionId}`);
      return config.sessionId;
    } else {
      throw new Error(`Failed to initialize session: ${response.message || 'Unknown error'}`);
    }
  } catch (error) {
    console.error('Error initializing session:', error);
    throw error;
  }
}

/**
 * Update session status
 * @param status - Status
 * @param message - Status message
 * @returns Promise<boolean> - Success status
 */
async function updateSessionStatus(status: string, message: string): Promise<boolean> {
  if (!config.sessionId) {
    console.error('No active session');
    return false;
  }
  
  const timestamp = Math.floor(Date.now() / 1000);
  const data = {
    timestamp: timestamp,
    status: status,
    message: message
  };
  
  try {
    const response = await callDevinApi(`/api/v1/sessions/${config.sessionId}/status`, 'POST', data);
    
    if (response.status === 'success') {
      console.log(`Session status updated: ${status} - ${message}`);
      return true;
    } else {
      console.error(`Failed to update session status: ${response.message || 'Unknown error'}`);
      return false;
    }
  } catch (error) {
    console.error('Error updating session status:', error);
    return false;
  }
}

/**
 * End session
 * @param status - Status
 * @returns Promise<boolean> - Success status
 */
async function endSession(status: string): Promise<boolean> {
  if (!config.sessionId) {
    console.error('No active session');
    return false;
  }
  
  const timestamp = Math.floor(Date.now() / 1000);
  const data = {
    timestamp: timestamp,
    event: 'v0_frontend_end',
    status: status,
    metadata: {
      integration_version: '1.0',
      environment: 'development',
      user: config.v0Email
    }
  };
  
  try {
    const response = await callDevinApi(`/api/v1/sessions/${config.sessionId}`, 'POST', data);
    
    if (response.status === 'success') {
      console.log(`Session ended: ${config.sessionId} with status: ${status}`);
      
      const homeDir = process.env.HOME || process.env.USERPROFILE;
      if (homeDir) {
        const sessionIdPath = path.join(homeDir, '.devin_session_id');
        if (fs.existsSync(sessionIdPath)) {
          fs.unlinkSync(sessionIdPath);
        }
      }
      
      config.sessionId = null;
      return true;
    } else {
      console.error(`Failed to end session: ${response.message || 'Unknown error'}`);
      return false;
    }
  } catch (error) {
    console.error('Error ending session:', error);
    return false;
  }
}

/**
 * Generate UUID
 * @returns string - UUID
 */
function generateUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

export {
  initializeSession,
  updateSessionStatus,
  endSession
};
