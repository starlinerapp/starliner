import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import {
  addEdge,
  applyNodeChanges,
  Background,
  BackgroundVariant,
  Controls,
  type Edge,
  type Node,
  type NodeTypes,
  type OnConnect,
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

export default function ArchitectureCanvas({
  environment,
}: ArchitectureCanvasProps) {
  const { fitView } = useReactFlow();

  const {
    slug,
    id: organizationId,
    deploymentId,
  } = useParams<{
    slug: string;
    id: string;
    deploymentId: string;
  }>();

  const navigate = useNavigate();

  const trpc = useTRPC();
  const { data: deployments } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: environment?.id },
      { refetchInterval: 2000 },
    ),
  );

  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  const selectedIdsRef = useRef<Set<string>>(new Set());
  useEffect(() => {
    selectedIdsRef.current = new Set(
      nodes.filter((n) => n.selected).map((n) => n.id),
    );
  }, [nodes]);

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

  const onConnect: OnConnect = useCallback(
    (params) => setEdges((edgesSnapshot) => addEdge(params, edgesSnapshot)),
    [],
  );

  const onSelectionChange = ({ nodes }: { nodes: Node[] }) => {
    selectedIdsRef.current = new Set(nodes.map((n) => n.id));
  };

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

  const graphFingerprint = useMemo(() => {
    if (!deployments) return null;

    const nodeIds = [
      ...deployments.databases.map((db) => `db:${db.id}`).sort(),
      ...deployments.images.map((img) => `img:${img.id}`).sort(),
      ...deployments.ingresses.map((ing) => `ing:${ing.id}`).sort(),
    ].join(",");

    const edgeParts: string[] = [];

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
          if (matchedTarget) {
            edgeParts.push(`${ing.id}->${matchedTarget.id}`);
          }
        }
      }
    }

    for (const src of deployments.images) {
      for (const v of src.envVars) {
        for (const db of databaseTargets) {
          if (referencesDatabase(v.value, db)) {
            edgeParts.push(`${src.id}->${db.id}`);
          }
        }
        for (const t of imageTargets) {
          if (t.id === String(src.id)) continue;
          if (referencesImage(v.value, t.serviceName, t.port)) {
            edgeParts.push(`${src.id}->${t.id}`);
          }
        }
      }
    }

    return `${nodeIds}|${edgeParts.sort().join(",")}`;
  }, [deployments]);

  const rawGraph = useMemo(() => {
    if (!deployments || graphFingerprint === null) return null;

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

    return { rawNodes, rawEdges };
  }, [graphFingerprint]);

  useEffect(() => {
    if (!rawGraph) return;

    let cancelled = false;
    (async () => {
      const laidOut = await getElkLayout(rawGraph.rawNodes, rawGraph.rawEdges);
      if (cancelled) return;

      const selectedIds = selectedIdsRef.current;
      setNodes(
        laidOut.nodes.map((n) => ({
          ...n,
          selected: deploymentId
            ? n.id === deploymentId
            : selectedIds.has(n.id),
        })),
      );
      setEdges(laidOut.edges);
      requestAnimationFrame(() => fitView({ maxZoom: 1, duration: 500 }));
    })();

    return () => {
      cancelled = true;
    };
  }, [rawGraph]);

  // Sync selection whenever the deploymentId URL param changes
  useEffect(() => {
    setNodes((prev) =>
      prev.map((n) => ({
        ...n,
        selected: deploymentId ? n.id === deploymentId : false,
      })),
    );
  }, [deploymentId]);

  const nodeDataMap = useMemo(() => {
    if (!deployments) return null;
    const map = new Map<string, Record<string, unknown>>();
    for (const db of deployments.databases) map.set(String(db.id), { ...db });
    for (const img of deployments.images) map.set(String(img.id), { ...img });
    for (const ing of deployments.ingresses)
      map.set(String(ing.id), { ...ing });
    return map;
  }, [deployments]);

  useEffect(() => {
    if (!nodeDataMap) return;
    setNodes((prev) =>
      prev.map((n) => {
        const fresh = nodeDataMap.get(n.id);
        return fresh ? { ...n, data: fresh } : n;
      }),
    );
  }, [nodeDataMap]);

  return (
    <div className="h-full w-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeOrigin={[0, 0.5]}
        nodeTypes={nodeTypes}
        nodesDraggable={false}
        onNodesChange={onNodesChange}
        onSelectionChange={onSelectionChange}
        onConnect={onConnect}
        onNodeClick={(_, node) => {
          if (node.type === "image" || node.type === "ingress") {
            handleNodeSelected(node.type, node.id);
          }
        }}
        onPaneClick={() => handlePlaneClick()}
        proOptions={{ hideAttribution: true }}
        maxZoom={1}
        fitViewOptions={{ maxZoom: 1 }}
      >
        <Background gap={20} color="#84828E" variant={BackgroundVariant.Dots} />
        <Controls />
      </ReactFlow>
    </div>
  );
}
