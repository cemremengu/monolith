import type { User } from "../users/types";
import type {
  LoginRequest,
  RegisterRequest,
  UpdatePreferencesRequest,
} from "./types";
import { API_BASE } from "../config";
import { httpClient } from "@/lib/http-client";

export const authApi = {
  login: async (data: LoginRequest): Promise<{ user: User }> => {
    return httpClient.post(`${API_BASE}/auth/login`, data);
  },

  register: async (data: RegisterRequest): Promise<{ user: User }> => {
    return httpClient.post(`${API_BASE}/auth/register`, data);
  },

  me: async (): Promise<User> => {
    return httpClient.get(`${API_BASE}/auth/me`);
  },

  logout: async (): Promise<void> => {
    return httpClient.post(`${API_BASE}/auth/logout`, undefined);
  },

  refresh: async (): Promise<void> => {
    return httpClient.post(`${API_BASE}/auth/refresh`, undefined);
  },

  updatePreferences: async (data: UpdatePreferencesRequest): Promise<User> => {
    return httpClient.patch(`${API_BASE}/auth/preferences`, data);
  },
};
