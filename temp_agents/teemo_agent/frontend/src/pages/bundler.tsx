import React, { useEffect, useState } from 'react';
import { startServer } from '@codesandbox/sandpack-bundler';

export default function BundlerPage() {
  const [status, setStatus] = useState<'starting' | 'running' | 'error'>('starting');
  const [error, setError] = useState<string | null>(null);
  const [port, setPort] = useState<number>(8008);

  useEffect(() => {
    const startBundler = async () => {
      try {
        const configModule = await import('../../sandpack.config.js');
        const config = configModule.default || {};
        
        const bundlerPort = config.bundler?.port || 8008;
        setPort(bundlerPort);
        
        await startServer({
          port: bundlerPort,
          host: config.bundler?.host || '0.0.0.0',
          maxSizeMB: config.bundler?.maxSizeMB || 500,
          timeout: config.bundler?.timeout || 60000,
          allowedOrigins: config.security?.allowedOrigins || ['http://localhost:3000'],
          logLevel: config.logging?.level || 'info',
        });
        
        setStatus('running');
      } catch (err) {
        console.error('Failed to start bundler:', err);
        setStatus('error');
        setError(err instanceof Error ? err.message : String(err));
      }
    };

    startBundler();
  }, []);

  return (
    <div className="container">
      <main className="main">
        <h1 className="title">
          Teemo <span className="highlight">Sandpack Bundler</span>
        </h1>
        
        <div className="status-container">
          {status === 'starting' && (
            <div className="status starting">
              <p>Starting Sandpack bundler...</p>
            </div>
          )}
          
          {status === 'running' && (
            <div className="status running">
              <p>Sandpack bundler is running on port {port}</p>
              <p>URL: <code>http://localhost:{port}</code></p>
            </div>
          )}
          
          {status === 'error' && (
            <div className="status error">
              <p>Failed to start Sandpack bundler</p>
              {error && <p className="error-message">{error}</p>}
            </div>
          )}
        </div>
      </main>

      <style jsx>{`
        .container {
          min-height: 100vh;
          padding: 0 0.5rem;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
          background-color: #1a1a1a;
          color: #f5f5f5;
        }
        
        .main {
          padding: 2rem 0;
          flex: 1;
          display: flex;
          flex-direction: column;
          justify-content: flex-start;
          align-items: center;
          width: 100%;
          max-width: 800px;
        }
        
        .title {
          margin: 0;
          line-height: 1.15;
          font-size: 3rem;
          text-align: center;
        }
        
        .highlight {
          color: #10b981; /* emerald-500 */
        }
        
        .status-container {
          margin-top: 2rem;
          width: 100%;
          padding: 1.5rem;
          border-radius: 0.5rem;
          background-color: #2a2a2a;
        }
        
        .status {
          padding: 1rem;
          border-radius: 0.25rem;
        }
        
        .starting {
          background-color: #374151; /* gray-700 */
        }
        
        .running {
          background-color: #065f46; /* emerald-800 */
        }
        
        .error {
          background-color: #7f1d1d; /* red-800 */
        }
        
        .error-message {
          margin-top: 0.5rem;
          padding: 0.5rem;
          background-color: #1f2937; /* gray-800 */
          border-radius: 0.25rem;
          font-family: monospace;
          white-space: pre-wrap;
        }
        
        code {
          background-color: #1f2937; /* gray-800 */
          padding: 0.2rem 0.4rem;
          border-radius: 0.25rem;
          font-family: monospace;
        }
      `}</style>
    </div>
  );
}
