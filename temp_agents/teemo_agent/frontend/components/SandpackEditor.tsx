import React, { useState, useEffect } from 'react';
import {
  SandpackProvider,
  SandpackLayout,
  SandpackCodeEditor,
  SandpackPreview,
  SandpackConsole,
  useSandpack
} from '@codesandbox/sandpack-react';
import { sandpackDark } from '@codesandbox/sandpack-themes';

interface SandpackEditorProps {
  initialCode?: string;
  language?: string;
  theme?: 'light' | 'dark';
  showPreview?: boolean;
  showConsole?: boolean;
  onCodeChange?: (code: string) => void;
  onExecute?: (code: string) => void;
}

const getFileExtension = (language: string): string => {
  const extensions: Record<string, string> = {
    javascript: 'js',
    typescript: 'ts',
    react: 'tsx',
    vue: 'vue',
    angular: 'ts',
    svelte: 'svelte',
    html: 'html',
    css: 'css',
    python: 'py',
    cpp: 'cpp',
    csharp: 'cs',
    go: 'go',
    rust: 'rs',
    php: 'php',
  };
  
  return extensions[language] || 'js';
};

const getSetupFiles = (language: string, code: string) => {
  const ext = getFileExtension(language);
  
  if (['javascript', 'typescript', 'react'].includes(language)) {
    return {
      [`/index.${ext}`]: { code },
    };
  } else if (language === 'vue') {
    return {
      '/App.vue': { code },
      '/main.js': { 
        code: `import { createApp } from 'vue';
import App from './App.vue';
createApp(App).mount('#app');` 
      },
    };
  } else {
    return {
      [`/index.${ext}`]: { code },
    };
  }
};

export const SandpackEditor: React.FC<SandpackEditorProps> = ({
  initialCode = '',
  language = 'javascript',
  theme = 'dark',
  showPreview = true,
  showConsole = true,
  onCodeChange,
  onExecute,
}) => {
  const [code, setCode] = useState(initialCode);
  
  const handleCodeChange = (newCode: string) => {
    setCode(newCode);
    if (onCodeChange) {
      onCodeChange(newCode);
    }
  };
  
  const ExecuteButton = () => {
    const { sandpack } = useSandpack();
    
    const handleExecute = () => {
      const currentCode = sandpack.files[`/index.${getFileExtension(language)}`]?.code || code;
      if (onExecute) {
        onExecute(currentCode);
      }
    };
    
    return (
      <button 
        onClick={handleExecute}
        className="execute-button"
        style={{
          backgroundColor: '#10b981', // emerald-500
          color: 'white',
          border: 'none',
          borderRadius: '0.375rem',
          padding: '0.5rem 1rem',
          cursor: 'pointer',
          fontWeight: 'bold',
        }}
      >
        Execute Code
      </button>
    );
  };
  
  return (
    <div className="sandpack-editor-container" style={{ width: '100%', height: '100%' }}>
      <SandpackProvider
        theme={theme === 'dark' ? sandpackDark : undefined}
        files={getSetupFiles(language, code)}
        template={language === 'react' ? 'react' : 'vanilla'}
        customSetup={{
          dependencies: language === 'vue' ? { vue: 'latest' } : {},
        }}
      >
        <SandpackLayout>
          <div className="editor-header" style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            padding: '0.5rem',
            backgroundColor: theme === 'dark' ? '#1e1e1e' : '#f5f5f5',
            borderBottom: '1px solid #333'
          }}>
            <span className="language-indicator" style={{ 
              fontWeight: 'bold',
              color: theme === 'dark' ? '#10b981' : '#047857'
            }}>
              {language.toUpperCase()}
            </span>
            <ExecuteButton />
          </div>
          <SandpackCodeEditor 
            showLineNumbers={true}
            showInlineErrors={true}
            wrapContent={true}
            onChange={handleCodeChange}
          />
          {showPreview && ['javascript', 'typescript', 'react', 'vue', 'html', 'css'].includes(language) && (
            <SandpackPreview />
          )}
          {showConsole && (
            <SandpackConsole />
          )}
        </SandpackLayout>
      </SandpackProvider>
    </div>
  );
};

export default SandpackEditor;
