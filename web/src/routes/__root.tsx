import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { Toaster } from "@/components/ui/sonner";
import type { QueryClient } from "@tanstack/react-query";

type MinimalUser = {
  id: string;
  username: string;
  email: string;
  avatar?: string;
};

type RouterContext = {
  queryClient: QueryClient;
  auth: {
    isAuthenticated: boolean;
    user: MinimalUser | null;
    isLoading: boolean;
    login: () => void;
    logout: () => Promise<void>;
    setUnauthenticated: () => void;
    fetchUser: () => Promise<void>;
    setUser: (user: MinimalUser) => void;
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
