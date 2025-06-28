import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authApi } from "@/api/auth";
import type { User } from "@/api/users/types";

type UserState = {
  user: User | null;
  isLoading: boolean;
  fetchUser: () => Promise<void>;
  setUser: (user: User) => void;
  clearUser: () => void;
};

export const useUser = create<UserState>()(
  persist(
    (set, get) => ({
      user: null,
      isLoading: false,

      fetchUser: async () => {
        if (get().isLoading) return;

        set({ isLoading: true });
        try {
          const user = await authApi.me();
          set({ user, isLoading: false });
        } catch {
          set({ user: null, isLoading: false });
        }
      },

      setUser: (user: User) => {
        set({ user, isLoading: false });
      },

      clearUser: () => {
        set({ user: null, isLoading: false });
      },
    }),
    {
      name: "user-store",
      partialize: (state) => ({
        id: state.user?.id,
        username: state.user?.username,
        email: state.user?.email,
      }),
    },
  ),
);
