import { useSuspenseQuery } from "@tanstack/react-query";
import { Monitor, Moon, Sun } from "lucide-react";
import { useCallback, type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { PageHeader } from "@/components/page-header";
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

import { profileQueryOptions, useUpdatePreferences } from "./api/queries";
import { TimezoneSelector } from "./timezone-selector";

const themeOptions: { value: Theme; icon: typeof Monitor }[] = [
  { value: "system", icon: Monitor },
  { value: "light", icon: Sun },
  { value: "dark", icon: Moon },
];

const languageItems = languages.map((language) => ({
  value: language.code,
  label: language.name,
}));

function SettingRow({
  title,
  description,
  control,
}: {
  title: string;
  description: string;
  control: ReactNode;
}) {
  return (
    <div className="flex items-center justify-between gap-4 px-4 py-4">
      <div className="space-y-0.5">
        <div className="text-sm font-medium">{title}</div>
        <div className="text-sm text-muted-foreground">{description}</div>
      </div>
      {control}
    </div>
  );
}

export function PreferencesSection() {
  useSuspenseQuery(profileQueryOptions);
  const { t } = useTranslation();
  const theme = usePreferences((state) => state.theme);
  const language = usePreferences((state) => state.language);
  const timezone = usePreferences((state) => state.timezone);
  const updatePreferences = useUpdatePreferences();

  const persist = useCallback(
    async (next: { theme?: Theme; language?: string; timezone?: string }) => {
      const current = usePreferences.getState();
      try {
        await updatePreferences.mutateAsync({
          theme: next.theme ?? current.theme,
          language: next.language ?? current.language,
          timezone: next.timezone ?? current.timezone,
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

  const handleTimezoneChange = useCallback(
    (value: string) => {
      usePreferences.getState().setTimezone(value);
      void persist({ timezone: value });
    },
    [persist],
  );

  return (
    <div className="space-y-6">
      <PageHeader title={t("profile.preferences.title")} />

      <section className="space-y-3">
        <h2 className="text-sm font-semibold text-muted-foreground">
          {t("profile.preferences.interfaceAndTheme")}
        </h2>
        <div className="divide-y rounded-xl border bg-card">
          <SettingRow
            title={t("profile.preferences.theme")}
            description={t("profile.preferences.themeDescription")}
            control={
              <div className="flex items-center gap-0.5 rounded-md border bg-background p-0.5">
                {themeOptions.map((option) => (
                  <ThemeOption
                    key={option.value}
                    value={option.value}
                    icon={option.icon}
                    active={theme === option.value}
                    onSelect={handleThemeChange}
                  />
                ))}
              </div>
            }
          />
          <SettingRow
            title={t("profile.preferences.language")}
            description={t("profile.preferences.languageDescription")}
            control={
              <Select items={languageItems} value={language} onValueChange={handleLanguageChange}>
                <SelectTrigger className="w-auto gap-2">
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
            }
          />
          <SettingRow
            title={t("profile.preferences.timezone")}
            description={t("profile.preferences.timezoneDescription")}
            control={<TimezoneSelector value={timezone} onChange={handleTimezoneChange} />}
          />
        </div>
      </section>
    </div>
  );
}

function ThemeOption({
  value,
  icon: Icon,
  active,
  onSelect,
}: {
  value: Theme;
  icon: typeof Monitor;
  active: boolean;
  onSelect: (value: Theme) => void;
}) {
  const handleClick = useCallback(() => onSelect(value), [onSelect, value]);

  return (
    <button
      type="button"
      aria-label={value}
      aria-pressed={active}
      onClick={handleClick}
      className={cn(
        "flex size-7 items-center justify-center rounded-sm text-muted-foreground transition-colors hover:text-foreground",
        active && "bg-accent text-foreground",
      )}
    >
      <Icon className="size-4" />
    </button>
  );
}
