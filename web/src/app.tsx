import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider, createRouter } from "@tanstack/react-router";

import { NotFound } from "@/components/not-found";
import { Spinner } from "@/components/ui/spinner";

import { TooltipProvider } from "./components/ui/tooltip";
import { useAuth } from "./hooks/use-auth";
import { ThemeProvider } from "./hooks/use-theme";
import { routeTree } from "./routeTree.gen";

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
      <Spinner className="size-10 text-primary" />
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

export function App() {
  return (
    <ThemeProvider>
      <TooltipProvider>
        <QueryClientProvider client={queryClient}>
          <RouterProvider router={router} context={{ auth: useAuth() }} />
        </QueryClientProvider>
      </TooltipProvider>
    </ThemeProvider>
  );
}
