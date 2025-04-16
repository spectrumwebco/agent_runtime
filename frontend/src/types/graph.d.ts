interface GraphNode {
  id: string;
  type: string;
  label: string;
  status: string;
}

interface GraphEdge {
  source: string;
  target: string;
  label: string;
}
