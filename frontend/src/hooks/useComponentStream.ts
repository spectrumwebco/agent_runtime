import { useState, useEffect } from 'react';

interface Component {
  id: string;
  type: string;
  props: Record<string, any>;
  agent_id?: string;
  action_id?: string;
  tool_id?: string;
  created_at: number;
  updated_at: number;
}

export function useComponentStream(agentId?: string) {
  const [components, setComponents] = useState<Component[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let eventSource: EventSource | null = null;
    
    const connectToStream = () => {
      setIsLoading(true);
      
      if (eventSource) {
        eventSource.close();
      }
      
      const url = agentId 
        ? `/api/rsc/stream?agent_id=${agentId}` 
        : '/api/rsc/stream';
        
      eventSource = new EventSource(url);
      
      eventSource.onopen = () => {
        setIsLoading(false);
        setError(null);
      };
      
      eventSource.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          if (data.type === 'ping') {
            return;
          }
          
          setComponents(prevComponents => {
            const existingIndex = prevComponents.findIndex(c => c.id === data.id);
            
            if (existingIndex >= 0) {
              const updatedComponents = [...prevComponents];
              updatedComponents[existingIndex] = {
                ...updatedComponents[existingIndex],
                ...data,
                updated_at: Date.now()
              };
              return updatedComponents;
            } else {
              return [...prevComponents, {
                ...data,
                created_at: data.created_at || Date.now(),
                updated_at: data.updated_at || Date.now()
              }];
            }
          });
        } catch (err) {
          console.error('Error parsing event data:', err);
        }
      };
      
      eventSource.onerror = (err) => {
        console.error('EventSource error:', err);
        setError(new Error('Failed to connect to component stream'));
        setIsLoading(false);
        
        setTimeout(() => {
          connectToStream();
        }, 5000);
      };
    };
    
    connectToStream();
    
    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  }, [agentId]);
  
  useEffect(() => {
    const fetchInitialComponents = async () => {
      try {
        const url = agentId 
          ? `/api/rsc/components/agent/${agentId}` 
          : '/api/rsc/components';
          
        const response = await fetch(url);
        
        if (!response.ok) {
          throw new Error(`Failed to fetch components: ${response.statusText}`);
        }
        
        const data = await response.json();
        setComponents(data.components || []);
      } catch (err) {
        console.error('Error fetching initial components:', err);
        setError(err instanceof Error ? err : new Error(String(err)));
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchInitialComponents();
  }, [agentId]);
  
  return { components, isLoading, error };
}
