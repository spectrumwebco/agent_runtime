import React from "react";
import { createRoot } from "react-dom/client";
import "tailwindcss/tailwind.css";
import App from "./App";
import reportWebVitals from "./reportWebVitals";

const container = document.getElementById("root");
if (!container) {
  throw new Error("Failed to find the root element");
}

const root = createRoot(container);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

reportWebVitals();
