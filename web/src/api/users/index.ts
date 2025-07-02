import type { User, CreateUserRequest } from "./types";
import { httpClient } from "@/lib/http-client";

export const usersApi = {
  getAll: async (params: {
    filterBy?: string;
    sortBy?: "name" | "email";
  }): Promise<User[]> => {
    const response = await httpClient.get("/users", { params });
    return response.data;
  },

  getById: async (id: string): Promise<User> => {
    const response = await httpClient.get(`/users/${id}`);
    return response.data;
  },

  create: async (data: CreateUserRequest): Promise<User> => {
    const response = await httpClient.post("/users", data);
    return response.data;
  },

  update: async (id: string, data: CreateUserRequest): Promise<User> => {
    const response = await httpClient.put(`/users/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await httpClient.delete(`/users/${id}`);
  },
};
