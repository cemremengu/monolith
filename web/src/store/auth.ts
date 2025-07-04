import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authApi } from "@/features/auth/api";
import { accountApi } from "@/features/profile/api";
import type { User, LoginRequest } from "@/types/api";

type AuthState = {
  user: User | null;
  isLoading: boolean;
  isLoggedIn: boolean;
};

type AuthActions = {
  login: (data: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  fetchUser: () => Promise<void>;
};

type AuthStore = AuthState & AuthActions;

export const useAuth = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      isLoading: false,
      isLoggedIn: false,

      login: async (data: LoginRequest) => {
        const response = await authApi.login(data);

        set({
          isLoggedIn: !!response.user,
          user: response.user,
          isLoading: false,
        });
      },

      logout: async () => {
        try {
          await authApi.logout();
        } finally {
          set({
            isLoggedIn: false,
            user: null,
            isLoading: false,
          });
        }

        window.location.replace("/login");
      },

      fetchUser: async () => {
        const { isLoading } = get();
        if (isLoading) return;

        set({ isLoading: true });
        try {
          const user = await accountApi.profile();
          set({
            user,
            isLoading: false,
            isLoggedIn: true,
          });
        } catch {
          set({
            user: null,
            isLoading: false,
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
