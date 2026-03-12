import React, { useState, useLayoutEffect, useRef, useCallback } from "react";
import { NavLink, useLocation } from "react-router";
import { motion } from "framer-motion";
import { cn } from "~/utils/cn";

type NavigationBarItem = {
  title: React.ReactNode;
  href: string;
};

interface NavigationBarProps {
  items: NavigationBarItem[];
}

export default function LinkNavigationBar({ items }: NavigationBarProps) {
  const location = useLocation();

  const [activeRect, setActiveRect] = useState({ left: 0, width: 0 });
  const containerRef = useRef<HTMLDivElement>(null);

  const updateActiveRect = useCallback(() => {
    if (!containerRef.current) return;

    const links = containerRef.current.querySelectorAll<HTMLAnchorElement>("a");
    let activeLink: HTMLSpanElement | null = null;

    for (const link of links) {
      const href = link.getAttribute("href");
      if (href && location.pathname.startsWith(href)) {
        const span = link.querySelector<HTMLSpanElement>("span");

        if (
          span &&
          (!activeLink ||
            href.length >
              (activeLink.closest("a")?.getAttribute("href")?.length || 0))
        ) {
          activeLink = span;
        }
      }
    }

    if (!activeLink) return;

    const rect = activeLink.getBoundingClientRect();
    const containerRect = containerRef.current.getBoundingClientRect();

    setActiveRect({
      left: rect.left - containerRect.left,
      width: rect.width,
    });
  }, [location.pathname]);

  useLayoutEffect(() => {
    updateActiveRect();

    window.addEventListener("resize", updateActiveRect);
    return () => window.removeEventListener("resize", updateActiveRect);
  }, [updateActiveRect, items]);

  return (
    <div
      ref={containerRef}
      className="border-mauve-6 text-mauve-11 relative flex w-full gap-4 border-b px-2 pt-2 pb-1 text-sm"
    >
      {items.map((link) => (
        <NavLink
          key={link.href}
          to={link.href}
          className="relative z-10 px-2 py-1.5"
        >
          {({ isActive }) => (
            <span
              className={cn(
                "truncate pb-2",
                isActive && "text-violet-11 font-semibold",
              )}
            >
              {link.title}
            </span>
          )}
        </NavLink>
      ))}

      <motion.div
        className="bg-violet-11 absolute bottom-0 h-[3px] rounded-md"
        animate={{ left: activeRect.left, width: activeRect.width }}
        transition={{ type: "spring", stiffness: 300, damping: 30 }}
      />
    </div>
  );
}
