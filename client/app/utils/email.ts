import { z } from "zod";

const emailSchema = z.email();

export function isValidEmail(email: string): boolean {
  return emailSchema.safeParse(email.trim()).success;
}

export function getInvalidEmails(emails: string[]): string[] {
  return emails.filter((email) => !isValidEmail(email));
}

export function formatInvalidEmailsError(invalidEmails: string[]): string {
  if (invalidEmails.length === 1) {
    return `"${invalidEmails[0]}" is not a valid email address.`;
  }

  return `Invalid email addresses: ${invalidEmails.join(", ")}.`;
}
