import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authApi } from "@/api/auth";
import type { User } from "@/api/users/types";

type AuthState = {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (user: User) => void;
  logout: () => Promise<void>;
  setUnauthenticated: () => void;
};

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,

      login: (user: User) => {
        set({ user, isAuthenticated: true, isLoading: false });
      },

      logout: async () => {
        try {
          await authApi.logout();
        } catch {
          // Ignore logout API errors, clear state anyway
        } finally {
          set({ user: null, isAuthenticated: false, isLoading: false });
        }
      },

      setUnauthenticated: () => {
        set({ user: null, isAuthenticated: false, isLoading: false });
      },
    }),
    {
      name: "auth-store",
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
);
