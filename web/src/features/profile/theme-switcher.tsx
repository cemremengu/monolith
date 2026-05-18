import { Monitor, Moon, Sun, type LucideIcon } from "lucide-react";
import { useCallback } from "react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useTheme } from "@/hooks/use-theme";

type ThemeSwitcherProps = {
  value?: string;
  onChange?: (theme: string) => void;
};

function ThemeOption({
  value,
  icon: Icon,
  label,
  onSelect,
}: {
  value: string;
  icon: LucideIcon;
  label: string;
  onSelect: (theme: string) => void;
}) {
  const handleClick = useCallback(() => onSelect(value), [onSelect, value]);

  return (
    <DropdownMenuItem onClick={handleClick}>
      <Icon className="mr-2 h-4 w-4" />
      <span>{label}</span>
    </DropdownMenuItem>
  );
}

export function ThemeSwitcher({ value, onChange }: ThemeSwitcherProps) {
  const { theme: contextTheme, setTheme } = useTheme();
  const theme = value ?? contextTheme;

  const handleSelect = useCallback(
    (selected: string) => {
      if (onChange) {
        onChange(selected);
      } else {
        setTheme(selected as "light" | "dark" | "system");
      }
    },
    [onChange, setTheme],
  );

  const getThemeIcon = () => {
    switch (theme) {
      case "light":
        return <Sun className="h-4 w-4" />;
      case "dark":
        return <Moon className="h-4 w-4" />;
      case "system":
        return <Monitor className="h-4 w-4" />;
      default:
        return <Monitor className="h-4 w-4" />;
    }
  };

  const getThemeLabel = () => {
    switch (theme) {
      case "light":
        return "Light";
      case "dark":
        return "Dark";
      case "system":
        return "System";
      default:
        return "System";
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="w-full justify-start gap-2">
          {getThemeIcon()}
          {getThemeLabel()}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <ThemeOption value="light" icon={Sun} label="Light" onSelect={handleSelect} />
        <ThemeOption value="dark" icon={Moon} label="Dark" onSelect={handleSelect} />
        <ThemeOption value="system" icon={Monitor} label="System" onSelect={handleSelect} />
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
