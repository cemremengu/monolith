import type { User, CreateUserRequest } from "./types";
import { API_BASE } from "../config";
import { httpClient } from "@/lib/http-client";

export const usersApi = {
  getAll: async (): Promise<User[]> => {
    return httpClient.get(`${API_BASE}/users`);
  },

  getById: async (id: string): Promise<User> => {
    return httpClient.get(`${API_BASE}/users/${id}`);
  },

  create: async (data: CreateUserRequest): Promise<User> => {
    return httpClient.post(`${API_BASE}/users`, data);
  },

  update: async (id: string, data: CreateUserRequest): Promise<User> => {
    return httpClient.put(`${API_BASE}/users/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return httpClient.delete(`${API_BASE}/users/${id}`);
  },
};
