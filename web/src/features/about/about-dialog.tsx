import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";

import { getVersionInfo } from "./api";

type AboutDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function AboutDialog({ open, onOpenChange }: AboutDialogProps) {
  const { t } = useTranslation();
  const { data: versionInfo, isLoading } = useQuery({
    queryKey: ["version-info"],
    queryFn: getVersionInfo,
    enabled: open,
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("about.title")}</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <div className="text-sm font-medium">{t("about.version")}</div>
            {isLoading ? (
              <Skeleton className="h-4 w-20" />
            ) : (
              <div className="font-mono text-sm text-muted-foreground">{versionInfo?.version}</div>
            )}
          </div>
          <div className="space-y-2">
            <div className="text-sm font-medium">{t("about.commit")}</div>
            {isLoading ? (
              <Skeleton className="h-4 w-32" />
            ) : (
              <div className="font-mono text-sm text-muted-foreground">{versionInfo?.commit}</div>
            )}
          </div>
          <div className="space-y-2">
            <div className="text-sm font-medium">{t("about.built")}</div>
            {isLoading ? (
              <Skeleton className="h-4 w-40" />
            ) : (
              <div className="font-mono text-sm text-muted-foreground">
                {versionInfo?.buildTime
                  ? new Date(versionInfo.buildTime).toLocaleString()
                  : t("about.unknown")}
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
