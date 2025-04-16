import React, { useState } from "react";
import { PlatformProvider } from "../../platform/platform-provider";
import { ElectronWrapper } from "../../platform/electron-wrapper";
import { LynxWrapper } from "../../platform/lynx-wrapper";
import { PlatformSpecific } from "../../platform/platform-specific";
import { GradientCard } from "../aceternity/gradient-card";
import { SpotlightButton } from "../aceternity/spotlight-button";
import { AnimatedBackground } from "../aceternity/animated-background";
import { TextGradient } from "../aceternity/text-gradient";
import { ThreeDCard } from "../aceternity/3d-card";
import { Button } from "../shadcn/button";
import { Card } from "../shadcn/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../shadcn/tabs";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "../shadcn/accordion";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "../radix/dialog";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../radix/dropdown-menu";
import { AgentCard } from "../tailwind/agent-card";
import { AgentDashboard } from "../tailwind/agent-dashboard";
import { AITaskTracker } from "../vercel-ai/ai-task-tracker";
import { AIStateVisualizer } from "../vercel-ai/ai-state-visualizer";
import { ModelSelector } from "../model-selector/model-selector";

export const UIKitDemo: React.FC = () => {
  const [platform, setPlatform] = useState<"web" | "electron" | "mobile">("web");
  const [activeTab, setActiveTab] = useState("radix");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  
  const sampleAgents = [
    {
      id: "agent-1",
      title: "Code Generation Agent",
      description: "Generates code based on natural language descriptions",
      status: "running" as const,
      progress: 65,
    },
    {
      id: "agent-2",
      title: "Data Analysis Agent",
      description: "Analyzes data and generates insights",
      status: "completed" as const,
      progress: 100,
    },
    {
      id: "agent-3",
      title: "Testing Agent",
      description: "Runs tests and reports results",
      status: "idle" as const,
      progress: 0,
    },
  ];
  
  const sampleTasks = [
    {
      id: "task-1",
      name: "Initialize environment",
      status: "completed" as const,
      progress: 100,
      startTime: new Date(Date.now() - 1000 * 60 * 5),
      endTime: new Date(Date.now() - 1000 * 60 * 4),
    },
    {
      id: "task-2",
      name: "Generate code",
      status: "in-progress" as const,
      progress: 75,
      startTime: new Date(Date.now() - 1000 * 60 * 3),
    },
    {
      id: "task-3",
      name: "Run tests",
      status: "pending" as const,
      progress: 0,
    },
  ];
  
  const sampleState = [
    {
      id: "state-1",
      name: "Agent State",
      value: null,
      children: [
        {
          id: "state-1-1",
          name: "currentTask",
          value: "Generate code",
        },
        {
          id: "state-1-2",
          name: "progress",
          value: 75,
        },
        {
          id: "state-1-3",
          name: "memory",
          value: null,
          children: [
            {
              id: "state-1-3-1",
              name: "context",
              value: "Building a UI Kit",
            },
            {
              id: "state-1-3-2",
              name: "history",
              value: ["Task started", "Environment initialized"],
            },
          ],
        },
      ],
    },
  ];
  
  const models = [
    { id: "gemini-2.5-pro", name: "Gemini 2.5 Pro", description: "Best for coding tasks" },
    { id: "llama-4-scout", name: "Llama 4 Scout", description: "Best for standard operations" },
    { id: "llama-4-maverick", name: "Llama 4 Maverick", description: "Best for reasoning tasks" },
  ];

  return (
    <PlatformProvider forcePlatform={platform}>
      <AnimatedBackground variant="grid" color="emerald">
        <div className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900">
          <header className="mb-8">
            <h1 className="text-4xl font-bold mb-2">
              <TextGradient>Kled UI Kit Demo</TextGradient>
            </h1>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              A showcase of UI components for building AI Agent interfaces
            </p>
            
            <div className="flex space-x-4 mb-6">
              <Button 
                variant={platform === "web" ? "default" : "outline"} 
                onClick={() => setPlatform("web")}
              >
                Web
              </Button>
              <Button 
                variant={platform === "electron" ? "default" : "outline"} 
                onClick={() => setPlatform("electron")}
              >
                Electron
              </Button>
              <Button 
                variant={platform === "mobile" ? "default" : "outline"} 
                onClick={() => setPlatform("mobile")}
              >
                Mobile
              </Button>
            </div>
            
            <PlatformSpecific
              web={<div className="text-sm p-2 bg-blue-100 dark:bg-blue-900 rounded">Running in Web mode</div>}
              electron={<div className="text-sm p-2 bg-purple-100 dark:bg-purple-900 rounded">Running in Electron mode</div>}
              mobile={<div className="text-sm p-2 bg-green-100 dark:bg-green-900 rounded">Running in Mobile mode</div>}
            />
          </header>
          
          <main>
            <Tabs value={activeTab} onValueChange={setActiveTab} className="mb-8">
              <TabsList>
                <TabsTrigger value="radix">Radix UI</TabsTrigger>
                <TabsTrigger value="shadcn">Shadcn UI</TabsTrigger>
                <TabsTrigger value="aceternity">Aceternity UI</TabsTrigger>
                <TabsTrigger value="tailwind">Tailwind UI</TabsTrigger>
                <TabsTrigger value="vercel-ai">Vercel AI</TabsTrigger>
                <TabsTrigger value="platform">Platform</TabsTrigger>
              </TabsList>
              
              <TabsContent value="radix" className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow mt-4">
                <h2 className="text-2xl font-bold mb-4">Radix UI Components</h2>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h3 className="text-lg font-medium mb-2">Dialog</h3>
                    <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                      <DialogTrigger asChild>
                        <Button>Open Dialog</Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Radix Dialog Example</DialogTitle>
                          <DialogDescription>
                            This is a dialog component from Radix UI.
                          </DialogDescription>
                        </DialogHeader>
                        <div className="py-4">
                          <p>Dialog content goes here.</p>
                        </div>
                        <div className="flex justify-end">
                          <Button onClick={() => setIsDialogOpen(false)}>Close</Button>
                        </div>
                      </DialogContent>
                    </Dialog>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">Dropdown Menu</h3>
                    <DropdownMenu open={isDropdownOpen} onOpenChange={setIsDropdownOpen}>
                      <DropdownMenuTrigger asChild>
                        <Button>Open Dropdown</Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent>
                        <DropdownMenuItem onSelect={() => console.log("Item 1 selected")}>
                          Item 1
                        </DropdownMenuItem>
                        <DropdownMenuItem onSelect={() => console.log("Item 2 selected")}>
                          Item 2
                        </DropdownMenuItem>
                        <DropdownMenuItem onSelect={() => console.log("Item 3 selected")}>
                          Item 3
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </TabsContent>
              
              <TabsContent value="shadcn" className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow mt-4">
                <h2 className="text-2xl font-bold mb-4">Shadcn UI Components</h2>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h3 className="text-lg font-medium mb-2">Buttons</h3>
                    <div className="flex flex-wrap gap-2">
                      <Button variant="default">Default</Button>
                      <Button variant="destructive">Destructive</Button>
                      <Button variant="outline">Outline</Button>
                      <Button variant="secondary">Secondary</Button>
                      <Button variant="ghost">Ghost</Button>
                      <Button variant="link">Link</Button>
                    </div>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">Card</h3>
                    <Card className="p-4">
                      <h4 className="font-medium mb-2">Card Title</h4>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        This is a card component from Shadcn UI.
                      </p>
                    </Card>
                  </div>
                  
                  <div className="col-span-1 md:col-span-2">
                    <h3 className="text-lg font-medium mb-2">Accordion</h3>
                    <Accordion type="single" collapsible>
                      <AccordionItem value="item-1">
                        <AccordionTrigger>Section 1</AccordionTrigger>
                        <AccordionContent>
                          Content for section 1 goes here.
                        </AccordionContent>
                      </AccordionItem>
                      <AccordionItem value="item-2">
                        <AccordionTrigger>Section 2</AccordionTrigger>
                        <AccordionContent>
                          Content for section 2 goes here.
                        </AccordionContent>
                      </AccordionItem>
                      <AccordionItem value="item-3">
                        <AccordionTrigger>Section 3</AccordionTrigger>
                        <AccordionContent>
                          Content for section 3 goes here.
                        </AccordionContent>
                      </AccordionItem>
                    </Accordion>
                  </div>
                </div>
              </TabsContent>
              
              <TabsContent value="aceternity" className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow mt-4">
                <h2 className="text-2xl font-bold mb-4">Aceternity UI Components</h2>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h3 className="text-lg font-medium mb-2">Gradient Card</h3>
                    <GradientCard className="p-4">
                      <h4 className="font-medium mb-2">Gradient Card</h4>
                      <p className="text-sm">
                        This card has a beautiful gradient border effect.
                      </p>
                    </GradientCard>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">Spotlight Button</h3>
                    <SpotlightButton>Spotlight Effect</SpotlightButton>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">Text Gradient</h3>
                    <h4 className="text-xl">
                      <TextGradient>This text has a gradient effect</TextGradient>
                    </h4>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">3D Card</h3>
                    <ThreeDCard className="p-4 h-40">
                      <div className="h-full flex items-center justify-center">
                        <h4 className="font-medium">3D Tilt Effect Card</h4>
                      </div>
                    </ThreeDCard>
                  </div>
                </div>
              </TabsContent>
              
              <TabsContent value="tailwind" className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow mt-4">
                <h2 className="text-2xl font-bold mb-4">Tailwind UI Components</h2>
                
                <div className="grid grid-cols-1 gap-6">
                  <div>
                    <h3 className="text-lg font-medium mb-2">Agent Card</h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <AgentCard
                        title="Code Generation Agent"
                        description="Generates code based on natural language descriptions"
                        status="running"
                        progress={65}
                      />
                      <AgentCard
                        title="Data Analysis Agent"
                        description="Analyzes data and generates insights"
                        status="completed"
                        progress={100}
                      />
                      <AgentCard
                        title="Testing Agent"
                        description="Runs tests and reports results"
                        status="idle"
                        progress={0}
                      />
                    </div>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">Agent Dashboard</h3>
                    <AgentDashboard
                      agents={sampleAgents}
                      onAgentClick={(id) => console.log(`Agent ${id} clicked`)}
                    />
                  </div>
                </div>
              </TabsContent>
              
              <TabsContent value="vercel-ai" className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow mt-4">
                <h2 className="text-2xl font-bold mb-4">Vercel AI Components</h2>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="col-span-1 md:col-span-2">
                    <h3 className="text-lg font-medium mb-2">Model Selector</h3>
                    <ModelSelector
                      models={models}
                      selectedModelId="gemini-2.5-pro"
                      onModelSelect={(modelId) => console.log(`Selected model: ${modelId}`)}
                    />
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">AI Task Tracker</h3>
                    <AITaskTracker
                      tasks={sampleTasks}
                      currentTaskId="task-2"
                      onTaskClick={(taskId) => console.log(`Task ${taskId} clicked`)}
                    />
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-medium mb-2">AI State Visualizer</h3>
                    <AIStateVisualizer
                      state={sampleState}
                      onStateNodeClick={(nodeId) => console.log(`State node ${nodeId} clicked`)}
                    />
                  </div>
                </div>
              </TabsContent>
              
              <TabsContent value="platform" className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow mt-4">
                <h2 className="text-2xl font-bold mb-4">Platform-Specific Components</h2>
                
                <div className="grid grid-cols-1 gap-6">
                  <ElectronWrapper>
                    <Card className="p-4 border-2 border-purple-500">
                      <h3 className="text-lg font-medium mb-2">Electron-Specific Component</h3>
                      <p>This component only renders in Electron desktop apps.</p>
                    </Card>
                  </ElectronWrapper>
                  
                  <LynxWrapper>
                    <Card className="p-4 border-2 border-green-500">
                      <h3 className="text-lg font-medium mb-2">Lynx-React Mobile Component</h3>
                      <p>This component only renders in Lynx-React mobile apps.</p>
                    </Card>
                  </LynxWrapper>
                  
                  <PlatformSpecific
                    web={
                      <Card className="p-4 border-2 border-blue-500">
                        <h3 className="text-lg font-medium mb-2">Web-Specific Component</h3>
                        <p>This component only renders in web apps.</p>
                      </Card>
                    }
                    electron={
                      <Card className="p-4 border-2 border-purple-500">
                        <h3 className="text-lg font-medium mb-2">Electron-Specific Component</h3>
                        <p>This component only renders in Electron desktop apps.</p>
                      </Card>
                    }
                    mobile={
                      <Card className="p-4 border-2 border-green-500">
                        <h3 className="text-lg font-medium mb-2">Mobile-Specific Component</h3>
                        <p>This component only renders in Lynx-React mobile apps.</p>
                      </Card>
                    }
                  />
                </div>
              </TabsContent>
            </Tabs>
          </main>
        </div>
      </AnimatedBackground>
    </PlatformProvider>
  );
};

export default UIKitDemo;
