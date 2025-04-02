import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import React from "react";

const Root: React.FC = () => {
    return (
      <StrictMode>
        <App />
      </StrictMode>
    );
};

createRoot(document.getElementById("root")!).render(<Root />);
