import { Link, useNavigate } from "@tanstack/react-router";
import { Info, LogOut, Settings, UserRoundIcon } from "lucide-react";
import { useCallback, useState } from "react";
import { useTranslation } from "react-i18next";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { AboutDialog } from "@/features/about/about-dialog";
import { PreferencesMenu } from "@/features/profile/preferences-menu";
import { useAuth } from "@/hooks/use-auth";

export function NavUser({
  user,
}: {
  user: {
    name: string;
    email: string;
    avatar: string;
    isAdmin: boolean;
  };
}) {
  const { t } = useTranslation();
  const { logout } = useAuth();
  const navigate = useNavigate();
  const [aboutOpen, setAboutOpen] = useState(false);

  const handleLogout = useCallback(async () => {
    await logout();
    void navigate({ to: "/login" });
  }, [logout, navigate]);

  const handleAboutOpen = useCallback(() => setAboutOpen(true), []);

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu modal={false}>
          <DropdownMenuTrigger
            render={
              <SidebarMenuButton className="group-data-[collapsible=icon]:justify-center hover:cursor-pointer data-popup-open:bg-sidebar-accent data-popup-open:text-sidebar-accent-foreground" />
            }
          >
            <Avatar className="size-4">
              <AvatarImage src={user.avatar} alt={user.name} />
              <AvatarFallback className="text-xs">{getInitials(user.name)}</AvatarFallback>
            </Avatar>
            <span className="flex-1 truncate font-medium group-data-[collapsible=icon]:hidden">
              {user.name}
            </span>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--anchor-width) min-w-56 rounded-lg"
            side="top"
            align="start"
            sideOffset={8}
          >
            <div className="grid px-2 py-1.5 text-left text-sm leading-tight">
              <span className="truncate font-medium">{user.name}</span>
              <span className="truncate text-xs text-muted-foreground">{user.email}</span>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuGroup className="space-y-1">
              <DropdownMenuItem render={<Link to="/profile" />}>
                <UserRoundIcon />
                {t("sidebar.user.profile")}
              </DropdownMenuItem>
              <DropdownMenuItem
                render={<Link to="/settings" search={{ section: "preferences" }} />}
              >
                <Settings />
                {t("settings.title")}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleAboutOpen}>
                <Info />
                {t("sidebar.user.about")}
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <PreferencesMenu />
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout}>
              <LogOut />
              {t("sidebar.user.signOut")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <AboutDialog open={aboutOpen} onOpenChange={setAboutOpen} />
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
