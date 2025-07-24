import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Info } from "lucide-react";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";

import { getVersionInfo } from "../api";

type AboutDialogProps = {
  trigger?: React.ReactNode;
};

export function AboutDialog({ trigger }: AboutDialogProps) {
  const [open, setOpen] = useState(false);

  const { data: versionInfo, isLoading } = useQuery({
    queryKey: ["version-info"],
    queryFn: getVersionInfo,
    enabled: open,
  });

  const defaultTrigger = (
    <Button variant="ghost" size="sm" className="w-full justify-start">
      <Info className="mr-2 h-4 w-4" />
      About
    </Button>
  );

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger || defaultTrigger}</DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>About</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <div className="text-sm font-medium">Version</div>
            {isLoading ? (
              <Skeleton className="h-4 w-20" />
            ) : (
              <div className="text-muted-foreground font-mono text-sm">
                {versionInfo?.version}
              </div>
            )}
          </div>
          <div className="space-y-2">
            <div className="text-sm font-medium">Commit</div>
            {isLoading ? (
              <Skeleton className="h-4 w-32" />
            ) : (
              <div className="text-muted-foreground font-mono text-sm">
                {versionInfo?.commit}
              </div>
            )}
          </div>
          <div className="space-y-2">
            <div className="text-sm font-medium">Built</div>
            {isLoading ? (
              <Skeleton className="h-4 w-40" />
            ) : (
              <div className="text-muted-foreground font-mono text-sm">
                {versionInfo?.dateBuilt
                  ? new Date(versionInfo.dateBuilt).toLocaleString()
                  : "Unknown"}
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
