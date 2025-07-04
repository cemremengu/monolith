import type { User, CreateUserRequest } from "@/types/api";
import { httpClient } from "@/lib/http-client";

export const usersApi = {
  getAll: (params: {
    filterBy?: string;
    sortBy?: "name" | "email";
  }): Promise<User[]> => {
    return httpClient.get(`users`, params);
  },

  getById: (id: string): Promise<User> => {
    return httpClient.get(`users/${id}`);
  },

  create: (data: CreateUserRequest): Promise<User> => {
    return httpClient.post(`users`, data);
  },

  update: (id: string, data: CreateUserRequest): Promise<User> => {
    return httpClient.put(`users/${id}`, data);
  },

  delete: (id: string): Promise<void> => {
    return httpClient.delete(`users/${id}`);
  },
};
