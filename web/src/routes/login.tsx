import { createFileRoute, redirect } from "@tanstack/react-router";
import { LoginForm } from "@/components/login-form";
import { z } from "zod";

const loginSearchSchema = z.object({
  redirect: z.string().optional(),
});

export const Route = createFileRoute("/login")({
  validateSearch: loginSearchSchema,
  beforeLoad: ({ context, search }) => {
    if (context.auth.isLoggedIn) {
      throw redirect({
        to: search.redirect || "/dashboard",
      });
    }
  },
  component: LoginPage,
});

function LoginPage() {
  const navigate = Route.useNavigate();
  const search = Route.useSearch();

  const handleLoginSuccess = () => {
    navigate({ to: search.redirect || "/dashboard" });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="w-full max-w-md">
        <LoginForm onSuccess={handleLoginSuccess} />
      </div>
    </div>
  );
}
