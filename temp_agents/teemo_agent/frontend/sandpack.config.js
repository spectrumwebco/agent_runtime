/**
 * Sandpack Bundler Configuration for Teemo
 * 
 * This configuration file sets up the self-hosted Sandpack bundler
 * for the Teemo UI/UX Agent.
 */

const path = require('path');

module.exports = {
  bundler: {
    port: 8008,
    
    host: '0.0.0.0',
    
    maxSizeMB: 500,
    
    timeout: 60000,
  },
  
  cache: {
    directory: path.join(__dirname, '.cache'),
    
    maxSizeMB: 1000,
    
    ttl: 1000 * 60 * 60 * 24, // 24 hours
  },
  
  security: {
    allowedOrigins: ['http://localhost:3000'],
    
    maxPayloadSizeMB: 50,
  },
  
  logging: {
    level: 'info',
    
    console: true,
    
    file: path.join(__dirname, 'logs', 'sandpack.log'),
  },
  
  templates: [
    'react',
    'react-ts',
    'vanilla',
    'vanilla-ts',
    'vue',
    'vue-ts',
    'angular',
    'svelte',
    'solid',
  ],
  
  customTemplates: {
    'python': {
      files: {
        '/index.py': { code: '# Python code goes here' },
      },
    },
    'cpp': {
      files: {
        '/main.cpp': { code: '// C++ code goes here' },
      },
    },
    'csharp': {
      files: {
        '/Program.cs': { code: '// C# code goes here' },
      },
    },
    'go': {
      files: {
        '/main.go': { code: '// Go code goes here' },
      },
    },
    'rust': {
      files: {
        '/main.rs': { code: '// Rust code goes here' },
      },
    },
    'php': {
      files: {
        '/index.php': { code: '<?php\n// PHP code goes here\n?>' },
      },
    },
  },
};
