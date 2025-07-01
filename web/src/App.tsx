import { RouterProvider, createRouter } from "@tanstack/react-router";

// Import the generated route tree
import { routeTree } from "./routeTree.gen";
import { useAuth } from "./store/auth";
import { useUser } from "./store/user";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      // staleTime: 1000 * 60 * 5, // 5 minutes
      // retry: 1,
    },
  },
});

// Create a new router instance with context
const router = createRouter({
  routeTree,
  context: { auth: undefined!, user: undefined!, queryClient },
  defaultPreload: "intent",
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function App() {
  const auth = useAuth();
  const user = useUser();

  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} context={{ auth, user }} />;
    </QueryClientProvider>
  );
}

export default App;
