import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Globe } from "lucide-react";
import { authApi } from "@/api/auth";
import { useAuth } from "@/lib/auth";

const languages = [
  { code: "en", name: "English", flag: "🇺🇸" },
  { code: "tr", name: "Türkçe", flag: "🇹🇷" },
];

type LanguageSwitcherProps = {
  value?: string;
  onChange?: (language: string) => void;
};

export function LanguageSwitcher({ value, onChange }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();
  const { isAuthenticated } = useAuth();

  const changeLanguage = async (languageCode: string) => {
    if (onChange) {
      onChange(languageCode);
      return;
    }

    i18n.changeLanguage(languageCode);

    // Save to backend if authenticated
    if (isAuthenticated) {
      try {
        await authApi.updatePreferences({ language: languageCode });
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
          <Globe className="h-4 w-4 mr-2" />
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
