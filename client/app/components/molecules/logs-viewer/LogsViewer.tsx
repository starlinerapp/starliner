import type { ReactNode } from "react";
import { useEffect, useMemo, useRef, useState } from "react";
import {
  ChevronDown,
  ChevronUp,
  Cross,
  MagnifyingGlass,
} from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

interface LogsViewerProps {
  logs: string[];
  resetKey?: string | number;
}

const BOTTOM_THRESHOLD_PX = 32;
const TOP_THRESHOLD_PX = 32;

export default function LogsViewer({ logs, resetKey }: LogsViewerProps) {
  const [search, setSearch] = useState("");
  const [isAtBottom, setIsAtBottom] = useState(true);
  const [isAtTop, setIsAtTop] = useState(true);

  const logsScrollRef = useRef<HTMLDivElement>(null);
  const hasLoadedInitial = useRef(false);
  const autoFollowRef = useRef(true);

  const lastScrollTopRef = useRef(0);

  const filteredLogs = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) {
      return logs;
    }
    return logs.filter((line) => line.toLowerCase().includes(query));
  }, [logs, search]);

  const updateButtonState = (el: HTMLDivElement) => {
    const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    setIsAtBottom(distanceFromBottom <= BOTTOM_THRESHOLD_PX);
    setIsAtTop(el.scrollTop <= TOP_THRESHOLD_PX);
  };

  useEffect(() => {
    hasLoadedInitial.current = false;
    autoFollowRef.current = true;
    lastScrollTopRef.current = 0;
    setIsAtBottom(true);
    setIsAtTop(true);
    setSearch("");
  }, [resetKey]);

  useEffect(() => {
    const el = logsScrollRef.current;
    if (!el) {
      return;
    }

    if (!hasLoadedInitial.current) {
      if (filteredLogs.length > 0) {
        hasLoadedInitial.current = true;
        requestAnimationFrame(() => {
          requestAnimationFrame(() => {
            el.scrollTop = el.scrollHeight;
            updateButtonState(el);
          });
        });
      }
      return;
    }

    if (autoFollowRef.current) {
      el.scrollTop = el.scrollHeight;
    }
    updateButtonState(el);
  }, [filteredLogs]);

  const handleScroll = () => {
    const el = logsScrollRef.current;
    if (!el) {
      return;
    }
    updateButtonState(el);

    const prev = lastScrollTopRef.current;
    const curr = el.scrollTop;
    lastScrollTopRef.current = curr;

    const distanceFromBottom = el.scrollHeight - curr - el.clientHeight;
    if (curr < prev - 1) {
      autoFollowRef.current = false;
    } else if (distanceFromBottom <= BOTTOM_THRESHOLD_PX) {
      autoFollowRef.current = true;
    }
  };

  const scrollToTop = () => {
    const el = logsScrollRef.current;
    if (!el) {
      return;
    }
    autoFollowRef.current = false;
    el.scrollTo({ top: 0, left: 0, behavior: "smooth" });
  };

  const scrollToBottom = () => {
    const el = logsScrollRef.current;
    if (!el) {
      return;
    }
    autoFollowRef.current = true;
    el.scrollTo({
      top: el.scrollHeight - el.clientHeight,
      left: 0,
      behavior: "smooth",
    });
  };

  const highlightLine = (line: string) => {
    const query = search.trim();
    if (!query) {
      return line;
    }
    const lowerLine = line.toLowerCase();
    const lowerQuery = query.toLowerCase();
    const parts: ReactNode[] = [];
    let lastIndex = 0;
    let matchIndex = lowerLine.indexOf(lowerQuery);
    while (matchIndex !== -1) {
      if (matchIndex > lastIndex) {
        parts.push(line.slice(lastIndex, matchIndex));
      }
      parts.push(
        <mark
          key={matchIndex}
          className="rounded-xs bg-violet-4 text-violet-12"
        >
          {line.slice(matchIndex, matchIndex + query.length)}
        </mark>,
      );
      lastIndex = matchIndex + query.length;
      matchIndex = lowerLine.indexOf(lowerQuery, lastIndex);
    }
    if (lastIndex < line.length) {
      parts.push(line.slice(lastIndex));
    }
    return parts;
  };

  return (
    <div className="flex h-full min-h-0 w-full flex-col gap-2">
      <div className="relative shrink-0">
        <MagnifyingGlass className="absolute top-1/2 left-2 h-4 w-4 -translate-y-1/2 text-mauve-11" />
        <input
          className="w-full rounded-md border border-mauve-6 p-2 pr-7 pl-7 text-xs shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)] placeholder:text-mauve-11"
          type="text"
          placeholder="Search logs"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        {search && (
          <button
            type="button"
            onClick={() => setSearch("")}
            aria-label="Clear search"
            className="absolute top-1/2 right-2 -translate-y-1/2 text-mauve-11 hover:text-mauve-12"
          >
            <Cross className="h-3.5 w-3.5" />
          </button>
        )}
      </div>
      <div className="relative min-h-0 flex-1">
        <div
          ref={logsScrollRef}
          onScroll={handleScroll}
          className="h-full min-h-0 w-full overflow-y-auto [overflow-anchor:none]"
        >
          <pre className="w-full whitespace-pre-wrap break-all font-mono text-mauve-11 text-sm">
            {filteredLogs.map((line, i) => (
              <span key={i} className="block">
                {highlightLine(line)}
              </span>
            ))}
          </pre>
        </div>
        <div className="absolute right-3 bottom-3 mr-4 flex flex-col gap-1.5">
          <button
            type="button"
            onClick={scrollToTop}
            disabled={isAtTop}
            aria-label="Scroll to top"
            className={cn(
              "flex h-7 w-7 items-center justify-center rounded-md border border-mauve-6 bg-white text-mauve-12 shadow-md transition-opacity cursor-pointer",
              isAtTop
                ? "cursor-default opacity-40"
                : "opacity-100 hover:bg-mauve-1",
            )}
          >
            <ChevronUp className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={scrollToBottom}
            disabled={isAtBottom}
            aria-label="Scroll to bottom"
            className={cn(
              "flex h-7 w-7 items-center justify-center rounded-md border border-mauve-6 bg-white text-mauve-12 shadow-md transition-opacity cursor-pointer",
              isAtBottom
                ? "cursor-default opacity-40"
                : "opacity-100 hover:bg-mauve-1",
            )}
          >
            <ChevronDown className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
