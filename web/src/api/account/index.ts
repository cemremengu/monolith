import type { User } from "../users/types";
import type { UpdatePreferencesRequest } from "./types";
import { httpClient } from "@/lib/http-client";

export const accountApi = {
  profile: (): Promise<User> => {
    return httpClient.get(`account/profile`);
  },

  updatePreferences: (data: UpdatePreferencesRequest): Promise<User> => {
    return httpClient.patch(`account/preferences`, data);
  },

  sessions: (): Promise<unknown[]> => {
    return httpClient.get(`account/sessions`);
  },

  revokeSession: (sessionId: string): Promise<void> => {
    return httpClient.delete(`account/sessions/${sessionId}`);
  },

  revokeAllOtherSessions: (): Promise<{
    message: string;
    revokedCount: number;
  }> => {
    return httpClient.post(`account/sessions/revoke-others`);
  },
};
