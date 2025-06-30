import { create } from "zustand";
import { persist } from "zustand/middleware";
import { authApi } from "@/api/auth";

type State = {
  isAuthenticated: boolean;
};

type Action = {
  login: () => void;
  logout: () => Promise<void>;
  setUnauthenticated: () => void;
};

export const useAuth = create<State & Action>()(
  persist(
    (set) => ({
      isAuthenticated: false,

      login: () => {
        set({ isAuthenticated: true });
      },

      logout: async () => {
        try {
          await authApi.logout();
        } finally {
          // Regardless of logout success, we set the state to unauthenticated
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
