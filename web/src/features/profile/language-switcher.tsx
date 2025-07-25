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
import { languages } from "@/i18n/language";

type LanguageSwitcherProps = {
  value?: string;
  onChange?: (language: string) => void;
};

export function LanguageSwitcher({ value, onChange }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();
  const { isLoggedIn, user } = useAuth();

  const changeLanguage = async (languageCode: string) => {
    i18n.changeLanguage(languageCode);

    if (onChange) {
      onChange(languageCode);
      return;
    }

    if (isLoggedIn && user) {
      try {
        await accountApi.updatePreferences({
          language: languageCode,
          theme: "system",
          timezone: "UTC",
        });
      } catch (error) {
        console.error("Failed to update language preference:", error);
      }
    }
  };

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
          <DropdownMenuItem
            key={language.code}
            onClick={() => changeLanguage(language.code)}
            className="cursor-pointer"
          >
            <span className="mr-2">{language.flag}</span>
            {language.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
