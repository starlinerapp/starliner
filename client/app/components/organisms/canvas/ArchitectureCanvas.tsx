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
  useReactFlow,
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
  const { fitView } = useReactFlow();

  const trpc = useTRPC();
  const { data: deployments } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: environment?.id },
      { refetchInterval: 2000 },
    ),
  );

  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  const topologyKey = useMemo(() => {
    if (!deployments) return "";

    return JSON.stringify({
      db: deployments.databases.map((d) => d.id),
      img: deployments.images.map((i) => i.id),
      ing: deployments.ingresses.map((i) => i.id),
    });
  }, [deployments]);

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
    fitView();
  }, [nodes.length]);

  useEffect(() => {
    if (!deployments) return;

    const baseNode = {
      position: { x: 0, y: 0 },
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
    };

    const rawNodes: Node[] = [
      ...deployments.databases.map((db) => ({
        id: String(db.id),
        type: "database",
        ...baseNode,
        data: { ...db },
      })),
      ...deployments.images.map((img) => ({
        id: String(img.id),
        type: "image",
        ...baseNode,
        data: { ...img },
      })),
      ...deployments.ingresses.map((ing) => ({
        id: String(ing.id),
        type: "ingress",
        ...baseNode,
        data: { ...ing },
      })),
    ];

    const imageIdByServiceName = new Map<string, string>(
      deployments.images.map((i) => [i.serviceName, String(i.id)]),
    );

    const rawEdges: Edge[] = [];

    for (const ing of deployments.ingresses) {
      const ingressNodeId = String(ing.id);

      for (const host of ing.hosts ?? []) {
        for (const path of host.paths ?? []) {
          const targetImageId = imageIdByServiceName.get(path.serviceName);

          if (!targetImageId) continue;

          rawEdges.push({
            id: `${ingressNodeId}->${targetImageId}:${host.host}:${path.path}`,
            source: ingressNodeId,
            target: targetImageId,
            type: "smoothstep",
          });
        }
      }
    }

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
  }, [topologyKey]);

  useEffect(() => {
    if (!deployments) return;

    setNodes((prev) =>
      prev.map((node) => {
        if (node.type === "image") {
          const img = deployments.images.find((i) => String(i.id) === node.id);
          if (!img) return node;

          return {
            ...node,
            data: { ...img },
          };
        }
        return node;
      }),
    );
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
