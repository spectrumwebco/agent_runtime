import React, { useState } from "react";
import { cn } from "../../utils/cn";
import { SpotlightCard } from "./aceternity/spotlight-card";
import { GradientButton } from "./aceternity/gradient-button";

interface MagicMCPServerProps {
  className?: string;
  onGenerateUI?: (generatedCode: string) => void;
}

export const MagicMCPServer: React.FC<MagicMCPServerProps> = ({
  className,
  onGenerateUI,
}) => {
  const [prompt, setPrompt] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedCode, setGeneratedCode] = useState("");
  const [uiComponents, setUIComponents] = useState<string[]>([]);

  const componentTemplates = [
    {
      name: "Dashboard Card",
      description: "A card component with spotlight effect and gradient border",
      template: `
import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { SpotlightCard } from "@/components/ui/aceternity/spotlight-card";

export const DashboardCard = ({ title, children }) => {
  return (
    <SpotlightCard>
      <Card className="border-0 bg-transparent">
        <CardHeader>
          <CardTitle>{title}</CardTitle>
        </CardHeader>
        <CardContent>{children}</CardContent>
      </Card>
    </SpotlightCard>
  );
};
      `,
    },
    {
      name: "Animated Button",
      description: "A button with gradient animation and hover effects",
      template: `
import React from "react";
import { Button } from "@/components/ui/button";
import { GradientButton } from "@/components/ui/aceternity/gradient-button";

export const AnimatedButton = ({ children, ...props }) => {
  return (
    <GradientButton {...props}>
      {children}
    </GradientButton>
  );
};
      `,
    },
    {
      name: "Navigation Menu",
      description: "A navigation menu with floating animation",
      template: `
import React from "react";
import { NavigationMenu, NavigationMenuList, NavigationMenuItem, NavigationMenuLink } from "@/components/ui/navigation-menu";
import { FloatingNav } from "@/components/ui/aceternity/floating-navbar";

export const AnimatedNavigation = ({ items }) => {
  const navItems = items.map(item => ({
    name: item.label,
    link: item.href,
    icon: item.icon,
  }));

  return (
    <FloatingNav navItems={navItems} />
  );
};
      `,
    },
  ];

  const handleGenerateUI = () => {
    if (!prompt.trim()) return;
    
    setIsGenerating(true);
    
    setTimeout(() => {
      const randomTemplate = componentTemplates[Math.floor(Math.random() * componentTemplates.length)];
      
      const customizedCode = `
import React from "react";
import { cn } from "@/lib/utils";
import { SpotlightCard } from "@/components/ui/aceternity/spotlight-card";
import { GradientButton } from "@/components/ui/aceternity/gradient-button";

export const CustomComponent = ({ className, ...props }) => {
  return (
    <SpotlightCard className={cn("p-4", className)}>
      <h3 className="text-lg font-medium mb-4">${prompt}</h3>
      <div className="space-y-4">
        <div className="border rounded-md p-3 bg-white/5">
          <p className="text-sm">
            Custom component with Aceternity UI styling and Shadcn UI structure
          </p>
        </div>
        <GradientButton>
          Interact
        </GradientButton>
      </div>
    </SpotlightCard>
  );
};
      `;
      
      setGeneratedCode(customizedCode);
      setUIComponents((prev) => [...prev, prompt]);
      
      if (onGenerateUI) {
        onGenerateUI(customizedCode);
      }
      
      setIsGenerating(false);
      setPrompt("");
    }, 2000);
  };

  return (
    <SpotlightCard className={cn("flex flex-col", className)}>
      <div className="flex justify-between items-center p-4 border-b">
        <div>
          <h3 className="text-lg font-medium">Magic MCP Server</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Generate UI components with Shadcn + Aceternity UI
          </p>
        </div>
        <span className={cn(
          "px-2 py-1 rounded-full text-xs",
          isGenerating ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300" : "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300"
        )}>
          {isGenerating ? "Generating..." : "Ready"}
        </span>
      </div>
      
      <div className="p-4 flex-1 overflow-y-auto">
        <div className="mb-4">
          <label htmlFor="ui-prompt" className="block text-sm font-medium mb-2">
            Describe the UI component you need
          </label>
          <textarea
            id="ui-prompt"
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 bg-white dark:bg-gray-800"
            rows={4}
            placeholder="E.g., Create a dashboard card with user statistics and a gradient border..."
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            disabled={isGenerating}
          />
        </div>
        
        <GradientButton
          onClick={handleGenerateUI}
          disabled={isGenerating || !prompt.trim()}
          className="w-full mb-4"
        >
          {isGenerating ? "Generating UI Component..." : "Generate UI Component"}
        </GradientButton>
        
        {generatedCode && (
          <div className="mt-4">
            <h4 className="text-sm font-medium mb-2">Generated Component</h4>
            <div className="bg-gray-100 dark:bg-gray-800 rounded-md p-3 overflow-x-auto">
              <pre className="text-xs">
                <code>{generatedCode}</code>
              </pre>
            </div>
          </div>
        )}
        
        {uiComponents.length > 0 && (
          <div className="mt-4">
            <h4 className="text-sm font-medium mb-2">Generated Components History</h4>
            <ul className="space-y-2">
              {uiComponents.map((comp, index) => (
                <li key={index} className="text-sm border-l-2 border-emerald-500 pl-2">
                  {comp}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </SpotlightCard>
  );
};

export default MagicMCPServer;
