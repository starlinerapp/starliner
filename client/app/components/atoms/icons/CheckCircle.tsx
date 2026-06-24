import type { SVGProps } from "react";
import * as React from "react";

const SvgCheckCircle = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={24}
    height={24}
    fill="none"
    stroke="currentColor"
    strokeLinecap="round"
    strokeLinejoin="round"
    strokeWidth={2}
    className="check-circle_svg__lucide check-circle_svg__lucide-circle-check-icon check-circle_svg__lucide-circle-check"
    {...props}
  >
    <circle cx={12} cy={12} r={10} />
    <path d="m9 12 2 2 4-4" />
  </svg>
);
export default SvgCheckCircle;
