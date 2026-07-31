import { CheckCircle2, Info } from "lucide-react";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";

const demoLicense = {
  licensee: "Acme Inc.",
  projectLimit: 25,
  issuedAt: "2026-01-15",
  expiresAt: "2027-01-15",
  used: 8,
};

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="space-y-0.5">
      <dt className="text-xs text-muted-foreground">{label}</dt>
      <dd className="text-sm font-medium">{value}</dd>
    </div>
  );
}

export function LicensePage() {
  const { t } = useTranslation();
  const percent = Math.round((demoLicense.used / demoLicense.projectLimit) * 100);

  return (
    <div className="space-y-8">
      <PageHeader title={t("license.title")} description={t("license.description")} />

      <div className="flex items-start gap-2 rounded-lg border border-dashed p-3 text-sm text-muted-foreground">
        <Info className="mt-0.5 size-4 shrink-0" />
        <span>{t("common.previewNotice")}</span>
      </div>

      <Card className="gap-0 space-y-5 p-5">
        <div className="flex items-center justify-between gap-3">
          <div className="space-y-1">
            <h2 className="text-sm font-medium">{t("license.status.title")}</h2>
            <p className="text-sm text-muted-foreground">
              {t("license.status.licensee", { name: demoLicense.licensee })}
            </p>
          </div>
          <Badge variant="outline">{t("license.states.valid")}</Badge>
        </div>

        <dl className="grid grid-cols-2 gap-4 sm:grid-cols-3">
          <Detail
            label={t("license.status.projectLimit")}
            value={String(demoLicense.projectLimit)}
          />
          <Detail
            label={t("license.status.issuedAt")}
            value={new Date(demoLicense.issuedAt).toLocaleDateString()}
          />
          <Detail
            label={t("license.status.expiresAt")}
            value={new Date(demoLicense.expiresAt).toLocaleDateString()}
          />
        </dl>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">{t("license.usage.title")}</span>
            <span className="font-medium">
              {t("license.usage.ofLimit", {
                used: demoLicense.used,
                limit: demoLicense.projectLimit,
              })}
            </span>
          </div>
          <div className="h-2 overflow-hidden rounded-full bg-muted">
            <div className="h-full rounded-full bg-primary" style={{ width: `${percent}%` }} />
          </div>
        </div>
      </Card>

      <section className="space-y-3">
        <div className="space-y-0.5">
          <h2 className="text-sm font-medium">{t("license.install.title")}</h2>
          <p className="text-sm text-muted-foreground">{t("license.install.subtitle")}</p>
        </div>

        <Textarea
          placeholder={t("license.install.placeholder")}
          rows={4}
          spellCheck={false}
          className="font-mono text-xs"
          disabled
        />

        <div className="flex justify-end">
          <Button disabled>
            <CheckCircle2 className="size-4" />
            {t("license.install.action")}
          </Button>
        </div>
      </section>
    </div>
  );
}
