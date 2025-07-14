import { Loader2 } from "lucide-react";

const Loading = () => {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center">
      <Loader2 className="text-primary h-10 w-10 animate-spin" />
    </div>
  );
};

Loading.displayName = "Loading";

export { Loading };
