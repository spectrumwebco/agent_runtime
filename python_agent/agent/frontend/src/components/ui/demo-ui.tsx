import React from "react";
import { Button } from "./button";
import { ThemeToggle } from "./theme-toggle";
import { ThemeProvider } from "./theme-provider";
import { Spotlight } from "./aceternity/spotlight";
import { Meteors } from "./aceternity/meteors";
import { TextReveal } from "./aceternity/text-reveal";
import { CardHoverEffect } from "./aceternity/card-hover";

interface DemoUIProps {
  className?: string;
}

export const DemoUI: React.FC<DemoUIProps> = ({ className = "" }) => {
  const cardItems = [
    {
      title: "Agent Runtime",
      description: "Autonomous software engineering agent system with multi-cloud Kubernetes infrastructure.",
      icon: <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z" />
      </svg>,
    },
    {
      title: "Shared State System",
      description: "Real-time state synchronization between frontend and backend components.",
      icon: <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" />
      </svg>,
    },
    {
      title: "WebSocket Communication",
      description: "Two-way communication between Go/Python backend and React frontend.",
      icon: <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
      </svg>,
    },
  ];

  return (
    <ThemeProvider defaultTheme="system" storageKey="ui-theme">
      <div className={`p-4 ${className}`}>
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-bold text-emerald-500">Spectrum Web Co</h1>
          <ThemeToggle />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
          <div className="rounded-lg overflow-hidden">
            <Spotlight className="h-[40vh]">
              <div className="flex flex-col items-center justify-center h-full">
                <TextReveal text="Agent Runtime" className="text-3xl font-bold mb-4" />
                <p className="text-center text-gray-500 dark:text-gray-400 max-w-md">
                  Autonomous software engineering agent system with real-time state synchronization.
                </p>
                <div className="mt-8">
                  <Button variant="emerald">Get Started</Button>
                </div>
              </div>
            </Spotlight>
          </div>

          <div className="rounded-lg overflow-hidden relative">
            <div className="bg-gray-900 h-[40vh] flex items-center justify-center">
              <div className="text-center z-10 relative">
                <h2 className="text-3xl font-bold text-white mb-4">WebSocket Communication</h2>
                <p className="text-gray-300 max-w-md mx-auto mb-8">
                  Two-way communication between Go/Python backend and React frontend.
                </p>
                <Button variant="outline" className="bg-transparent border-white text-white hover:bg-white hover:text-black">
                  Learn More
                </Button>
              </div>
              <Meteors number={20} />
            </div>
          </div>
        </div>

        <div className="mb-8">
          <h2 className="text-2xl font-bold mb-4">Features</h2>
          <CardHoverEffect items={cardItems} />
        </div>
      </div>
    </ThemeProvider>
  );
};
