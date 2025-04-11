import React, { useState } from "react";
import Footer from "./components/Footer";
import Header from "./components/Header";
import Run from "./Run";
import "./static/font.css";
import "./static/index.css";
import "./static/devin-layout.css";
import { SharedStateProvider } from "./contexts/SharedStateContext";
import { BrowserRouter, Routes, Route } from "react-router-dom";

const App: React.FC = () => {
  return (
    <SharedStateProvider serverUrl="ws://localhost:8080/ws">
      <BrowserRouter>
        <div className="app-container">
          <Header />
          <Routes>
            <Route path="/" element={<Run />} />
          </Routes>
          <Footer />
          <div className="fixed bottom-4 right-4 z-50">
            <button 
              className="bg-emerald-500 hover:bg-emerald-600 text-white px-4 py-2 rounded-md shadow-md"
            >
              Theme Toggle
            </button>
          </div>
        </div>
      </BrowserRouter>
    </SharedStateProvider>
  );
};

export default App;
