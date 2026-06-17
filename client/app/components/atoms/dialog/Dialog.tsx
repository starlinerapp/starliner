import * as DialogPrimitive from "@radix-ui/react-dialog";
import * as React from "react";
import { Cross } from "~/components/atoms/icons";

export const DialogContent = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Content>
>(({ children, ...props }, forwardedRef) => (
  <DialogPrimitive.Portal>
    <DialogPrimitive.Overlay className="data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-998 bg-black-a6 backdrop-blur-[2px] duration-200 data-[state=closed]:animate-out data-[state=open]:animate-in" />
    <DialogPrimitive.Content
      className="data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 fixed top-1/2 left-1/2 z-999 max-h-[85vh] w-[90vw] max-w-150 -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6.25 shadow-xl duration-200 focus:outline-none data-[state=closed]:animate-out data-[state=open]:animate-in"
      {...props}
      ref={forwardedRef}
    >
      {children}
      <DialogPrimitive.Close className="cursor-pointer" asChild>
        <button
          type="button"
          className="absolute top-2.5 right-2.5 inline-flex size-6.25 appearance-none items-center justify-center rounded-md bg-white text-violet11 hover:bg-gray-4 focus:outline-none"
          aria-label="Close"
        >
          <Cross
            width={20}
            height={20}
            strokeWidth={2}
            className="stroke-gray-11"
          />
        </button>
      </DialogPrimitive.Close>
    </DialogPrimitive.Content>
  </DialogPrimitive.Portal>
));
DialogContent.displayName = "DialogContent";

export const Dialog = DialogPrimitive.Root;
export const DialogTrigger = DialogPrimitive.Trigger;
