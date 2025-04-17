import type { NextApiRequest, NextApiResponse } from 'next';
import path from 'path';
import fs from 'fs';

interface ExecuteCodeResponse {
  success: boolean;
  output?: string;
  error?: string;
  execution_time?: number;
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<ExecuteCodeResponse>
) {
  if (req.method !== 'POST') {
    return res.status(405).json({ success: false, error: 'Method not allowed' });
  }

  try {
    const { code, language, dependencies } = req.body;

    if (!code || !language) {
      return res.status(400).json({ 
        success: false, 
        error: 'Code and language are required' 
      });
    }

    const apiKey = process.env.LIBRECHAT_CODE_API_KEY;
    if (!apiKey) {
      return res.status(500).json({ 
        success: false, 
        error: 'LibreChat API key not configured' 
      });
    }

    const response = await fetch('http://185.192.220.224:8000/api/v1/execute', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-api-key': apiKey
      },
      body: JSON.stringify({
        code,
        language,
        timeout: 30
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      return res.status(response.status).json({
        success: false,
        error: errorData.message || 'Error executing code'
      });
    }

    const data = await response.json();
    return res.status(200).json({
      success: true,
      output: data.output,
      execution_time: data.execution_time
    });
  } catch (error) {
    console.error('Error executing code:', error);
    return res.status(500).json({
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    });
  }
}
