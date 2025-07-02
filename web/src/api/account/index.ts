import type { User } from "../users/types";
import type { RegisterRequest, UpdatePreferencesRequest } from "./types";
import { httpClient } from "@/lib/http-client";

export const accountApi = {
  profile: async (): Promise<User> => {
    const response = await httpClient.get("/account/profile");
    return response.data;
  },

  updatePreferences: async (data: UpdatePreferencesRequest): Promise<User> => {
    const response = await httpClient.patch("/account/preferences", data);
    return response.data;
  },

  sessions: async (): Promise<unknown[]> => {
    const response = await httpClient.get("/account/sessions");
    return response.data;
  },

  revokeSession: async (sessionId: string): Promise<void> => {
    await httpClient.delete(`/account/sessions/${sessionId}`);
  },

  register: async (data: RegisterRequest): Promise<{ user: User }> => {
    const response = await httpClient.post("/auth/register", data);
    return response.data;
  },

  revokeAllOtherSessions: async (): Promise<{
    message: string;
    revokedCount: number;
  }> => {
    const response = await httpClient.post("/account/sessions/revoke-others");
    return response.data;
  },
};
