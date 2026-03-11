import React, { useCallback, useEffect, useRef, useState } from "react";
import { XTerm } from "react-xtermjs";
import { AttachAddon } from "@xterm/addon-attach";
import { FitAddon } from "@xterm/addon-fit";
import type { ITerminalAddon } from "@xterm/xterm";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/client/generated";

type Deployment =
  | ResponseGitDeployment
  | ResponseImageDeployment
  | ResponseIngressDeployment
  | ResponseDatabaseDeployment;

interface TerminalProps {
  deployment: Deployment | undefined;
}

export default function TerminalClient({ deployment }: TerminalProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const [addons, setAddons] = useState<ITerminalAddon[]>([]);

  const getTtySize = useCallback(() => {
    const el = containerRef.current;
    if (!el) return { rows: 24, cols: 80 };
    const rect = el.getBoundingClientRect();
    return {
      rows: Math.max(1, Math.floor(rect.height / 16)),
      cols: Math.max(1, Math.floor(rect.width / 8)),
    };
  }, []);

  useEffect(() => {
    const { rows, cols } = getTtySize();

    const ws = new WebSocket(
      `wss://${window.location.host}/ws/${deployment?.id}?tty_height=${rows}&tty_width=${cols}`,
    );
    ws.binaryType = "arraybuffer";

    ws.onopen = (event) => {
      const socket = event.target as WebSocket;
      const fitAddon = new FitAddon();
      const attachAddon = new AttachAddon(socket);
      fitAddonRef.current = fitAddon;
      setAddons([fitAddon, attachAddon]);
      // Defer fit until XTerm has had a chance to mount/render
      setTimeout(() => fitAddon.fit(), 50);
    };

    ws.onclose = () => {
      fitAddonRef.current = null;
      setAddons([]);
    };

    ws.onerror = (err) => console.error("WebSocket error:", err);

    return () => ws.close();
  }, [deployment, getTtySize]);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const observer = new ResizeObserver(() => {
      fitAddonRef.current?.fit();
    });

    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  return (
    <div ref={containerRef} className="h-full w-full p-4">
      <XTerm
        className="h-full w-full"
        addons={addons}
        options={{
          cursorStyle: "block",
          theme: {
            background: "#ffffff",
            foreground: "#65636D",
            cursor: "#65636D",
          },
        }}
      />
    </div>
  );
}
