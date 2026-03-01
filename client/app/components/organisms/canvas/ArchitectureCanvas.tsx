import React, { useCallback, useEffect, useMemo, useState } from "react";
import {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Background,
  BackgroundVariant,
  Controls,
  type Edge,
  type Node,
  type NodeTypes,
  type OnConnect,
  type OnEdgesChange,
  type OnNodesChange,
  Position,
  ReactFlow,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import type { ResponseEnvironment } from "~/server/api/client/generated";
import DatabaseNode from "~/components/atoms/nodes/DatabaseNode";
import ImageNode from "~/components/atoms/nodes/ImageNode";
import IngressNode from "~/components/atoms/nodes/IngressNode";
import getElkLayout from "~/utils/reactflow/getElkLayout";

interface ArchitectureCanvasProps {
  environment: ResponseEnvironment;
}

export default function ArchitectureCanvas({
  environment,
}: ArchitectureCanvasProps) {
  const trpc = useTRPC();
  const { data: deployments } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      {
        id: environment?.id,
      },
      {
        refetchInterval: 2000, // 2 seconds
      },
    ),
  );

  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  const nodeTypes: NodeTypes = useMemo(() => {
    return {
      database: DatabaseNode,
      image: ImageNode,
      ingress: IngressNode,
    };
  }, []);

  const onNodesChange: OnNodesChange = useCallback(
    (changes) =>
      setNodes((nodesSnapshot) => applyNodeChanges(changes, nodesSnapshot)),
    [],
  );
  const onEdgesChange: OnEdgesChange = useCallback(
    (changes) =>
      setEdges((edgesSnapshot) => applyEdgeChanges(changes, edgesSnapshot)),
    [],
  );
  const onConnect: OnConnect = useCallback(
    (params) => setEdges((edgesSnapshot) => addEdge(params, edgesSnapshot)),
    [],
  );

  useEffect(() => {
    if (!deployments) return;

    const rawNodes: Node[] = [
      ...deployments.databases.map((d) => ({
        id: `database:${d.id}`,
        type: "database",
        position: { x: 0, y: 0 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
        data: {
          id: d.id,
          serviceName: d.name,
          status: d.status,
          port: d.port,
          username: d.username,
          password: d.password,
        },
      })),
      ...deployments.images.map((i) => ({
        id: `image:${i.id}`,
        type: "image",
        position: { x: 0, y: 0 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
        data: {
          id: i.id,
          serviceName: i.serviceName,
          status: i.status,
          port: i.port,
          imageName: i.imageName,
          tag: i.tag,
        },
      })),
      ...deployments.ingresses.map((ing) => ({
        id: `ingress:${ing.id}`,
        type: "ingress",
        position: { x: 0, y: 0 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
        data: {
          id: ing.id,
          serviceName: ing.serviceName,
          status: ing.status,
          port: ing.port,
        },
      })),
    ];

    const rawEdges: Edge[] = [];
    // rawEdges.push({
    //   id: "e1",
    //   source: "ingress:76",
    //   target: "image:75",
    //   type: "smoothstep",
    // });

    let cancelled = false;

    (async () => {
      const laidOut = await getElkLayout(rawNodes, rawEdges);
      if (cancelled) return;
      setNodes(laidOut.nodes);
      setEdges(laidOut.edges);
    })();

    return () => {
      cancelled = true;
    };
  }, [deployments]);

  return (
    <div className="h-full w-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeOrigin={[0, 0.5]}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        proOptions={{ hideAttribution: true }}
        fitView
        maxZoom={1}
        fitViewOptions={{ maxZoom: 1 }}
      >
        <Background gap={20} color="#84828E" variant={BackgroundVariant.Dots} />
        <Controls />
      </ReactFlow>
    </div>
  );
}
