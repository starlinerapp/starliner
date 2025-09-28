import React, { useState, useEffect, useRef } from "react";
import { NavLink, useLocation } from "react-router";
import { motion } from "framer-motion";
import { cn } from "~/utils/cn";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

type NavigationBarItem = {
  title: string;
  href: string;
};

interface NavigationBarProps {
  title: string;
  items: NavigationBarItem[];
  isLoading: boolean;
}

export default function NavigationBar({
  title,
  items,
  isLoading,
}: NavigationBarProps) {
  const location = useLocation();

  const [activeRect, setActiveRect] = useState({ left: 0, width: 0 });
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    const activeLink = containerRef.current.querySelector<HTMLAnchorElement>(
      `a[href="${location.pathname}"] span`,
    );

    if (activeLink) {
      const rect = activeLink.getBoundingClientRect();
      const containerRect = containerRef.current.getBoundingClientRect();
      setActiveRect({
        left: rect.left - containerRect.left,
        width: rect.width,
      });
    }
  }, [location.pathname]);

  return (
    <div className="bg-violet-1">
      {isLoading ? (
        <div className="px-4 pt-4 pb-3">
          <Skeleton className="h-7 w-32" />
        </div>
      ) : (
        <h1 className="text-mauve-12 px-4 pt-4 pb-3 text-xl font-bold">
          {title}
        </h1>
      )}

      <div
        ref={containerRef}
        className="border-mauve-6 text-mauve-11 relative flex w-full gap-4 border-b px-2 pb-1 text-sm"
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
                  "pb-2",
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
    </div>
  );
}
