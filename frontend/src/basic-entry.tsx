import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';

const App = () => {
  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <header className="bg-gray-800 border-b border-gray-700 p-4">
        <div className="container mx-auto flex justify-between items-center">
          <div className="flex items-center">
            <h1 className="text-xl font-bold text-emerald-500">Kled.io</h1>
            <span className="ml-2 text-sm text-gray-400">Agent Runtime</span>
          </div>
        </div>
      </header>
      
      <main className="container mx-auto p-4">
        <div className="grid grid-cols-1 gap-6 mb-6">
          <div className="border rounded-lg p-6 bg-gray-800">
            <h3 className="text-lg font-medium mb-4 text-emerald-500">CoPilot Integration</h3>
            <p className="text-gray-400 mb-4">AI-powered agent coordination</p>
            <button className="px-4 py-2 bg-emerald-500 hover:bg-emerald-600 rounded-md text-white">
              Connect Agent
            </button>
          </div>
        </div>
        
        <div className="border rounded-lg p-6 bg-gray-800">
          <h3 className="text-lg font-medium mb-4 text-emerald-500">Agent View Switcher</h3>
          <p className="text-gray-400 mb-4">Toggle between Control Plane Agent and Worker Agents</p>
          <div className="flex gap-4">
            <button className="px-4 py-2 bg-emerald-500 hover:bg-emerald-600 rounded-md text-white">
              Control Plane
            </button>
            <button className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md text-white">
              Worker Agents
            </button>
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
