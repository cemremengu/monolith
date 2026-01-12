import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { NotFound } from "@/components/not-found";
import { Spinner } from "@/components/ui/spinner";

import { useAuth } from "./hooks/use-auth";
import { routeTree } from "./routeTree.gen";
import { ThemeProvider } from "./hooks/use-theme";

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
  // defaultPreload: "intent",
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
  defaultPendingComponent: () => (
    <div className="flex min-h-screen flex-col items-center justify-center">
      <Spinner className="text-primary size-10" />
    </div>
  ),
  defaultNotFoundComponent: () => <NotFound />,
  defaultErrorComponent: ({ error }) => (
    <div>
      <h2>Something went wrong!</h2>
      <p>{error.message}</p>
    </div>
  ),
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
