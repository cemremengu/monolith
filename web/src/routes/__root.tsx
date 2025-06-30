import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
// import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { ThemeProvider } from "@/context/theme";
import { Toaster } from "@/components/ui/sonner";
import type { User } from "@/api/users/types";
import type { QueryClient } from "@tanstack/react-query";

type RouterContext = {
  queryClient: QueryClient;
  auth: {
    isAuthenticated: boolean;
    login: () => void;
    logout: () => Promise<void>;
    setUnauthenticated: () => void;
  };
  user: {
    user: User | null;
    isLoading: boolean;
    fetchUser: () => Promise<void>;
    setUser: (user: User) => void;
    clearUser: () => void;
  };
};

function Root() {
  return (
    <ThemeProvider>
      <Outlet />
      <Toaster />
      {/* <TanStackRouterDevtools /> */}
    </ThemeProvider>
  );
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: Root,
});
