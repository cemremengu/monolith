import { Loader2 } from "lucide-react";

const Spinner = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <Loader2 className="w-10 h-10 text-primary animate-spin" />
    </div>
  );
};

Spinner.displayName = "Spinner";

export { Spinner };
