import type { User } from "../users/types";
import type {
  LoginRequest,
  RegisterRequest,
  UpdatePreferencesRequest,
} from "./types";
import { API_BASE } from "../config";
import { httpClient } from "@/lib/http-client";

export const authApi = {
  login: (data: LoginRequest): Promise<{ user: User }> => {
    return httpClient.post(`${API_BASE}/auth/login`, data);
  },

  register: (data: RegisterRequest): Promise<{ user: User }> => {
    return httpClient.post(`${API_BASE}/auth/register`, data);
  },

  logout: (): Promise<void> => {
    return httpClient.post(`${API_BASE}/auth/logout`, undefined);
  },

  refresh: (): Promise<void> => {
    return httpClient.post(`${API_BASE}/auth/refresh`, undefined);
  },

  updatePreferences: (data: UpdatePreferencesRequest): Promise<User> => {
    return httpClient.patch(`${API_BASE}/auth/preferences`, data);
  },
};
