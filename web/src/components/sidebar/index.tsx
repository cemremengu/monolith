import { useLocation } from "@tanstack/react-router";
import { HomeIcon } from "lucide-react";
import * as React from "react";

import { Logo } from "@/components/logo";
import { Menu } from "@/components/sidebar/menu";
import { SettingsNav } from "@/components/sidebar/settings-nav";
import { NavUser } from "@/components/sidebar/user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useAuth } from "@/hooks/use-auth";

const navMain = [
  {
    title: "Dashboard",
    url: "/dashboard",
    icon: HomeIcon,
    isActive: true,
  },
];

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth();
  const { pathname } = useLocation();
  const inSettings = pathname.startsWith("/settings");

  const userData = React.useMemo(
    () => ({
      name: user?.username || "User",
      email: user?.email || "user@example.com",
      avatar: user?.avatar || "",
      isAdmin: user?.isAdmin || false,
    }),
    [user],
  );

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="gap-3.5">
        <Logo />
      </SidebarHeader>
      <SidebarContent>{inSettings ? <SettingsNav /> : <Menu items={navMain} />}</SidebarContent>
      <SidebarFooter>
        <NavUser user={userData} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
