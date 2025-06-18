import { createFileRoute } from "@tanstack/react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useState } from "react";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const [count, setCount] = useState(0);

  return (
    <div className="p-6 flex flex-col items-center justify-center gap-2">
      <Card className="max-w-md mx-auto">
        <CardHeader>
          <CardTitle>Welcome to My App Template</CardTitle>
          <CardDescription>
            A simple template app built with Go, Echo, PostgreSQL, React, and
            shadcn/ui.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Navigate to the Users page to see the CRUD functionality in action.
          </p>
        </CardContent>
      </Card>

      <Button
        onClick={() => setCount((count) => count + 1)}
        className="hover:cursor-pointer"
      >
        count is {count}
      </Button>
    </div>
  );
}
