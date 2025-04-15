import React from "react";
import { ThemeProvider } from "./theme-provider";
import { DemoUI } from "./demo-ui";

export const AppDemo: React.FC = () => {
  return (
    <ThemeProvider defaultTheme="system" storageKey="ui-theme">
      <div className="min-h-screen bg-background text-foreground">
        <DemoUI />
      </div>
    </ThemeProvider>
  );
};
