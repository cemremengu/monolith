import { create } from "zustand";
import { persist } from "zustand/middleware";
import { changeLanguage } from "i18next";

export type Theme = "dark" | "light" | "system";

type Preferences = {
  theme: Theme;
  language: string;
  timezone: string;
};

type PreferencesActions = {
  setTheme: (theme: Theme) => void;
  setLanguage: (language: string) => void;
  setTimezone: (timezone: string) => void;
  syncPreferences: (prefs: Preferences) => void;
};

type PreferencesStore = Preferences & PreferencesActions;

export const usePreferences = create<PreferencesStore>()(
  persist(
    (set) => ({
      theme: "system",
      language: "en-US",
      timezone: "UTC",

      setTheme: (theme: Theme) => {
        set({ theme });
      },

      setLanguage: (language: string) => {
        changeLanguage(language);
        set({ language });
      },

      setTimezone: (timezone: string) => {
        set({ timezone });
      },

      syncPreferences: (prefs: Preferences) => {
        const { theme, language, timezone } = usePreferences.getState();
        if (
          theme === prefs.theme &&
          language === prefs.language &&
          timezone === prefs.timezone
        )
          return;
        changeLanguage(prefs.language);
        set(prefs);
      },
    }),
    {
      name: "preferences",
    },
  ),
);
