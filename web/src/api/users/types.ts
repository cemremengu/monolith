export type User = {
  id: string;
  username: string;
  email: string;
  name?: string;
  avatar?: string;
  isAdmin: boolean;
  language?: string;
  theme?: string;
  timezone?: string;
  lastSeenAt?: string;
  isDisabled: boolean;
  createdAt: string;
  updatedAt: string;
};

export type CreateUserRequest = {
  username: string;
  name: string;
  email: string;
};
