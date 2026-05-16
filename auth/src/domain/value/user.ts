import type { User } from "~/domain/entity/user";

export type UserProfile = {
  id: string;
  name: string;
  email: string;
};

export function toUserProfile(user: User): UserProfile {
  return {
    id: user.id,
    name: user.name,
    email: user.email,
  };
}

export function toUserProfiles(users: User[]): UserProfile[] {
  return users.map(toUserProfile);
}
