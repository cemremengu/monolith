import type { User, CreateUserRequest } from "./types";
import { API_BASE } from "../config";
import { httpClient } from "@/lib/http-client";

export const usersApi = {
  getAll: (): Promise<User[]> => {
    return httpClient.get(`${API_BASE}/users`);
  },

  getById: (id: string): Promise<User> => {
    return httpClient.get(`${API_BASE}/users/${id}`);
  },

  create: (data: CreateUserRequest): Promise<User> => {
    return httpClient.post(`${API_BASE}/users`, data);
  },

  update: (id: string, data: CreateUserRequest): Promise<User> => {
    return httpClient.put(`${API_BASE}/users/${id}`, data);
  },

  delete: (id: string): Promise<void> => {
    return httpClient.delete(`${API_BASE}/users/${id}`);
  },
};
