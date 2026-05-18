import { useCallback } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { accountApi } from "@/features/profile/api";
import { useAuth } from "@/hooks/use-auth";
import { usePreferences } from "@/hooks/use-preferences";
import { languages } from "@/i18n/language";

type Language = (typeof languages)[number];

type LanguageSwitcherProps = {
  value?: string;
  onChange?: (language: string) => void;
};

function LanguageItem({
  language,
  onSelect,
}: {
  language: Language;
  onSelect: (code: string) => void;
}) {
  const handleClick = useCallback(() => onSelect(language.code), [onSelect, language.code]);

  return (
    <DropdownMenuItem onClick={handleClick} className="cursor-pointer">
      <span className="mr-2">{language.flag}</span>
      {language.name}
    </DropdownMenuItem>
  );
}

export function LanguageSwitcher({ value, onChange }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();
  const { isLoggedIn, user } = useAuth();

  const changeLanguage = useCallback(
    async (languageCode: string) => {
      usePreferences.getState().setLanguage(languageCode);

      if (onChange) {
        onChange(languageCode);
        return;
      }

      if (isLoggedIn && user) {
        try {
          const { theme, timezone } = usePreferences.getState();
          await accountApi.updatePreferences({
            language: languageCode,
            theme,
            timezone,
          });
        } catch (error) {
          console.error("Failed to update language preference:", error);
        }
      }
    },
    [onChange, isLoggedIn, user],
  );

  const currentLanguageCode = value ?? i18n.language;
  const currentLanguage =
    languages.find((lang) => lang.code === currentLanguageCode) || languages[0];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="w-full justify-start">
          <span className="mr-2">{currentLanguage.flag}</span>
          {currentLanguage.name}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {languages.map((language) => (
          <LanguageItem key={language.code} language={language} onSelect={changeLanguage} />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
