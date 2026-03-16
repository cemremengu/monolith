import { useEffect } from "react";
import { createFileRoute, redirect, Outlet } from "@tanstack/react-router";

import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useSessionRotation } from "@/hooks/use-session-rotation";
import { usePreferences } from "@/hooks/use-preferences";
import { useAuth } from "@/hooks/use-auth";
import { AppSidebar } from "@/components/sidebar";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async ({ context, location }) => {
    if (!context.auth.isLoggedIn) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }

    if (!context.auth.user) {
      await context.auth.fetchUser();
    }
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  useSessionRotation();
  const { user } = useAuth();
  const syncPreferences = usePreferences((s) => s.syncPreferences);

  useEffect(() => {
    if (user) {
      syncPreferences({
        theme: (user.theme as "light" | "dark" | "system") || "system",
        language: user.language || "en-US",
        timezone: user.timezone || "UTC",
      });
    }
  }, [user, syncPreferences]);

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="flex-1">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
