import { useCallback, useEffect, useRef, useState } from "react";
import { throttle } from "throttle-debounce";

export default function useMouseLeave() {
  const [mouseLeft, setMouseLeft] = useState(true);
  const elementRef = useRef<HTMLElement | null>(null);

  const handleMouseMove = useRef(
    throttle(50, (e: MouseEvent) => {
      if (!elementRef?.current) return;

      const rect = elementRef.current.getBoundingClientRect();

      if (
        e.clientX < rect.left ||
        e.clientX > rect.right ||
        e.clientY < rect.top ||
        e.clientY > rect.bottom
      ) {
        setMouseLeft(true);
      } else {
        setMouseLeft(false);
      }
    }),
  ).current;

  const handleMouseEnter = useRef(() => {
    setMouseLeft(false);
    window.addEventListener("mousemove", handleMouseMove);
  }).current;

  const setRef = useCallback(
    (node: HTMLElement | null) => {
      if (elementRef?.current) {
        elementRef.current.removeEventListener("mouseenter", handleMouseEnter);
      }

      if (node !== null) {
        node.addEventListener("mouseenter", handleMouseEnter);

        elementRef.current = node;
      }
    },
    [handleMouseEnter],
  );

  useEffect(() => {
    if (mouseLeft) {
      window.removeEventListener("mousemove", handleMouseMove);
    }
  }, [mouseLeft, handleMouseMove]);

  useEffect(() => {
    return () => {
      if (elementRef?.current) {
        elementRef.current.removeEventListener("mouseenter", handleMouseEnter);
      }
      window.removeEventListener("mousemove", handleMouseMove);
    };
  }, [handleMouseMove, handleMouseEnter]);

  return [mouseLeft, setRef, elementRef] as const;
}
