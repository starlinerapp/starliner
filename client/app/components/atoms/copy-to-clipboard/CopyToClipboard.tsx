import React, { useEffect, useState } from "react";
import useMouseLeave from "~/hooks/useMouseLeave";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "~/components/atoms/tooltip/Tooltip";
import { cn } from "~/utils/cn";

interface CopyToClipboardProps {
  text: string;
  className?: string;
}

const CopyToClipboard = ({ text, className }: CopyToClipboardProps) => {
  const [mouseLeft, ref] = useMouseLeave();
  const [copied, setCopied] = useState(false);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (mouseLeft) {
      setCopied(false);
      setOpen(false);
    }
  }, [mouseLeft]);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setOpen(true);
      setTimeout(() => setOpen(false), 750);
    } catch (err) {
      console.error("Failed to copy text: ", err);
    }
  };

  return (
    <TooltipProvider>
      <Tooltip open={open} onOpenChange={setOpen} delayDuration={300}>
        <TooltipTrigger ref={ref}>
          <p
            onClick={handleCopy}
            className={cn(
              "hover:bg-gray-4 flex cursor-pointer flex-row gap-1 rounded-md px-2 align-middle",
              className,
            )}
          >
            {text}
          </p>
        </TooltipTrigger>
        <TooltipContent>
          <p>{copied ? "Copied" : "Click to copy"}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default CopyToClipboard;
