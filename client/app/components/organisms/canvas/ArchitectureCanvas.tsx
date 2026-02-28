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

    setNodes((prevNodes) => {
      const nextIds = new Set<string>([
        ...deployments.databases.map((d) => `database:${d.id}`),
        ...deployments.images.map((i) => `image:${i.id}`),
      ]);

      const remainingNodes = prevNodes.filter((n) => nextIds.has(n.id));

      const dbNodes: Node[] = deployments.databases.map((deployment) => {
        const id = `database:${deployment.id}`;
        const existing = remainingNodes.find((n) => n.id === id);

        return {
          id,
          type: "database",
          position: existing?.position ?? { x: 0, y: 0 },
          sourcePosition: Position.Right,
          targetPosition: Position.Left,
          data: {
            id: deployment.id,
            serviceName: deployment.name,
            status: deployment.status,
            port: deployment.port,
            username: deployment.username,
            password: deployment.password,
          },
        };
      });

      const imageNodes: Node[] = deployments.images.map((deployment) => {
        const id = `image:${deployment.id}`;
        const existing = remainingNodes.find((n) => n.id === id);

        return {
          id,
          type: "image",
          position: existing?.position ?? { x: 0, y: 0 },
          sourcePosition: Position.Right,
          targetPosition: Position.Left,
          data: {
            id: deployment.id,
            serviceName: deployment.serviceName,
            status: deployment.status,
            port: deployment.port,
            imageName: deployment.imageName,
            tag: deployment.tag,
          },
        };
      });

      return [...dbNodes, ...imageNodes];
    });
  }, [deployments]);

  return (
    <div className="h-full w-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        proOptions={{ hideAttribution: true }}
        fitView
        maxZoom={1}
        fitViewOptions={{
          maxZoom: 1,
        }}
      >
        <Background gap={20} color="#84828E" variant={BackgroundVariant.Dots} />
        <Controls />
      </ReactFlow>
    </div>
  );
}
