import { createFileRoute, redirect, Outlet } from "@tanstack/react-router";

import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useSessionRotation } from "@/hooks/use-session-rotation";
import { AppSidebar } from "@/components/sidebar";
import i18next from "@/i18n";

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

    const lang = context.auth.user?.language;
    if (lang && lang !== i18next.language) {
      i18next.changeLanguage(lang);
    }
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  useSessionRotation();

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
