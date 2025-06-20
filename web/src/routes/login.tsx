import { createFileRoute, Link } from "@tanstack/react-router";
import { LoginForm } from "@/components/login-form";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  const navigate = Route.useNavigate();

  const handleLoginSuccess = () => {
    navigate({ to: "/" });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-md">
        <LoginForm onSuccess={handleLoginSuccess} />
        <div className="mt-4 text-center">
          <p className="text-sm text-gray-600">
            Don't have an account?{" "}
            <Link to="/register" className="text-blue-600 hover:underline">
              Register here
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}