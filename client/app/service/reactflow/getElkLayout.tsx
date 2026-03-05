// NOTE: Initial draft generated with AI assistance.
// Code has been reviewed and modified.

import ELK from "elkjs/lib/elk.bundled.js";
import type { ElkNode } from "elkjs/lib/elk.bundled";
import type { Edge, Node } from "@xyflow/react";

const DEFAULT_H = 120;
const DEFAULT_W = 220;
const GAP = 150;

function getNodeSize(n: Node) {
  const w = n.measured?.width ?? n.style?.width ?? n.width;
  const h = n.measured?.height ?? n.style?.height ?? n.height;
  return { width: Number(w) || DEFAULT_W, height: Number(h) };
}

type Component = {
  nodeIds: string[];
  edgeIds: string[];
};

function buildComponents(nodes: Node[], edges: Edge[]): Component[] {
  const nodeIds = new Set(nodes.map((n) => n.id));

  const adj = new Map<string, Set<string>>();
  for (const n of nodes) adj.set(n.id, new Set());

  const edgeIdsByNode = new Map<string, string[]>();
  for (const n of nodes) edgeIdsByNode.set(n.id, []);

  for (const e of edges) {
    if (!nodeIds.has(e.source) || !nodeIds.has(e.target)) continue;
    adj.get(e.source)!.add(e.target);
    adj.get(e.target)!.add(e.source);
    edgeIdsByNode.get(e.source)!.push(e.id);
    edgeIdsByNode.get(e.target)!.push(e.id);
  }

  const visited = new Set<string>();
  const comps: Component[] = [];

  for (const n of nodes) {
    if (visited.has(n.id)) continue;

    const queue = [n.id];
    visited.add(n.id);

    const compNodes: string[] = [];
    const compEdgeIds = new Set<string>();

    while (queue.length) {
      const cur = queue.pop()!;
      compNodes.push(cur);

      for (const eid of edgeIdsByNode.get(cur) ?? []) compEdgeIds.add(eid);

      for (const nb of adj.get(cur) ?? []) {
        if (!visited.has(nb)) {
          visited.add(nb);
          queue.push(nb);
        }
      }
    }
    comps.push({ nodeIds: compNodes, edgeIds: [...compEdgeIds] });
  }

  return comps;
}

function boundsOfLaidOutNodes(
  laidOut: {
    id: string;
    x?: number;
    y?: number;
    width?: number;
    height?: number;
  }[],
) {
  let minX = Infinity,
    minY = Infinity,
    maxX = -Infinity,
    maxY = -Infinity;
  for (const n of laidOut) {
    const x = n.x ?? 0;
    const y = n.y ?? 0;
    const w = n.width ?? DEFAULT_W;
    const h = n.height ?? DEFAULT_H;
    minX = Math.min(minX, x);
    minY = Math.min(minY, y);
    maxX = Math.max(maxX, x + w);
    maxY = Math.max(maxY, y + h);
  }
  if (!isFinite(minX)) return { minX: 0, minY: 0, width: 0, height: 0 };
  return { minX, minY, width: maxX - minX, height: maxY - minY };
}

function packRects(
  rects: { id: string; width: number; height: number }[],
  gap = GAP,
  maxColHeight = 2400,
) {
  let x = 0,
    y = 0,
    colW = 0;

  const placements = new Map<string, { x: number; y: number }>();

  for (const r of rects) {
    if (y > 0 && y + r.height > maxColHeight) {
      y = 0;
      x += colW + gap;
      colW = 0;
    }

    placements.set(r.id, { x, y });

    y += r.height + gap;
    colW = Math.max(colW, r.width);
  }

  return placements;
}

export default async function getElkLayout(
  nodes: Node[] = [],
  edges: Edge[] = [],
) {
  const targetPosition = "left" as const;
  const sourcePosition = "right" as const;

  const elk = new ELK();

  const comps = buildComponents(nodes, edges);

  const laidOutByComp = await Promise.all(
    comps.map(async (c, idx) => {
      const compNodes = nodes.filter((n) => c.nodeIds.includes(n.id));
      const compEdges = edges.filter((e) => c.edgeIds.includes(e.id));

      const hasEdges = compEdges.length > 0;

      const graph: ElkNode = {
        id: `comp-${idx}`,
        layoutOptions: hasEdges
          ? {
              "elk.algorithm": "org.eclipse.elk.layered",
              "elk.direction": "RIGHT",
              "elk.edgeRouting": "POLYLINE",
              "elk.spacing.nodeNode": "200",
              "elk.spacing.edgeNode": "400",
              "elk.layered.spacing.nodeNodeBetweenLayers": "200",
            }
          : {
              "elk.algorithm": "org.eclipse.elk.rectpacking",
              "elk.direction": "RIGHT",
              "elk.spacing.nodeNode": "200",
            },
        children: compNodes.map((n) => {
          const { width, height } = getNodeSize(n);
          return { id: n.id, width, height };
        }),
        edges: compEdges.map((e) => ({
          id: e.id,
          sources: [e.source],
          targets: [e.target],
        })),
      };

      const layout = await elk.layout(graph);
      const children = layout.children ?? [];
      const b = boundsOfLaidOutNodes(children);

      return {
        id: `comp-${idx}`,
        nodeIds: c.nodeIds,
        edgeIds: c.edgeIds,
        layout,
        bounds: b,
      };
    }),
  );

  const placements = packRects(
    laidOutByComp.map((c) => ({
      id: c.id,
      width: c.bounds.width + GAP,
      height: c.bounds.height + GAP,
    })),
  );

  const nodePos = new Map<string, { x: number; y: number }>();

  for (const comp of laidOutByComp) {
    const place = placements.get(comp.id)!;
    const children = comp.layout.children ?? [];

    for (const ln of children) {
      const x = (ln.x ?? 0) - comp.bounds.minX + place.x;
      const y = (ln.y ?? 0) - comp.bounds.minY + place.y;
      nodePos.set(ln.id, { x, y });
    }
  }

  return {
    nodes: nodes.map((n) => ({
      ...n,
      position: nodePos.get(n.id) ?? n.position ?? { x: 0, y: 0 },
      sourcePosition,
      targetPosition,
    })) as Node[],
    edges,
  };
}
