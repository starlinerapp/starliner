import { AttachAddon } from "@xterm/addon-attach";
import { FitAddon } from "@xterm/addon-fit";
import type { Terminal } from "@xterm/xterm";
import { useCallback, useEffect, useMemo, useRef } from "react";
import { useXTerm } from "react-xtermjs";
import { debounce } from "throttle-debounce";

interface TerminalProps {
  webSocketUrl: string;
}

type TtyResizeMessage = {
  type: "resize";
  cols: number;
  rows: number;
};

const DEFAULT_TTY_ROWS = 24;
const DEFAULT_TTY_COLS = 80;

function sendTerminalResize(ws: WebSocket, cols: number, rows: number) {
  if (ws.readyState !== WebSocket.OPEN || cols <= 0 || rows <= 0) {
    return;
  }

  const message: TtyResizeMessage = { type: "resize", cols, rows };
  ws.send(JSON.stringify(message));
}

export default function TerminalClient({ webSocketUrl }: TerminalProps) {
  const wsRef = useRef<WebSocket | null>(null);
  const instanceRef = useRef<Terminal | null>(null);
  const attachAddonRef = useRef<AttachAddon | null>(null);
  const fitAddon = useMemo(() => new FitAddon(), []);
  const terminalAddons = useMemo(() => [fitAddon], [fitAddon]);

  const terminalOptions = useMemo(
    () => ({
      cursorStyle: "block" as const,
      scrollback: 1000,
      overviewRulerWidth: 0,
      theme: {
        background: "#ffffff",
        foreground: "#65636D",
        cursor: "#65636D",
        selectionBackground: "#bbd6fb",
      },
    }),
    [],
  );

  const syncTerminalSize = useCallback(
    (ws: WebSocket) => {
      const term = instanceRef.current;
      if (!term) return;

      fitAddon.fit();

      sendTerminalResize(ws, term.cols, term.rows);
      term.scrollToBottom();
    },
    [fitAddon],
  );

  const debouncedSyncTerminalSize = useMemo(
    () => debounce(50, (ws: WebSocket) => syncTerminalSize(ws)),
    [syncTerminalSize],
  );

  const { ref: terminalRef, instance } = useXTerm({
    addons: terminalAddons,
    options: terminalOptions,
    listeners: {
      onResize: () => {
        instanceRef.current?.scrollToBottom();
      },
    },
  });

  useEffect(() => {
    instanceRef.current = instance;
  }, [instance]);

  useEffect(() => {
    if (!instance) return;

    const ws = new WebSocket(
      `${webSocketUrl}?tty_height=${DEFAULT_TTY_ROWS}&tty_width=${DEFAULT_TTY_COLS}`,
    );
    ws.binaryType = "arraybuffer";
    wsRef.current = ws;

    ws.onopen = () => {
      attachAddonRef.current?.dispose();
      const attachAddon = new AttachAddon(ws);
      attachAddonRef.current = attachAddon;
      instance.loadAddon(attachAddon);

      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          syncTerminalSize(ws);
        });
      });
    };

    ws.onclose = () => {
      attachAddonRef.current?.dispose();
      attachAddonRef.current = null;
      wsRef.current = null;
    };

    ws.onerror = (err) => console.error("WebSocket error:", err);

    return () => {
      debouncedSyncTerminalSize.cancel();
      attachAddonRef.current?.dispose();
      attachAddonRef.current = null;
      ws.close();
    };
  }, [webSocketUrl, instance, syncTerminalSize, debouncedSyncTerminalSize]);

  useEffect(() => {
    const el = terminalRef.current;
    if (!el) return;

    const observer = new ResizeObserver(() => {
      const ws = wsRef.current;
      if (!ws) return;
      debouncedSyncTerminalSize(ws);
    });

    observer.observe(el);
    return () => observer.disconnect();
  }, [terminalRef, debouncedSyncTerminalSize]);

  return (
    <div className="flex min-h-0 flex-1 flex-col p-4">
      <div ref={terminalRef} className="min-h-0 flex-1" />
    </div>
  );
}
