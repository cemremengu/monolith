import * as React from "react";
import { Home, Users, Settings2, BarChart3, User } from "lucide-react";
import { useAuth } from "@/lib/auth";

import { NavMain } from "@/components/nav-main";
import { NavUser } from "@/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth();

  const navMain = [
    {
      title: "Dashboard",
      url: "/dashboard",
      icon: Home,
      isActive: true,
    },
    {
      title: "Users",
      url: "/users",
      icon: Users,
    },
    {
      title: "Analytics",
      url: "#",
      icon: BarChart3,
      items: [
        {
          title: "Overview",
          url: "#",
        },
        {
          title: "Reports",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "#",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "#",
        },
        {
          title: "Security",
          url: "#",
        },
      ],
    },
  ];

  const userData = {
    name: user?.name || user?.username || "User",
    email: user?.email || "user@example.com",
    avatar: user?.avatar || "",
  };

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <div className="flex items-center gap-2 px-2 py-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <User className="h-4 w-4" />
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-semibold">My App</span>
            <span className="truncate text-xs">Dashboard</span>
          </div>
        </div>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={userData} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
