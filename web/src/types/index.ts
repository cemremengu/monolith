export interface User {
  id: string;
  username: string;
  email: string;
  name?: string;
  isAdmin: boolean;
  language?: string;
  theme?: string;
  timezone?: string;
  lastSeenAt?: string;
  isDisabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserRequest {
  username: string;
  name: string;
  email: string;
}

export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  name: string;
}