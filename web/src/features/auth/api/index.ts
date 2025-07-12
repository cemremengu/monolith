import type { LoginRequest } from "@/types/api";
import { httpClient } from "@/lib/http-client";

export const authApi = {
  login: (data: LoginRequest): Promise<{ message: string }> => {
    return httpClient.post(`auth/login`, data);
  },

  logout: (): Promise<void> => {
    return httpClient.post(`auth/logout`);
  },
};
