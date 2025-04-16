import React from "react";
import ReactDOM from "react-dom/client";
import { KledApp } from "./components/KledApp";
import "./styles/globals.css";
import "./styles/aceternity-ui.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <KledApp />
  </React.StrictMode>
);
