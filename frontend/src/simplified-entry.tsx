import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';

const App = () => {
  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <header className="bg-gray-800 border-b border-gray-700">
        <div className="container mx-auto px-4 py-3 flex justify-between items-center">
          <div className="flex items-center">
            <h1 className="text-xl font-bold text-emerald-500">Kled.io</h1>
            <span className="ml-2 text-sm text-gray-400">Agent Runtime</span>
          </div>
        </div>
      </header>
      
      <main className="container mx-auto p-4">
        <div className="grid grid-cols-1 gap-4">
          <div className="border rounded-lg p-4 bg-gray-800">
            <h3 className="text-lg font-medium mb-4">CoPilot Integration</h3>
            <p className="text-gray-400">AI-powered agent coordination</p>
          </div>
        </div>
      </main>
      
      <footer className="bg-gray-800 border-t border-gray-700 py-3 mt-8">
        <div className="container mx-auto px-4 text-center text-sm text-gray-400">
          &copy; {new Date().getFullYear()} Kled.io - Agent Runtime System
        </div>
      </footer>
    </div>
  );
};

const root = document.getElementById('root');
if (root) {
  ReactDOM.createRoot(root).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
}
