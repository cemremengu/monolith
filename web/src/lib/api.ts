import type {
  User,
  CreateUserRequest,
  LoginRequest,
  RegisterRequest,
} from "@/types";

const API_BASE = "/api";

function getHeaders(): HeadersInit {
  return {
    "Content-Type": "application/json",
  };
}

export const api = {
  auth: {
    login: async (data: LoginRequest): Promise<{ user: User }> => {
      const response = await fetch(`${API_BASE}/auth/login`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(data),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to login");
      return response.json();
    },

    register: async (data: RegisterRequest): Promise<{ user: User }> => {
      const response = await fetch(`${API_BASE}/auth/register`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(data),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to register");
      return response.json();
    },

    me: async (): Promise<User> => {
      const response = await fetch(`${API_BASE}/auth/me`, {
        headers: getHeaders(),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to fetch user profile");
      return response.json();
    },

    logout: async (): Promise<void> => {
      const response = await fetch(`${API_BASE}/auth/logout`, {
        method: "POST",
        headers: getHeaders(),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to logout");
    },
  },

  users: {
    getAll: async (): Promise<User[]> => {
      const response = await fetch(`${API_BASE}/users`, {
        headers: getHeaders(),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to fetch users");
      return response.json();
    },

    getById: async (id: string): Promise<User> => {
      const response = await fetch(`${API_BASE}/users/${id}`, {
        headers: getHeaders(),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to fetch user");
      return response.json();
    },

    create: async (data: CreateUserRequest): Promise<User> => {
      const response = await fetch(`${API_BASE}/users`, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(data),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to create user");
      return response.json();
    },

    update: async (id: string, data: CreateUserRequest): Promise<User> => {
      const response = await fetch(`${API_BASE}/users/${id}`, {
        method: "PUT",
        headers: getHeaders(),
        body: JSON.stringify(data),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to update user");
      return response.json();
    },

    delete: async (id: string): Promise<void> => {
      const response = await fetch(`${API_BASE}/users/${id}`, {
        method: "DELETE",
        headers: getHeaders(),
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to delete user");
    },
  },
};
