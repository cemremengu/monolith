import { create } from "zustand";
import { authApi } from "@/api/auth";
import type { User } from "@/api/users/types";

type State = {
  user: User | null;
  isLoading: boolean;
};

type Action = {
  fetchUser: () => Promise<void>;
  setUser: (user: User) => void;
  clearUser: () => void;
};

export const useUser = create<State & Action>()((set, get) => ({
  user: null,
  isLoading: false,

  fetchUser: async () => {
    if (get().isLoading) return;

    set({ isLoading: true });
    try {
      const user = await authApi.profile();
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
}));
