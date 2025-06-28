import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authApi } from "@/api/auth";

type AuthState = {
  isAuthenticated: boolean;
  login: () => void;
  logout: () => Promise<void>;
  setUnauthenticated: () => void;
};

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      isAuthenticated: false,

      login: () => {
        set({ isAuthenticated: true });
      },

      logout: async () => {
        try {
          await authApi.logout();
        } catch {
          // Ignore logout API errors, clear state anyway
        } finally {
          set({ isAuthenticated: false });
        }
      },

      setUnauthenticated: () => {
        set({ isAuthenticated: false });
      },
    }),
    {
      name: "auth-store",
    },
  ),
);
