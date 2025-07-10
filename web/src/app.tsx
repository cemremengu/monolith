import { RouterProvider, createRouter } from "@tanstack/react-router";

import { routeTree } from "./routeTree.gen";
import { useAuth } from "./store/auth";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Loading } from "@/components/loading";
import { NotFound } from "@/components/not-found";
import { ThemeProvider } from "@/context/theme";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: false,
    },
  },
});

// Create a new router instance with context
const router = createRouter({
  routeTree,
  context: { auth: undefined!, queryClient },
  defaultPreload: "intent",
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
  defaultPendingComponent: Loading,
  defaultNotFoundComponent: () => <NotFound />,
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function Router() {
  return (
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} context={{ auth: useAuth() }} />
      </QueryClientProvider>
    </ThemeProvider>
  );
}

export default Router;
