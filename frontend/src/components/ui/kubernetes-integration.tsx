import React, { useState } from "react";
import { cn } from "../../utils/cn";
import { SpotlightCard } from "./aceternity/spotlight-card";
import { GradientButton } from "./aceternity/gradient-button";
import { AnimatedTabs } from "./aceternity/animated-tabs";

interface KubernetesIntegrationProps {
  className?: string;
}

export const KubernetesIntegration: React.FC<KubernetesIntegrationProps> = ({
  className,
}) => {
  const [isRefreshing, setIsRefreshing] = useState(false);
  
  const resources = {
    pods: [
      { name: "kled-agent-control-plane-7d8f9c6b5-x2v4r", status: "Running", namespace: "kled-system", age: "2h", ready: "1/1" },
      { name: "kled-agent-worker-5f7d8c9b6-zxc43", status: "Running", namespace: "kled-system", age: "2h", ready: "1/1" },
      { name: "kled-mcp-server-6c5d4b3a2-qwe12", status: "Running", namespace: "kled-system", age: "2h", ready: "1/1" },
    ],
    services: [
      { name: "kled-agent-api", type: "ClusterIP", clusterIP: "10.96.0.1", externalIP: "-", ports: "8000/TCP", age: "2h" },
      { name: "kled-mcp-server", type: "ClusterIP", clusterIP: "10.96.0.2", externalIP: "-", ports: "9000/TCP", age: "2h" },
    ],
    deployments: [
      { name: "kled-agent-control-plane", ready: "1/1", upToDate: "1", available: "1", age: "2h" },
      { name: "kled-agent-worker", ready: "1/1", upToDate: "1", available: "1", age: "2h" },
      { name: "kled-mcp-server", ready: "1/1", upToDate: "1", available: "1", age: "2h" },
    ],
  };
  
  const metrics = {
    cpu: {
      usage: "250m",
      limit: "1000m",
      percentage: 25,
    },
    memory: {
      usage: "256Mi",
      limit: "1Gi",
      percentage: 25,
    },
    storage: {
      usage: "2Gi",
      limit: "10Gi",
      percentage: 20,
    },
  };
  
  const handleRefresh = () => {
    setIsRefreshing(true);
    setTimeout(() => {
      setIsRefreshing(false);
    }, 1000);
  };
  
  const tabs = [
    {
      id: "pods",
      label: "Pods",
      content: (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="px-4 py-2 text-left">Name</th>
                <th className="px-4 py-2 text-left">Ready</th>
                <th className="px-4 py-2 text-left">Status</th>
                <th className="px-4 py-2 text-left">Namespace</th>
                <th className="px-4 py-2 text-left">Age</th>
              </tr>
            </thead>
            <tbody>
              {resources.pods.map((pod) => (
                <tr key={pod.name} className="border-b border-gray-200 dark:border-gray-700">
                  <td className="px-4 py-2 font-medium">{pod.name}</td>
                  <td className="px-4 py-2">{pod.ready}</td>
                  <td className="px-4 py-2">
                    <span className={cn(
                      "px-2 py-1 rounded-full text-xs",
                      pod.status === "Running" ? "bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300" :
                      pod.status === "Pending" ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300" :
                      "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300"
                    )}>
                      {pod.status}
                    </span>
                  </td>
                  <td className="px-4 py-2">{pod.namespace}</td>
                  <td className="px-4 py-2">{pod.age}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ),
    },
    {
      id: "services",
      label: "Services",
      content: (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="px-4 py-2 text-left">Name</th>
                <th className="px-4 py-2 text-left">Type</th>
                <th className="px-4 py-2 text-left">Cluster IP</th>
                <th className="px-4 py-2 text-left">External IP</th>
                <th className="px-4 py-2 text-left">Ports</th>
                <th className="px-4 py-2 text-left">Age</th>
              </tr>
            </thead>
            <tbody>
              {resources.services.map((service) => (
                <tr key={service.name} className="border-b border-gray-200 dark:border-gray-700">
                  <td className="px-4 py-2 font-medium">{service.name}</td>
                  <td className="px-4 py-2">{service.type}</td>
                  <td className="px-4 py-2">{service.clusterIP}</td>
                  <td className="px-4 py-2">{service.externalIP}</td>
                  <td className="px-4 py-2">{service.ports}</td>
                  <td className="px-4 py-2">{service.age}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ),
    },
    {
      id: "deployments",
      label: "Deployments",
      content: (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700">
                <th className="px-4 py-2 text-left">Name</th>
                <th className="px-4 py-2 text-left">Ready</th>
                <th className="px-4 py-2 text-left">Up-to-date</th>
                <th className="px-4 py-2 text-left">Available</th>
                <th className="px-4 py-2 text-left">Age</th>
              </tr>
            </thead>
            <tbody>
              {resources.deployments.map((deployment) => (
                <tr key={deployment.name} className="border-b border-gray-200 dark:border-gray-700">
                  <td className="px-4 py-2 font-medium">{deployment.name}</td>
                  <td className="px-4 py-2">{deployment.ready}</td>
                  <td className="px-4 py-2">{deployment.upToDate}</td>
                  <td className="px-4 py-2">{deployment.available}</td>
                  <td className="px-4 py-2">{deployment.age}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ),
    },
    {
      id: "metrics",
      label: "Metrics",
      content: (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 p-4">
          <div className="p-4 border rounded-md bg-white/5">
            <div className="flex justify-between items-center mb-2">
              <h4 className="font-medium">CPU Usage</h4>
              <span className="text-xs text-gray-500 dark:text-gray-400">{metrics.cpu.usage} / {metrics.cpu.limit}</span>
            </div>
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
              <div 
                className="bg-emerald-500 h-2.5 rounded-full" 
                style={{ width: `${metrics.cpu.percentage}%` }}
              ></div>
            </div>
          </div>
          
          <div className="p-4 border rounded-md bg-white/5">
            <div className="flex justify-between items-center mb-2">
              <h4 className="font-medium">Memory Usage</h4>
              <span className="text-xs text-gray-500 dark:text-gray-400">{metrics.memory.usage} / {metrics.memory.limit}</span>
            </div>
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
              <div 
                className="bg-emerald-500 h-2.5 rounded-full" 
                style={{ width: `${metrics.memory.percentage}%` }}
              ></div>
            </div>
          </div>
          
          <div className="p-4 border rounded-md bg-white/5">
            <div className="flex justify-between items-center mb-2">
              <h4 className="font-medium">Storage Usage</h4>
              <span className="text-xs text-gray-500 dark:text-gray-400">{metrics.storage.usage} / {metrics.storage.limit}</span>
            </div>
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
              <div 
                className="bg-emerald-500 h-2.5 rounded-full" 
                style={{ width: `${metrics.storage.percentage}%` }}
              ></div>
            </div>
          </div>
        </div>
      ),
    },
  ];

  return (
    <SpotlightCard className={cn("flex flex-col", className)}>
      <div className="flex justify-between items-center p-4 border-b">
        <h3 className="text-lg font-medium">Kubernetes Integration</h3>
        <GradientButton
          size="sm"
          variant="outline"
          onClick={handleRefresh}
          disabled={isRefreshing}
        >
          {isRefreshing ? "Refreshing..." : "Refresh"}
        </GradientButton>
      </div>
      
      <div className="flex-1 overflow-y-auto">
        <AnimatedTabs tabs={tabs} defaultTabId="pods" />
      </div>
    </SpotlightCard>
  );
};

export default KubernetesIntegration;
