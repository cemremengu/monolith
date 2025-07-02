import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { Toaster } from "@/components/ui/sonner";
import type { QueryClient } from "@tanstack/react-query";
import type { User } from "@/api/users/types";

type RouterContext = {
  queryClient: QueryClient;
  auth: {
    isAuthenticated: boolean;
    user: User | null;
    isLoading: boolean;
    login: () => void;
    logout: () => Promise<void>;
    setUnauthenticated: () => void;
    fetchUser: () => Promise<void>;
    setUser: (user: User) => void;
    clearUser: () => void;
  };
};

function Root() {
  return (
    <>
      <Outlet />
      <Toaster />
    </>
  );
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: Root,
});
