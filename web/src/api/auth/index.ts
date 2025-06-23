import type { User } from "../users/types";
import type { LoginRequest, RegisterRequest } from "./types";
import { API_BASE } from "../config";
import { httpClient } from "@/lib/http-client";

export const authApi = {
  login: async (data: LoginRequest): Promise<{ user: User }> => {
    return httpClient.post(`${API_BASE}/auth/login`, data, { skipAuth: true });
  },

  register: async (data: RegisterRequest): Promise<{ user: User }> => {
    return httpClient.post(`${API_BASE}/auth/register`, data, { skipAuth: true });
  },

  me: async (): Promise<User> => {
    return httpClient.get(`${API_BASE}/auth/me`);
  },

  logout: async (): Promise<void> => {
    return httpClient.post(`${API_BASE}/auth/logout`, undefined, { skipAuth: true });
  },

  refresh: async (): Promise<void> => {
    return httpClient.post(`${API_BASE}/auth/refresh`, undefined, { skipAuth: true });
  },
};
