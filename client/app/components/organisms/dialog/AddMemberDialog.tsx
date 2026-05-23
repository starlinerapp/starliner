import React, { useRef, useState } from "react";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { useTRPC } from "~/utils/trpc/react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { Cross } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";
import {
  formatInvalidEmailsError,
  getInvalidEmails,
  isValidEmail,
} from "~/utils/email";

interface FormInput {
  emails: string;
}

interface AddMemberDialogProps {
  organizationId: number;
  teamId?: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

function parseInviteEmails(raw: string): string[] {
  const emails = raw
    .trim()
    .split(/\s+/)
    .map((email) => email.trim())
    .filter(Boolean);

  const seen = new Set<string>();
  const uniqueEmails: string[] = [];

  for (const email of emails) {
    const key = email.toLowerCase();
    if (seen.has(key)) continue;
    seen.add(key);
    uniqueEmails.push(email);
  }

  return uniqueEmails;
}

function getRecognizedAndCurrent(raw: string): {
  recognized: string[];
  current: string;
} {
  if (!raw) {
    return { recognized: [], current: "" };
  }

  if (/\s$/.test(raw)) {
    return { recognized: parseInviteEmails(raw), current: "" };
  }

  const parts = raw.split(/\s+/).filter(Boolean);
  if (parts.length === 0) {
    return { recognized: [], current: "" };
  }

  const current = parts[parts.length - 1] ?? "";
  const recognizedPart = parts.slice(0, -1).join(" ");

  return {
    recognized: parseInviteEmails(recognizedPart),
    current,
  };
}

function buildEmailsValue(recognized: string[], current: string): string {
  const recognizedPart = recognized.join(" ");

  if (!current) {
    return recognizedPart ? `${recognizedPart} ` : "";
  }

  return recognizedPart ? `${recognizedPart} ${current}` : current;
}

export default function AddMemberDialog({
  organizationId,
  teamId,
  open,
  onOpenChange,
}: AddMemberDialogProps) {
  const trpc = useTRPC();
  const inputRef = useRef<HTMLInputElement>(null);
  const [validationError, setValidationError] = useState<string | null>(null);

  const { handleSubmit, reset, watch, setValue } = useForm<FormInput>({
    defaultValues: { emails: "" },
  });

  const emailsInput = watch("emails", "");
  const { recognized, current } = getRecognizedAndCurrent(emailsInput);
  const invalidRecognizedEmails = getInvalidEmails(recognized);
  const hasInvalidEmails = invalidRecognizedEmails.length > 0;
  const hasEmailsToInvite = recognized.length > 0 || current.trim().length > 0;

  const sendInviteMutation = useMutation(
    trpc.organization.sendInvite.mutationOptions(),
  );

  function updateEmails(nextRecognized: string[], nextCurrent: string) {
    setValue("emails", buildEmailsValue(nextRecognized, nextCurrent), {
      shouldDirty: true,
    });
    setValidationError(null);
    sendInviteMutation.reset();
  }

  function removeRecognizedEmail(emailToRemove: string) {
    updateEmails(
      recognized.filter(
        (email) => email.toLowerCase() !== emailToRemove.toLowerCase(),
      ),
      current,
    );
  }

  function handleInputKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    if (
      event.key !== "Backspace" ||
      current.length > 0 ||
      recognized.length === 0
    ) {
      return;
    }

    event.preventDefault();
    updateEmails(recognized.slice(0, -1), "");
  }

  function onInviteMember(data: FormInput) {
    const toEmails = parseInviteEmails(data.emails);

    if (toEmails.length === 0) {
      setValidationError("Enter at least one email address.");
      return;
    }

    const invalidEmails = getInvalidEmails(toEmails);
    if (invalidEmails.length > 0) {
      setValidationError(formatInvalidEmailsError(invalidEmails));
      return;
    }

    setValidationError(null);
    sendInviteMutation.mutate(
      {
        organizationId: organizationId,
        toEmails,
        inviteUrlPrefix: `${window.location.origin}/organizations/invite/`,
        ...(teamId != null ? { teamId } : {}),
      },
      {
        onSuccess: () => {
          reset();
          onOpenChange(false);
        },
      },
    );
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        onOpenChange(nextOpen);
        if (!nextOpen) {
          reset();
          setValidationError(null);
          sendInviteMutation.reset();
        }
      }}
    >
      <DialogContent>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h1>Invite Member</h1>
            <p className="text-mauve-11 text-sm">
              Invite members via email. They&apos;ll each receive a link to join
              your organization.
            </p>
          </div>
          {(validationError || sendInviteMutation.isError) && (
            <ErrorBanner
              text={validationError ?? sendInviteMutation.error?.message ?? ""}
            />
          )}
          <form
            className="flex flex-col gap-3"
            onSubmit={handleSubmit(onInviteMember)}
          >
            <div
              className="border-mauve-6 bg-gray-2 flex h-10 items-center gap-1.5 overflow-x-auto rounded-md border px-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
              onClick={() => inputRef.current?.focus()}
            >
              {recognized.map((email) => {
                const valid = isValidEmail(email);

                return (
                  <span
                    key={email}
                    className={cn(
                      "inline-flex h-6 max-w-48 shrink-0 items-center gap-1 rounded-md border px-2 text-xs font-medium shadow-sm",
                      valid
                        ? "border-mauve-6 text-mauve-12 bg-white"
                        : "border-red-6 bg-red-3 text-red-11",
                    )}
                    title={valid ? email : `${email} (invalid email)`}
                  >
                    <span className="truncate">{email}</span>
                    <button
                      type="button"
                      className={cn(
                        "shrink-0 rounded hover:bg-black/5",
                        valid ? "text-violet-11" : "text-red-11",
                      )}
                      aria-label={`Remove ${email}`}
                      onClick={(event) => {
                        event.stopPropagation();
                        removeRecognizedEmail(email);
                      }}
                    >
                      <Cross className="h-3.5 w-3.5 stroke-2" />
                    </button>
                  </span>
                );
              })}
              <input
                ref={inputRef}
                type="text"
                value={current}
                onChange={(event) =>
                  updateEmails(recognized, event.target.value)
                }
                onKeyDown={handleInputKeyDown}
                className="text-mauve-12 placeholder:text-mauve-11 h-full min-w-32 flex-1 bg-transparent text-sm outline-none"
                placeholder={
                  recognized.length === 0 && !current
                    ? "Enter one or more emails"
                    : undefined
                }
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button
                type="button"
                intent="secondary"
                className="w-24"
                onClick={() => {
                  onOpenChange(false);
                  reset();
                  setValidationError(null);
                  sendInviteMutation.reset();
                }}
              >
                Cancel
              </Button>
              <Button
                className="h-10 w-24"
                type="submit"
                disabled={
                  !hasEmailsToInvite ||
                  hasInvalidEmails ||
                  sendInviteMutation.isPending
                }
              >
                Invite
              </Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
}
