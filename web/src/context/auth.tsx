import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { authApi } from "@/api/auth";
import { accountApi } from "@/api/account";
import type { User } from "@/api/users/types";

type AuthContextValue = {
  isAuthenticated: boolean;
  user: User | null;
  isLoading: boolean;
  login: () => void;
  logout: () => Promise<void>;
  setUnauthenticated: () => void;
  fetchUser: () => Promise<void>;
  setUser: (user: User) => void;
  clearUser: () => void;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

type AuthProviderProps = {
  children: ReactNode;
  storageKey?: string;
};

export function AuthProvider({
  children,
  storageKey = "auth",
}: AuthProviderProps) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    try {
      const stored = localStorage.getItem(storageKey);
      return stored === "true";
    } catch {
      return false;
    }
  });
  const [user, setUserState] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Persist authentication state to localStorage
  useEffect(() => {
    localStorage.setItem(storageKey, isAuthenticated.toString());
  }, [isAuthenticated, storageKey]);

  const fetchUser = async () => {
    if (isLoading) return;

    setIsLoading(true);
    try {
      const user = await accountApi.profile();

      setUserState(user);
    } catch {
      setUserState(null);
    } finally {
      setIsLoading(false);
    }
  };

  const setUser = (user: User) => {
    setUserState(user);
    setIsLoading(false);
  };

  const clearUser = () => {
    setUserState(null);
    setIsLoading(false);
  };

  const login = () => {
    setIsAuthenticated(true);
  };

  const logout = async () => {
    try {
      await authApi.logout();
    } finally {
      // Regardless of logout success, we set the state to unauthenticated
      setIsAuthenticated(false);
      clearUser();
    }
  };

  const setUnauthenticated = () => {
    setIsAuthenticated(false);
    clearUser();
  };

  const value = {
    isAuthenticated,
    user,
    isLoading,
    login,
    logout,
    setUnauthenticated,
    fetchUser,
    setUser,
    clearUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export const useAuth = () => {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
};

// Utility function to get auth state without React context (for http-client usage)
export const getAuthState = () => {
  try {
    const stored = localStorage.getItem("auth-store");
    return {
      isAuthenticated: stored === "true",
      setUnauthenticated: () => {
        localStorage.setItem("auth-store", "false");
      },
    };
  } catch {
    // Silently fail if localStorage is not available
    return {
      isAuthenticated: false,
      setUnauthenticated: () => {
        try {
          localStorage.setItem("auth-store", "false");
        } catch {
          // Silently fail if localStorage is not available
        }
      },
    };
  }
};
