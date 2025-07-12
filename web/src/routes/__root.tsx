import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { Toaster } from "@/components/ui/sonner";
import type { QueryClient } from "@tanstack/react-query";
import type { User } from "@/types/api";

type RouterContext = {
  queryClient: QueryClient;
  auth: {
    isLoggedIn: boolean;
    user: User | null;
    fetchUser: () => Promise<void>;
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
