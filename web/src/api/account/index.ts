import type { User } from "../users/types";
import type { UpdatePreferencesRequest } from "./types";
import { API_BASE } from "../config";
import { httpClient } from "@/lib/http-client";

export const accountApi = {
  profile: (): Promise<User> => {
    return httpClient.get(`${API_BASE}/account/profile`);
  },

  updatePreferences: (data: UpdatePreferencesRequest): Promise<User> => {
    return httpClient.patch(`${API_BASE}/account/preferences`, data);
  },

  sessions: (): Promise<unknown[]> => {
    return httpClient.get(`${API_BASE}/account/sessions`);
  },

  revokeSession: (sessionId: string): Promise<void> => {
    return httpClient.delete(`${API_BASE}/account/sessions/${sessionId}`);
  },

  revokeAllOtherSessions: (): Promise<{
    message: string;
    revokedCount: number;
  }> => {
    return httpClient.post(
      `${API_BASE}/account/sessions/revoke-others`,
      undefined,
    );
  },
};
