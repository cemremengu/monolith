import type { User, CreateUserRequest } from "./types";
import { API_BASE, getHeaders } from "../config";

export const usersApi = {
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
};
