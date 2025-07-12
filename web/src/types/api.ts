export type LoginRequest = {
  login: string;
  password: string;
};

export type UpdatePreferencesRequest = {
  language?: string;
  theme?: string;
  timezone?: string;
};

export type RegisterRequest = {
  username: string;
  email: string;
  password: string;
  name: string;
};

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
