import type { User } from "../users/types";
import type { LoginRequest } from "./types";
import { httpClient } from "@/lib/http-client";

export const authApi = {
  login: async (data: LoginRequest): Promise<{ user: User }> => {
    const response = await httpClient.post("/auth/login", data);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await httpClient.post("/auth/logout");
  },

  refresh: async (): Promise<void> => {
    await httpClient.post("/auth/refresh");
  },
};
