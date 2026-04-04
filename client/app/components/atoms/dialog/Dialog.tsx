import * as React from "react";
import * as DialogPrimitive from "@radix-ui/react-dialog";
import { Cross } from "~/components/atoms/icons";

export const DialogContent = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Content>
>(({ children, ...props }, forwardedRef) => (
  <DialogPrimitive.Portal>
    <DialogPrimitive.Overlay className="bg-black-a6 fixed inset-0" />
    <DialogPrimitive.Content
      className="fixed top-1/2 left-1/2 max-h-[85vh] w-[90vw] max-w-[600px] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-[25px] shadow-xl focus:outline-none"
      {...props}
      ref={forwardedRef}
    >
      {children}
      <DialogPrimitive.Close className="cursor-pointer" asChild>
        <button
          className="text-violet11 hover:bg-gray-4 absolute top-2.5 right-2.5 inline-flex size-[25px] appearance-none items-center justify-center rounded-md bg-white focus:outline-none"
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
