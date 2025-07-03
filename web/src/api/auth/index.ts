import type { User } from "../users/types";
import type { LoginRequest } from "./types";
import { httpClient } from "@/lib/http-client";

export const authApi = {
  login: (data: LoginRequest): Promise<{ user: User }> => {
    return httpClient.post(`auth/login`, data);
  },

  logout: (): Promise<void> => {
    return httpClient.post(`auth/logout`, undefined);
  },

  refresh: (): Promise<void> => {
    return httpClient.post(`auth/refresh`, undefined);
  },
};
