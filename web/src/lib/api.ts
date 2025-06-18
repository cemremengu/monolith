import type { User, CreateUserRequest } from "@/types";

const API_BASE = "/api";

export const api = {
  users: {
    getAll: async (): Promise<User[]> => {
      const response = await fetch(`${API_BASE}/users`);
      if (!response.ok) throw new Error("Failed to fetch users");
      return response.json();
    },

    getById: async (id: number): Promise<User> => {
      const response = await fetch(`${API_BASE}/users/${id}`);
      if (!response.ok) throw new Error("Failed to fetch user");
      return response.json();
    },

    create: async (data: CreateUserRequest): Promise<User> => {
      const response = await fetch(`${API_BASE}/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      if (!response.ok) throw new Error("Failed to create user");
      return response.json();
    },

    update: async (id: number, data: CreateUserRequest): Promise<User> => {
      const response = await fetch(`${API_BASE}/users/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      if (!response.ok) throw new Error("Failed to update user");
      return response.json();
    },

    delete: async (id: number): Promise<void> => {
      const response = await fetch(`${API_BASE}/users/${id}`, {
        method: "DELETE",
      });
      if (!response.ok) throw new Error("Failed to delete user");
    },
  },
};
