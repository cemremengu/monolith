import { LoginForm } from "./login-form";

interface LoginPageProps {
  onSuccess: () => void;
}

export function LoginPage({ onSuccess }: LoginPageProps) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-md">
        <LoginForm onSuccess={onSuccess} />
      </div>
    </div>
  );
}
