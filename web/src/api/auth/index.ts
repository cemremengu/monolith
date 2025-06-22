import type { User } from "../users/types";
import type { LoginRequest, RegisterRequest } from "./types";
import { API_BASE, getHeaders } from "../config";

export const authApi = {
  login: async (data: LoginRequest): Promise<{ user: User }> => {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(data),
      credentials: "include",
    });
    if (!response.ok) throw new Error("Failed to login");
    return response.json();
  },

  register: async (data: RegisterRequest): Promise<{ user: User }> => {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(data),
      credentials: "include",
    });
    if (!response.ok) throw new Error("Failed to register");
    return response.json();
  },

  me: async (): Promise<User> => {
    const response = await fetch(`${API_BASE}/auth/me`, {
      headers: getHeaders(),
      credentials: "include",
    });
    if (!response.ok) throw new Error("Failed to fetch user profile");
    return response.json();
  },

  logout: async (): Promise<void> => {
    const response = await fetch(`${API_BASE}/auth/logout`, {
      method: "POST",
      headers: getHeaders(),
      credentials: "include",
    });
    if (!response.ok) throw new Error("Failed to logout");
  },
};
