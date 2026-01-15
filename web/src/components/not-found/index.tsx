import { Home, RotateCcw } from "lucide-react";
import { Link } from "@tanstack/react-router";

import { Button } from "@/components/ui/button";

export function NotFound() {
  const handleGoBack = () => {
    window.history.back();
  };

  return (
    <div className="bg-background flex min-h-screen items-center justify-center p-4">
      <div className="w-full max-w-md space-y-8 text-center">
        {/* Error Icon */}

        {/* Error Message */}
        <div className="space-y-4">
          <h1 className="text-foreground text-6xl font-bold">404</h1>
          <h2 className="text-foreground text-2xl font-semibold">
            Page Not Found
          </h2>
          <p className="text-muted-foreground">
            The page you're looking for doesn't exist or has been moved.
          </p>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col justify-center gap-3 sm:flex-row">
          <Button onClick={handleGoBack} variant="outline">
            <RotateCcw className="size-4" />
            Go Back
          </Button>
          <Link to="/">
            <Button>
              <Home className="size-4" />
              Start Over
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
