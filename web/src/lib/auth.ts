import { create } from "zustand";
import type { User } from "@/types";
import { api } from "./api";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  login: (user: User) => void;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

export const useAuth = create<AuthState>()((set) => ({
  user: null,
  isAuthenticated: false,

  login: (user: User) => {
    set({ user, isAuthenticated: true });
  },

  logout: async () => {
    try {
      await api.auth.logout();
    } catch (error) {
      // Ignore logout API errors, clear state anyway
    } finally {
      set({ user: null, isAuthenticated: false });
    }
  },

  checkAuth: async () => {
    try {
      const user = await api.auth.me();
      set({ user, isAuthenticated: true });
    } catch (error) {
      set({ user: null, isAuthenticated: false });
    }
  },
}));