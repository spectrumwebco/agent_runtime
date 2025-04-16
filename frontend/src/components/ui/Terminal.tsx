import React, { useEffect, useRef } from 'react';

interface TerminalProps {
  content: string;
  title?: string;
  className?: string;
  autoScroll?: boolean;
}

export const Terminal: React.FC<TerminalProps> = ({
  content,
  title = 'Terminal',
  className = '',
  autoScroll = true,
}) => {
  const terminalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (autoScroll && terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [content, autoScroll]);

  return (
    <div className={`rounded-lg overflow-hidden shadow-md ${className}`}>
      <div className="flex items-center justify-between px-4 py-2 bg-gray-800 dark:bg-gray-900">
        <div className="flex items-center space-x-2">
          <div className="flex space-x-1">
            <div className="w-3 h-3 rounded-full bg-red-500"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
            <div className="w-3 h-3 rounded-full bg-green-500"></div>
          </div>
          <span className="text-sm font-medium text-gray-200">{title}</span>
        </div>
      </div>
      <div
        ref={terminalRef}
        className="p-4 bg-black text-gray-200 font-mono text-sm overflow-auto"
        style={{ maxHeight: '400px' }}
      >
        <pre className="whitespace-pre-wrap">{content}</pre>
      </div>
    </div>
  );
};

export default Terminal;
