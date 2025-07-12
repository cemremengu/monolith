import { useTranslation } from "react-i18next";
import { profileQueryOptions, useUpdatePreferences } from "../api/queries";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Globe, Palette, Clock } from "lucide-react";
import { ThemeSwitcher } from "./theme-switcher";
import { LanguageSwitcher } from "./language-switcher";
import { TimezoneSelector } from "./timezone-selector";
import { useTheme } from "@/context/theme";
import { toast } from "sonner";
import { useSuspenseQuery } from "@tanstack/react-query";

export function SettingsPage() {
  const { data: user } = useSuspenseQuery(profileQueryOptions);
  const { t } = useTranslation();
  const updatePreferences = useUpdatePreferences();
  const { setTheme } = useTheme();

  const handleLanguageChange = async (language: string) => {
    try {
      await updatePreferences.mutateAsync({
        language,
        theme: user?.theme || "system",
        timezone: user?.timezone || "UTC",
      });
      toast.success(t("profile.messages.updateSuccess"));
    } catch (error) {
      console.error("Failed to update language preference:", error);
      toast.error(t("profile.messages.updateError"));
    }
  };

  const handleThemeChange = async (theme: string) => {
    try {
      setTheme(theme as "light" | "dark" | "system");
      await updatePreferences.mutateAsync({
        language: user?.language || "en",
        theme,
        timezone: user?.timezone || "UTC",
      });
      toast.success(t("profile.messages.updateSuccess"));
    } catch (error) {
      console.error("Failed to update theme preference:", error);
      toast.error(t("profile.messages.updateError"));
    }
  };

  const handleTimezoneChange = async (timezone: string) => {
    try {
      await updatePreferences.mutateAsync({
        language: user?.language || "en",
        theme: user?.theme || "system",
        timezone,
      });
      toast.success(t("profile.messages.updateSuccess"));
    } catch (error) {
      console.error("Failed to update timezone preference:", error);
      toast.error(t("profile.messages.updateError"));
    }
  };

  return (
    <div className="p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">
            {t("settings.title", "Settings")}
          </h1>
          <p className="text-muted-foreground">
            {t("settings.subtitle", "Manage your preferences and settings")}
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>{t("profile.preferences.title")}</CardTitle>
            <CardDescription>
              {t(
                "settings.preferences.subtitle",
                "Customize your language, theme, and timezone preferences",
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Globe className="h-4 w-4" />
                  <label className="text-sm font-medium">
                    {t("profile.preferences.language")}
                  </label>
                </div>
                <LanguageSwitcher
                  value={user?.language || "en"}
                  onChange={handleLanguageChange}
                />
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Palette className="h-4 w-4" />
                  <label className="text-sm font-medium">
                    {t("profile.preferences.theme")}
                  </label>
                </div>
                <ThemeSwitcher
                  value={user?.theme || "system"}
                  onChange={handleThemeChange}
                />
              </div>
              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <Clock className="h-4 w-4" />
                  <label className="text-sm font-medium">
                    {t("profile.preferences.timezone")}
                  </label>
                </div>
                <TimezoneSelector
                  value={user?.timezone || "UTC"}
                  onChange={handleTimezoneChange}
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
