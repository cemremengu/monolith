import { Monitor, Moon, Sun } from "lucide-react";
import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { usePreferences, type Theme } from "@/hooks/use-preferences";
import { languages } from "@/i18n/language";
import { cn } from "@/lib/utils";

import { useUpdatePreferences } from "./api/queries";

const themeOptions: { value: Theme; icon: typeof Monitor }[] = [
  { value: "system", icon: Monitor },
  { value: "light", icon: Sun },
  { value: "dark", icon: Moon },
];

const languageItems = languages.map((language) => ({
  value: language.code,
  label: language.name,
}));

export function PreferencesMenu() {
  const { t } = useTranslation();
  const theme = usePreferences((state) => state.theme);
  const language = usePreferences((state) => state.language);
  const updatePreferences = useUpdatePreferences();

  const persist = useCallback(
    async (next: { theme?: Theme; language?: string }) => {
      const current = usePreferences.getState();
      try {
        await updatePreferences.mutateAsync({
          theme: next.theme ?? current.theme,
          language: next.language ?? current.language,
          timezone: current.timezone,
        });
        toast.success(t("profile.messages.updateSuccess"));
      } catch (error) {
        console.error("Failed to update preference:", error);
        toast.error(t("profile.updateError"));
      }
    },
    [updatePreferences, t],
  );

  const handleThemeChange = useCallback(
    (value: Theme) => {
      usePreferences.getState().setTheme(value);
      void persist({ theme: value });
    },
    [persist],
  );

  const handleLanguageChange = useCallback(
    (value: string | null) => {
      if (!value) {
        return;
      }
      usePreferences.getState().setLanguage(value);
      void persist({ language: value });
    },
    [persist],
  );

  return (
    <div className="px-1 py-1.5">
      <div className="px-1.5 pb-1 text-xs font-medium text-muted-foreground">
        {t("profile.preferences.title")}
      </div>

      <div className="flex items-center justify-between gap-2 px-1.5 py-1">
        <span className="text-sm">{t("profile.preferences.theme")}</span>
        <div className="flex items-center gap-0.5 rounded-md border bg-background p-0.5">
          {themeOptions.map((option) => (
            <button
              key={option.value}
              type="button"
              aria-label={option.value}
              aria-pressed={theme === option.value}
              onClick={() => handleThemeChange(option.value)}
              className={cn(
                "flex size-6 items-center justify-center rounded-sm text-muted-foreground transition-colors hover:text-foreground",
                theme === option.value && "bg-accent text-foreground",
              )}
            >
              <option.icon className="size-3.5" />
            </button>
          ))}
        </div>
      </div>

      <div className="flex items-center justify-between gap-2 px-1.5 py-1">
        <span className="text-sm">{t("profile.preferences.language")}</span>
        <Select items={languageItems} value={language} onValueChange={handleLanguageChange}>
          <SelectTrigger size="sm" className="w-auto gap-1 px-2 text-xs data-[size=sm]:h-7.5">
            <SelectValue />
          </SelectTrigger>
          <SelectContent align="end">
            {languageItems.map((item) => (
              <SelectItem key={item.value} value={item.value}>
                {item.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
