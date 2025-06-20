import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { useAuth } from "@/lib/auth";
import { Button } from "@/components/ui/button";
import { useEffect } from "react";

function Root() {
  const { user, isAuthenticated, logout, checkAuth } = useAuth();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const handleLogout = async () => {
    await logout();
  };

  return (
    <>
      <div className="p-4 border-b">
        <div className="flex justify-between items-center">
          <div className="flex gap-4">
            <Link to="/" className="[&.active]:font-bold">
              Home
            </Link>
            {isAuthenticated && (
              <Link to="/users" className="[&.active]:font-bold">
                Users
              </Link>
            )}
          </div>
          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <>
                <span className="text-sm text-gray-600">
                  Welcome, {user?.name || user?.username}
                </span>
                <Button variant="outline" size="sm" onClick={handleLogout}>
                  Logout
                </Button>
              </>
            ) : (
              <div className="flex gap-2">
                <Link to="/login">
                  <Button variant="outline" size="sm">
                    Login
                  </Button>
                </Link>
                <Link to="/register">
                  <Button variant="default" size="sm">
                    Register
                  </Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
      <Outlet />
      <TanStackRouterDevtools />
    </>
  );
}

export const Route = createRootRoute({
  component: Root,
});
