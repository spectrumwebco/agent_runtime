import React, { useState } from 'react';
import SandpackEditor from '../../components/SandpackEditor';

export default function Home() {
  const [code, setCode] = useState('// Write your code here');
  const [language, setLanguage] = useState('javascript');
  const [output, setOutput] = useState('');
  const [isExecuting, setIsExecuting] = useState(false);

  const handleCodeChange = (newCode: string) => {
    setCode(newCode);
  };

  const handleExecute = async (codeToExecute: string) => {
    setIsExecuting(true);
    setOutput('Executing code...');
    
    try {
      const response = await fetch('/api/execute-code', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          code: codeToExecute,
          language,
        }),
      });
      
      const data = await response.json();
      
      if (data.success) {
        setOutput(data.output);
      } else {
        setOutput(`Error: ${data.error}`);
      }
    } catch (error) {
      setOutput(`Error: ${error instanceof Error ? error.message : String(error)}`);
    } finally {
      setIsExecuting(false);
    }
  };

  return (
    <div className="container">
      <main className="main">
        <h1 className="title">
          Teemo <span className="highlight">Code Editor</span>
        </h1>
        
        <div className="description">
          <p>
            Write and execute code in multiple programming languages using the LibreChat Code Interpreter API.
          </p>
        </div>
        
        <div className="language-selector">
          <label htmlFor="language-select">Language:</label>
          <select 
            id="language-select"
            value={language}
            onChange={(e) => setLanguage(e.target.value)}
          >
            <option value="javascript">JavaScript</option>
            <option value="typescript">TypeScript</option>
            <option value="python">Python</option>
            <option value="c++">C++</option>
            <option value="c#">C#</option>
            <option value="go">Go</option>
            <option value="rust">Rust</option>
            <option value="php">PHP</option>
          </select>
        </div>
        
        <div className="editor-container">
          <SandpackEditor
            initialCode={code}
            language={language}
            theme="dark"
            showPreview={['javascript', 'typescript'].includes(language)}
            showConsole={true}
            onCodeChange={handleCodeChange}
            onExecute={handleExecute}
          />
        </div>
        
        <div className="output-container">
          <h2>Output</h2>
          <div className="output">
            {isExecuting ? (
              <div className="loading">Executing code...</div>
            ) : (
              <pre>{output}</pre>
            )}
          </div>
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
          max-width: 1200px;
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
        
        .description {
          text-align: center;
          line-height: 1.5;
          font-size: 1.5rem;
          margin: 1rem 0;
        }
        
        .language-selector {
          margin: 1rem 0;
          display: flex;
          align-items: center;
          gap: 0.5rem;
        }
        
        .language-selector select {
          padding: 0.5rem;
          border-radius: 0.25rem;
          background-color: #2a2a2a;
          color: #f5f5f5;
          border: 1px solid #444;
        }
        
        .editor-container {
          width: 100%;
          height: 400px;
          margin: 1rem 0;
          border-radius: 0.5rem;
          overflow: hidden;
        }
        
        .output-container {
          width: 100%;
          margin: 1rem 0;
        }
        
        .output {
          background-color: #2a2a2a;
          padding: 1rem;
          border-radius: 0.5rem;
          min-height: 100px;
          max-height: 300px;
          overflow: auto;
        }
        
        .loading {
          display: flex;
          justify-content: center;
          align-items: center;
          height: 100px;
        }
      `}</style>
    </div>
  );
}
