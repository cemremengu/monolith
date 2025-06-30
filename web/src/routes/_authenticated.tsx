import {
  createFileRoute,
  redirect,
  Outlet,
  useNavigate,
} from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { AppSidebar } from "@/components/app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: ({ context, location }) => {
    if (!context.auth.isAuthenticated) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  const { user, auth } = Route.useRouteContext();
  const { i18n } = useTranslation();
  const navigate = useNavigate();

  // Fetch user data when authenticated but user is not loaded
  useEffect(() => {
    if (auth.isAuthenticated && !user.user && !user.isLoading) {
      user.fetchUser();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [auth.isAuthenticated, user.user, user.isLoading]);

  // Load user's language preference when authenticated
  useEffect(() => {
    if (
      auth.isAuthenticated &&
      user.user?.language &&
      i18n.language !== user.user.language
    ) {
      i18n.changeLanguage(user.user.language);
    }
  }, [auth.isAuthenticated, user.user?.language, i18n]);

  const handleLogout = async () => {
    await auth.logout();
    user.clearUser();
    navigate({ to: "/login" });
  };

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mr-2 data-[orientation=vertical]:h-4"
          />
          <div className="flex-1" />
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-600">
              Welcome, {user.user?.name || user.user?.username}
            </span>
            <Button variant="outline" size="sm" onClick={handleLogout}>
              Logout
            </Button>
          </div>
        </header>
        <div className="flex-1">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
