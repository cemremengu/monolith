import { createFileRoute, redirect } from "@tanstack/react-router";
import { LoginPage } from "@/features/auth/login-page";
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
  component: LoginRouteComponent,
});

function LoginRouteComponent() {
  const navigate = Route.useNavigate();
  const search = Route.useSearch();

  const handleLoginSuccess = () => {
    navigate({ to: search.redirect || "/dashboard" });
  };

  return <LoginPage onSuccess={handleLoginSuccess} />;
}
