import { Link, useSearch } from "@tanstack/react-router";
import {
  Brain,
  ChevronLeft,
  KeyRound,
  SlidersHorizontal,
  Users,
  type LucideIcon,
} from "lucide-react";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import type { SettingsSection } from "@/features/profile/settings-page";
import { useAuth } from "@/hooks/use-auth";

type NavItem = { id: SettingsSection; label: string; icon: LucideIcon };
type NavGroup = { label: string; items: NavItem[] };

export function SettingsNav() {
  const { t } = useTranslation();
  const { user } = useAuth();
  const { section } = useSearch({ strict: false });

  const groups = useMemo<NavGroup[]>(() => {
    const result: NavGroup[] = [
      {
        label: t("settings.account"),
        items: [
          {
            id: "preferences",
            label: t("profile.preferences.title"),
            icon: SlidersHorizontal,
          },
        ],
      },
    ];
    if (user?.isAdmin) {
      result.push({
        label: t("admin.title"),
        items: [
          { id: "users", label: t("admin.users.title"), icon: Users },
          { id: "providers", label: t("aiProviders.title"), icon: Brain },
          { id: "license", label: t("license.title"), icon: KeyRound },
        ],
      });
    }
    return result;
  }, [t, user?.isAdmin]);

  return (
    <>
      <div className="px-2 py-1 group-data-[collapsible=icon]:hidden">
        <Button
          render={<Link to="/dashboard" />}
          variant="ghost"
          className="relative w-full justify-center"
        >
          <ChevronLeft className="absolute left-2" />
          <span className="font-semibold">{t("settings.title")}</span>
        </Button>
      </div>
      {groups.map((group) => (
        <SidebarGroup key={group.label}>
          <SidebarGroupLabel>{group.label}</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              {group.items.map((item) => (
                <SidebarMenuItem key={item.id}>
                  <SidebarMenuButton
                    render={<Link to="/settings" search={{ section: item.id }} />}
                    tooltip={item.label}
                    isActive={section === item.id}
                  >
                    <item.icon />
                    <span>{item.label}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      ))}
    </>
  );
}
