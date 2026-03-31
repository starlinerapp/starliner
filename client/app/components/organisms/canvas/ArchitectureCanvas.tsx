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
import getElkLayout from "~/service/reactflow/getElkLayout";
import { useMatch, useNavigate, useParams } from "react-router";
import GitNode from "~/components/atoms/nodes/GitNode";
import {
  buildEdgePairs,
  buildTargets,
  makeBaseNode,
} from "~/service/reactflow/helpers";

interface ArchitectureCanvasProps {
  environment: ResponseEnvironment;
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
      git: GitNode,
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

  const isOnDetailPage = !!useMatch(
    "/starliner/projects/:id/:environment/architecture/:type/:deploymentId",
  );
  function handlePlaneClick() {
    if (isOnDetailPage) {
      navigate(
        `/${slug}/projects/${organizationId}/${environment.slug}/architecture`,
        { relative: "path", replace: true },
      );
    }
  }

  const graphFingerprint = useMemo(() => {
    if (!deployments) return null;

    const nodeIds = [
      ...deployments.databases.map((db) => `db:${db.id}`).sort(),
      ...deployments.images.map((img) => `img:${img.id}`).sort(),
      ...deployments.ingresses.map((ing) => `ing:${ing.id}`).sort(),
      ...deployments.gitDeployments.map((git) => `git:${git.id}`).sort(),
    ].join(",");

    const targets = buildTargets(deployments);
    const edgeKey = buildEdgePairs(deployments, targets)
      .map((e) => `${e.source}->${e.target}`)
      .sort()
      .join(",");

    return `${nodeIds}|${edgeKey}`;
  }, [deployments]);

  const rawGraph = useMemo(() => {
    if (!deployments || graphFingerprint === null) return null;

    const rawNodes: Node[] = [
      ...deployments.databases.map((db) => makeBaseNode(db.id, "database", db)),
      ...deployments.images.map((img) => makeBaseNode(img.id, "image", img)),
      ...deployments.ingresses.map((ing) =>
        makeBaseNode(ing.id, "ingress", ing),
      ),
      ...deployments.gitDeployments.map((git) =>
        makeBaseNode(git.id, "git", git),
      ),
    ];

    const targets = buildTargets(deployments);
    const rawEdges: Edge[] = buildEdgePairs(deployments, targets).map((p) => ({
      id: `${p.source}->${p.target}:${p.kind}`,
      source: p.source,
      target: p.target,
    }));

    return { rawNodes, rawEdges };
  }, [graphFingerprint]);

  const didInitialFitViewRef = useRef(false);

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

      requestAnimationFrame(() => {
        const duration = didInitialFitViewRef.current ? 500 : 0;
        fitView({ maxZoom: 1, duration });
        didInitialFitViewRef.current = true;
      });
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
    for (const git of deployments.gitDeployments)
      map.set(String(git.id), { ...git });
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
          if (node.type) {
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
