// NOTE: Initial draft generated with AI assistance.
// Code has been reviewed and modified.

import type { ResponseDeployments } from "../../server/api/clients/server/generated";
import { Position } from "@xyflow/react";

// ─── Types ────────────────────────────────────────────────────────────────────

interface DatabaseTarget {
  id: string;
  serviceName: string;
  rwHost: string;
}

interface ServiceTarget {
  id: string;
  serviceName: string;
  port: string;
}

interface Targets {
  databaseTargets: DatabaseTarget[];
  imageTargets: ServiceTarget[];
  gitTargets: ServiceTarget[];
}

interface EdgePair {
  source: string;
  target: string;
  kind: string;
}

type ServiceSource = { id: number | string; envVars: Array<{ value: string }> };

// ─── Reference Matching ───────────────────────────────────────────────────────

const serviceHostCandidates = (serviceName: string, port: string) => [
  serviceName,
  `${serviceName}:${port}`,
  `http://${serviceName}:${port}`,
];

const referencesService = (value: string, serviceName: string, port: string) =>
  serviceHostCandidates(serviceName, port).some(
    (candidate) => value === candidate || value.startsWith(candidate),
  );

const isPostgresUrlProtocol = (protocol: string) =>
  protocol === "postgres:" || protocol === "postgresql:";

/** Host from a postgres[ql] or jdbc:postgres[ql] URL, or undefined if not parseable as one. */
const hostFromPostgresConnectionString = (
  value: string,
): string | undefined => {
  if (!/^[a-z+.-]+:\/\//i.test(value) && !/^jdbc:/i.test(value)) {
    return undefined;
  }
  const withoutJdbc = value.replace(/^jdbc:/i, "");
  if (!/^[a-z+.-]+:\/\//i.test(withoutJdbc)) {
    return undefined;
  }
  try {
    const url = new URL(withoutJdbc);
    if (!isPostgresUrlProtocol(url.protocol) || !url.hostname) {
      return undefined;
    }
    return url.hostname;
  } catch {
    return undefined;
  }
};

const referencesDatabase = (value: string, { rwHost }: DatabaseTarget) => {
  if (value === rwHost) return true;
  return hostFromPostgresConnectionString(value) === rwHost;
};

// ─── Target Builders ──────────────────────────────────────────────────────────

export const buildTargets = (deployments: ResponseDeployments): Targets => ({
  databaseTargets: deployments.databases.map((db) => ({
    id: String(db.id),
    serviceName: db.serviceName,
    rwHost: `${db.serviceName}-rw`,
  })),
  imageTargets: deployments.images.map((t) => ({
    id: String(t.id),
    serviceName: t.serviceName,
    port: t.port,
  })),
  gitTargets: deployments.gitDeployments.map((t) => ({
    id: String(t.id),
    serviceName: t.serviceName,
    port: t.port,
  })),
});

// ─── Edge Builders ────────────────────────────────────────────────────────────

const findServiceMatch = (
  serviceName: string,
  pools: ServiceTarget[][],
): ServiceTarget | undefined => {
  for (const pool of pools) {
    const match = pool.find((t) =>
      referencesService(serviceName, t.serviceName, t.port),
    );
    if (match) return match;
  }
};

const buildIngressEdges = (
  deployments: ResponseDeployments,
  { imageTargets, gitTargets }: Targets,
): EdgePair[] => {
  const servicePools = [imageTargets, gitTargets];
  return deployments.ingresses.flatMap((ing) =>
    (ing.hosts ?? []).flatMap((host) =>
      (host.paths ?? []).flatMap((path) => {
        const match = findServiceMatch(path.serviceName, servicePools);
        if (!match) return [];
        return [
          {
            source: String(ing.id),
            target: match.id,
            kind: `ing:${host.host}:${path.path}`,
          },
        ];
      }),
    ),
  );
};

const buildServiceEdges = (
  sources: ServiceSource[],
  sourcePrefix: string,
  selfTargets: ServiceTarget[],
  databaseTargets: DatabaseTarget[],
): EdgePair[] =>
  sources.flatMap((src) => {
    const sourceId = String(src.id);
    return src.envVars.flatMap(({ value }) => {
      const dbEdges = databaseTargets
        .filter((db) => referencesDatabase(value, db))
        .map((db) => ({
          source: sourceId,
          target: db.id,
          kind: `${sourcePrefix}->db:${db.serviceName}`,
        }));

      const serviceEdges = selfTargets
        .filter(
          (t) =>
            t.id !== sourceId &&
            referencesService(value, t.serviceName, t.port),
        )
        .map((t) => ({
          source: sourceId,
          target: t.id,
          kind: `${sourcePrefix}->${sourcePrefix}:${t.serviceName}:${t.port}`,
        }));

      return [...dbEdges, ...serviceEdges];
    });
  });

export const buildEdgePairs = (
  deployments: ResponseDeployments,
  targets: Targets,
): EdgePair[] => {
  const { databaseTargets, imageTargets, gitTargets } = targets;
  return [
    ...buildIngressEdges(deployments, targets),
    ...buildServiceEdges(
      deployments.images,
      "img",
      imageTargets,
      databaseTargets,
    ),
    ...buildServiceEdges(
      deployments.gitDeployments,
      "git",
      gitTargets,
      databaseTargets,
    ),
  ];
};

// ─── Node Factory ─────────────────────────────────────────────────────────────

export const makeBaseNode = <T extends object>(
  id: string | number,
  type: string,
  data: T,
) => ({
  id: String(id),
  type,
  position: { x: 0, y: 0 },
  sourcePosition: Position.Right,
  targetPosition: Position.Left,
  data: data as Record<string, unknown>,
});
