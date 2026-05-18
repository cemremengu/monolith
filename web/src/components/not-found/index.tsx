import { Link } from "@tanstack/react-router";
import { Home, RotateCcw } from "lucide-react";

import { Button } from "@/components/ui/button";

function handleGoBack() {
  window.history.back();
}

export function NotFound() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <div className="w-full max-w-md space-y-8 text-center">
        {/* Error Icon */}

        {/* Error Message */}
        <div className="space-y-4">
          <h1 className="text-6xl font-bold text-foreground">404</h1>
          <h2 className="text-2xl font-semibold text-foreground">Page Not Found</h2>
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
