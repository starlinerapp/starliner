import { motion } from "framer-motion";
import { useEffect, useRef, useState } from "react";
import { cn } from "~/utils/cn";

interface NavigationBarProps<T extends string> {
  items: readonly T[];
  onSelect: (item: T) => void;
  selected: T;
}

export default function NavigationBar<T extends string>({
  items,
  onSelect,
  selected,
}: NavigationBarProps<T>) {
  const [activeRect, setActiveRect] = useState({ left: 0, width: 0 });
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    const activeLink = Array.from(
      containerRef.current.querySelectorAll<HTMLSpanElement>("span"),
    ).find((span) => span.textContent === selected);

    if (activeLink) {
      const rect = activeLink.getBoundingClientRect();
      const containerRect = containerRef.current.getBoundingClientRect();
      setActiveRect({
        left: rect.left - containerRect.left,
        width: rect.width,
      });
    }
  }, [selected]);

  return (
    <div className="bg-violet-1">
      <div
        ref={containerRef}
        className="relative flex w-full gap-4 border-mauve-6 border-b px-2 pt-2 pb-1 text-mauve-11 text-sm"
      >
        {items.map((item, i) => (
          <button
            type="button"
            key={i}
            className="relative z-10 cursor-pointer px-2 py-1.5"
            onClick={() => onSelect(item)}
          >
            <span
              className={cn(
                "pb-2",
                selected === item && "font-semibold text-violet-11",
              )}
            >
              {item}
            </span>
          </button>
        ))}

        <motion.div
          className="absolute bottom-0 h-[3px] rounded-md bg-violet-11"
          animate={{ left: activeRect.left, width: activeRect.width }}
          transition={{ type: "spring", stiffness: 300, damping: 30 }}
        />
      </div>
    </div>
  );
}
