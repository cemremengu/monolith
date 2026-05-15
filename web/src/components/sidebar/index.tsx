import { HomeIcon } from "lucide-react";
import * as React from "react";

import { Logo } from "@/components/logo";
import { Menu } from "@/components/sidebar/menu";
import { NavUser } from "@/components/sidebar/user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useAuth } from "@/hooks/use-auth";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth();

  const navMain = [
    {
      title: "Dashboard",
      url: "/dashboard",
      icon: HomeIcon,
      isActive: true,
    },
  ];

  const userData = {
    name: user?.username || "User",
    email: user?.email || "user@example.com",
    avatar: user?.avatar || "",
    isAdmin: user?.isAdmin || false,
  };

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="gap-3.5">
        <Logo />
      </SidebarHeader>
      <SidebarContent>
        <Menu items={navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={userData} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
