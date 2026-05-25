import * as React from "react";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";

interface DestructiveDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  onConfirm: () => void;
}

export default function DestructiveDialog({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel = "Delete",
  cancelLabel = "Cancel",
  onConfirm,
}: DestructiveDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h1>{title}</h1>
            <p className="text-mauve-11 text-sm">{description}</p>
          </div>
          <div className="flex justify-end gap-2">
            <Button
              type="button"
              intent="secondary"
              className="w-24"
              onClick={() => onOpenChange(false)}
            >
              {cancelLabel}
            </Button>
            <Button className="w-24" intent="danger" onClick={onConfirm}>
              {confirmLabel}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
