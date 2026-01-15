import type {
  User,
  CreateUserRequest,
  InviteUsersRequest,
  InviteUsersResponse,
} from "@/types/api";
import { httpClient } from "@/lib/http-client";

export const usersApi = {
  getAll: (params: {
    filterBy?: string;
    sortBy?: "name" | "email";
  }): Promise<User[]> => {
    return httpClient.get(`accounts`, params);
  },

  getById: (id: string): Promise<User> => {
    return httpClient.get(`accounts/${id}`);
  },

  create: (data: CreateUserRequest): Promise<User> => {
    return httpClient.post(`accounts`, data);
  },

  invite: (data: InviteUsersRequest): Promise<InviteUsersResponse> => {
    return httpClient.post(`accounts/invite`, data);
  },

  update: (id: string, data: CreateUserRequest): Promise<User> => {
    return httpClient.put(`accounts/${id}`, data);
  },

  disable: (id: string): Promise<void> => {
    return httpClient.patch(`accounts/${id}/disable`);
  },

  enable: (id: string): Promise<void> => {
    return httpClient.patch(`accounts/${id}/enable`);
  },

  delete: (id: string): Promise<void> => {
    return httpClient.delete(`accounts/${id}`);
  },
};
