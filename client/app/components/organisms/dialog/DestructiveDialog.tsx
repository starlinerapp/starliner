import * as React from "react";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

interface DestructiveDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  bannerText: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  isPending?: boolean;
  onConfirm: () => void;
}

export default function DestructiveDialog({
  open,
  onOpenChange,
  title,
  bannerText,
  description,
  confirmLabel = "Delete",
  cancelLabel = "Cancel",
  isPending = false,
  onConfirm,
}: DestructiveDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h1>{title}</h1>
            <ErrorBanner text={bannerText} />
            <p className="text-mauve-11 text-sm">{description}</p>
          </div>
          <div className="flex justify-end gap-2">
            <Button
              type="button"
              intent="secondary"
              className="w-24"
              disabled={isPending}
              onClick={() => onOpenChange(false)}
            >
              {cancelLabel}
            </Button>
            <Button
              className="w-24"
              intent="primary"
              onClick={onConfirm}
              disabled={isPending}
            >
              {confirmLabel}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
