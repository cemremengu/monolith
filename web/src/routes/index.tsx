import { createFileRoute, Link } from "@tanstack/react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const { user, isAuthenticated } = useAuth();

  return (
    <div className="p-6 flex flex-col items-center justify-center gap-6">
      <Card className="max-w-md mx-auto">
        <CardHeader>
          <CardTitle>
            {isAuthenticated ? `Welcome back, ${user?.name || user?.username}!` : "Welcome to Monolith"}
          </CardTitle>
          <CardDescription>
            A full-stack application with JWT authentication, built with Go, Echo, PostgreSQL, React, and shadcn/ui.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {isAuthenticated ? (
            <>
              <p className="text-sm text-muted-foreground">
                You are successfully logged in! Explore the features available to you.
              </p>
              <div className="flex gap-2">
                <Link to="/users">
                  <Button>View Users</Button>
                </Link>
              </div>
            </>
          ) : (
            <>
              <p className="text-sm text-muted-foreground">
                Please log in or register to access the full functionality.
              </p>
              <div className="flex gap-2">
                <Link to="/login">
                  <Button>Login</Button>
                </Link>
                <Link to="/register">
                  <Button variant="outline">Register</Button>
                </Link>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      {isAuthenticated && (
        <Card className="max-w-md mx-auto">
          <CardHeader>
            <CardTitle>Features</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="text-sm space-y-1">
              <li>✅ JWT Authentication with HTTP-only cookies</li>
              <li>✅ Login with username or email</li>
              <li>✅ User management (CRUD)</li>
              <li>✅ Secure password hashing</li>
              <li>✅ Protected routes</li>
              <li>✅ Responsive UI with Tailwind CSS</li>
            </ul>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
