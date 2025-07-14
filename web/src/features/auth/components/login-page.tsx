import { LoginForm } from "./login-form";

interface LoginPageProps {
  onSuccess: () => void;
}

export function LoginPage({ onSuccess }: LoginPageProps) {
  return (
    <div className="bg-background flex min-h-screen items-center justify-center">
      <div className="w-full max-w-md">
        <LoginForm onSuccess={onSuccess} />
      </div>
    </div>
  );
}
