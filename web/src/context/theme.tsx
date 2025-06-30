import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { authApi } from "@/api/auth";
import { useAuth } from "@/store/auth";
import { useUser } from "@/store/user";

type Theme = "dark" | "light" | "system";

type ThemeProviderContextValue = {
  theme: Theme;
  setTheme: (theme: Theme) => Promise<void>;
};

const ThemeProviderContext = createContext<
  ThemeProviderContextValue | undefined
>(undefined);

type ThemeProviderProps = {
  children: ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "theme",
}: ThemeProviderProps) {
  const { isAuthenticated } = useAuth();
  const { user } = useUser();
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window !== "undefined") {
      return (localStorage.getItem(storageKey) as Theme) || defaultTheme;
    }
    return defaultTheme;
  });

  // Load theme from user preferences when authenticated
  useEffect(() => {
    if (isAuthenticated && user?.theme) {
      const userTheme = user.theme as Theme;
      setTheme(userTheme);

      // Update localStorage when user preferences change (e.g., from profile save)
      if (typeof window !== "undefined") {
        localStorage.setItem(storageKey, userTheme);
      }
    }
  }, [isAuthenticated, user?.theme, storageKey]);

  useEffect(() => {
    const root = window.document.documentElement;

    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light";

      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(theme);
  }, [theme]);

  const value = {
    theme,
    setTheme: async (newTheme: Theme) => {
      setTheme(newTheme);

      // Always save to localStorage for fallback
      if (typeof window !== "undefined") {
        localStorage.setItem(storageKey, newTheme);
      }

      // Save to backend if authenticated
      if (isAuthenticated) {
        try {
          await authApi.updatePreferences({ theme: newTheme });
        } catch (error) {
          console.error("Failed to update theme preference:", error);
        }
      }
    },
  };

  return (
    <ThemeProviderContext.Provider value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }

  return context;
};
