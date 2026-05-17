import type { User } from "../entity/user";

export type UserProfile = {
  id: string;
  name: string;
  email: string;
  image: string | null;
};

export function toUserProfile(user: User): UserProfile {
  return {
    id: user.id,
    name: user.name,
    email: user.email,
    image: user.image,
  };
}

export function toUserProfiles(users: User[]): UserProfile[] {
  return users.map(toUserProfile);
}
