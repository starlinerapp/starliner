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
import { useNavigate, useParams } from "react-router";

interface ArchitectureCanvasProps {
  environment: ResponseEnvironment;
}

export default function ArchitectureCanvas({
  environment,
}: ArchitectureCanvasProps) {
  const { slug, id: organizationId } = useParams<{
    slug: string;
    id: string;
  }>();

  const navigate = useNavigate();

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
      db: deployments.databases.map((d) => ({
        id: d.id,
        serviceName: d.serviceName,
      })),
      img: deployments.images.map((i) => ({
        id: i.id,
        serviceName: i.serviceName,
        port: i.port,
        envValues: i.envVars,
      })),
      ing: deployments.ingresses.map((ing) => ({
        id: ing.id,
        hosts: ing.hosts,
      })),
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

  function referencesImage(v: string, service: string, port: string): boolean {
    const candidates = [
      service,
      `${service}:${port}`,
      `http://${service}:${port}`,
    ];
    return candidates.some((prefix) => v === prefix || v.startsWith(prefix));
  }

  function referencesDatabase(v: string, db: { rwHost: string }): boolean {
    return v === db.rwHost;
  }

  function handleNodeSelected(type: string, id: string) {
    navigate(
      `${slug}/projects/${organizationId}/${environment.slug}/architecture/${type}/${id}`,
    );
  }

  function handlePlaneClick() {
    navigate(
      `${slug}/projects/${organizationId}/${environment.slug}/architecture`,
    );
  }

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

    const rawEdges: Edge[] = [];
    const pushEdge = (source: string, target: string, kind: string) => {
      rawEdges.push({
        id: `${source}->${target}:${kind}`,
        source,
        target,
        type: "smoothstep",
      });
    };

    const databaseTargets = deployments.databases.map((db) => ({
      id: String(db.id),
      serviceName: db.serviceName,
      rwHost: `${db.serviceName}-rw`,
    }));

    const imageTargets = deployments.images.map((t) => ({
      id: String(t.id),
      serviceName: t.serviceName,
      port: t.port,
    }));

    for (const ing of deployments.ingresses) {
      for (const host of ing.hosts ?? []) {
        for (const path of host.paths ?? []) {
          const matchedTarget = imageTargets.find((t) =>
            referencesImage(path.serviceName, t.serviceName, t.port),
          );
          if (!matchedTarget) continue;

          pushEdge(
            String(ing.id),
            matchedTarget.id,
            `ing:${host.host}:${path.path}`,
          );
        }
      }
    }

    for (const src of deployments.images) {
      const sourceId = String(src.id);

      for (const v of src.envVars) {
        for (const db of databaseTargets) {
          if (referencesDatabase(v.value, db)) {
            pushEdge(sourceId, db.id, `img->db:${db.serviceName}`);
          }
        }
        for (const t of imageTargets) {
          if (t.id === sourceId) continue;

          if (referencesImage(v.value, t.serviceName, t.port)) {
            pushEdge(sourceId, t.id, `img->img:${t.serviceName}:${t.port}`);
          }
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
        const updated = [
          ...deployments.images,
          ...deployments.databases,
          ...deployments.ingresses,
        ].find((i) => String(i.id) === node.id);

        if (!updated) return node;
        return {
          ...node,
          data: { ...updated },
        };
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
        multiSelectionKeyCode={null}
        selectionKeyCode={null}
        onNodeClick={(_, node) => {
          if (node.type === "image" || node.type === "ingress") {
            handleNodeSelected(node.type, node.id);
          }
        }}
        onPaneClick={() => handlePlaneClick()}
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
