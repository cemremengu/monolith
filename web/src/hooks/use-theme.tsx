/* eslint-disable react-refresh/only-export-components */
import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  type ReactNode,
} from "react";

import { usePreferences, type Theme } from "./use-preferences";

type ThemeProviderContextValue = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  isDarkTheme: boolean;
};

const ThemeProviderContext = createContext<
  ThemeProviderContextValue | undefined
>(undefined);

type ThemeProviderProps = {
  children: ReactNode;
};

export function ThemeProvider({ children }: ThemeProviderProps) {
  const theme = usePreferences((s) => s.theme);
  const setTheme = usePreferences((s) => s.setTheme);

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

  const isDarkTheme = useMemo(() => {
    if (theme === "dark") {
      return true;
    }
    if (theme === "light") {
      return false;
    }
    return window.matchMedia("(prefers-color-scheme: dark)").matches;
  }, [theme]);

  const value = useMemo(
    () => ({ theme, setTheme, isDarkTheme }),
    [theme, setTheme, isDarkTheme],
  );

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
