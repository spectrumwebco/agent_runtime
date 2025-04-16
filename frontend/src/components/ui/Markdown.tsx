import React from 'react';

interface MarkdownProps {
  content: string;
  className?: string;
}

export const Markdown: React.FC<MarkdownProps> = ({
  content,
  className = '',
}) => {
  
  const processMarkdown = (markdown: string) => {
    let processed = markdown
      .replace(/^# (.*$)/gm, '<h1>$1</h1>')
      .replace(/^## (.*$)/gm, '<h2>$1</h2>')
      .replace(/^### (.*$)/gm, '<h3>$1</h3>')
      .replace(/^#### (.*$)/gm, '<h4>$1</h4>')
      .replace(/^##### (.*$)/gm, '<h5>$1</h5>')
      .replace(/^###### (.*$)/gm, '<h6>$1</h6>');
    
    processed = processed
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/__(.*?)__/g, '<strong>$1</strong>')
      .replace(/_(.*?)_/g, '<em>$1</em>');
    
    processed = processed.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" class="text-emerald-500 hover:underline" target="_blank" rel="noopener noreferrer">$1</a>');
    
    processed = processed.replace(/```([\s\S]*?)```/g, '<pre class="bg-gray-900 p-4 rounded-md overflow-x-auto my-4"><code>$1</code></pre>');
    
    processed = processed.replace(/`(.*?)`/g, '<code class="bg-gray-200 dark:bg-gray-800 px-1 py-0.5 rounded text-sm">$1</code>');
    
    processed = processed.replace(/^\s*\*\s(.*$)/gm, '<li>$1</li>');
    processed = processed.replace(/^\s*-\s(.*$)/gm, '<li>$1</li>');
    processed = processed.replace(/^\s*\d\.\s(.*$)/gm, '<li>$1</li>');
    
    processed = processed.replace(/^([^<].*)\n$/gm, '<p>$1</p>');
    
    processed = processed.replace(/\n/g, '<br>');
    
    return processed;
  };
  
  return (
    <div 
      className={`prose dark:prose-invert max-w-none ${className}`}
      dangerouslySetInnerHTML={{ __html: processMarkdown(content) }}
    />
  );
};

export default Markdown;
