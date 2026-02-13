import { create } from "zustand";
import { persist } from "zustand/middleware";

import { authApi } from "@/features/auth/api";
import { accountApi } from "@/features/profile/api";
import type { User, LoginRequest } from "@/types/api";

type AuthState = {
  user: User | null;
  isLoggedIn: boolean;
};

type LogoutOptions = {
  redirectToLogin?: boolean;
};

type AuthActions = {
  login: (data: LoginRequest) => Promise<void>;
  logout: (options?: LogoutOptions) => Promise<void>;
  fetchUser: () => Promise<void>;
};

type AuthStore = AuthState & AuthActions;

export const useAuth = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      isLoggedIn: false,

      login: async (data: LoginRequest) => {
        await authApi.login(data);
        set({
          isLoggedIn: true,
        });
      },

      logout: async (options?: LogoutOptions) => {
        try {
          await authApi.logout();
        } finally {
          set({
            isLoggedIn: false,
            user: null,
          });
        }

        if (options?.redirectToLogin) {
          window.location.replace("/login");
        }
      },

      fetchUser: async () => {
        try {
          const user = await accountApi.profile();
          set({
            user,
          });
        } catch {
          set({
            user: null,
            isLoggedIn: false,
          });
        }
      },
    }),
    {
      name: "auth-store",
      partialize: (state) => ({ isLoggedIn: state.isLoggedIn }),
    },
  ),
);
