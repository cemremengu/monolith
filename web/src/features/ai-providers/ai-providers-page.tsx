import { Boxes, Bot, Brain, Info, Plus, Sparkles, type LucideIcon } from "lucide-react";
import { useTranslation } from "react-i18next";

import { PageHeader } from "@/components/page-header";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";

type DemoProvider = {
  name: string;
  typeLabel: string;
  models: number;
  isDefault: boolean;
  icon: LucideIcon;
  accent: string;
};

const configuredProviders: DemoProvider[] = [
  {
    name: "OpenAI",
    typeLabel: "OpenAI",
    models: 4,
    isDefault: true,
    icon: Sparkles,
    accent: "bg-foreground/5 text-foreground",
  },
  {
    name: "Anthropic",
    typeLabel: "Anthropic",
    models: 3,
    isDefault: false,
    icon: Brain,
    accent: "bg-orange-500/10 text-orange-600 dark:text-orange-400",
  },
];

const availableProviders: DemoProvider[] = [
  {
    name: "DeepSeek",
    typeLabel: "DeepSeek",
    models: 2,
    isDefault: false,
    icon: Bot,
    accent: "bg-blue-500/10 text-blue-600 dark:text-blue-400",
  },
  {
    name: "Custom",
    typeLabel: "Custom",
    models: 0,
    isDefault: false,
    icon: Boxes,
    accent: "bg-amber-500/10 text-amber-600 dark:text-amber-400",
  },
];

function ProviderIcon({ icon: Icon, accent }: { icon: LucideIcon; accent: string }) {
  return (
    <div className={cn("flex size-8 shrink-0 items-center justify-center rounded-lg", accent)}>
      <Icon className="size-4" />
    </div>
  );
}

export function AIProvidersPage() {
  const { t } = useTranslation();

  return (
    <div className="space-y-8">
      <PageHeader
        title={t("aiProviders.title")}
        description={t("aiProviders.description")}
        actions={
          <Button disabled>
            <Plus className="size-4" />
            {t("aiProviders.addProvider")}
          </Button>
        }
      />

      <div className="flex items-start gap-2 rounded-lg border border-dashed p-3 text-sm text-muted-foreground">
        <Info className="mt-0.5 size-4 shrink-0" />
        <span>{t("common.previewNotice")}</span>
      </div>

      <Card className="gap-0 p-5">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="space-y-1">
            <h2 className="text-sm font-medium">{t("aiProviders.activeModel.title")}</h2>
            <p className="text-sm text-muted-foreground">
              {t("aiProviders.activeModel.current", { model: "gpt-4o", provider: "OpenAI" })}
            </p>
          </div>
          <Button variant="outline" disabled>
            {t("aiProviders.activeModel.placeholder")}
          </Button>
        </div>
      </Card>

      <section className="space-y-3">
        <div className="space-y-0.5">
          <h2 className="text-sm font-medium">{t("aiProviders.configuredTitle")}</h2>
          <p className="text-sm text-muted-foreground">{t("aiProviders.configuredSubtitle")}</p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          {configuredProviders.map((provider) => (
            <Card key={provider.name} className="gap-0 p-4">
              <div className="flex min-w-0 items-center gap-3">
                <ProviderIcon icon={provider.icon} accent={provider.accent} />
                <div className="min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="truncate text-sm font-medium">{provider.name}</p>
                    {provider.isDefault && (
                      <Badge variant="outline">{t("aiProviders.defaultBadge")}</Badge>
                    )}
                  </div>
                  <p className="truncate text-xs text-muted-foreground">
                    {provider.typeLabel} · {t("aiProviders.modelCount", { count: provider.models })}
                  </p>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </section>

      <section className="space-y-3">
        <div className="space-y-0.5">
          <h2 className="text-sm font-medium">{t("aiProviders.availableTitle")}</h2>
          <p className="text-sm text-muted-foreground">{t("aiProviders.availableSubtitle")}</p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          {availableProviders.map((provider) => (
            <Card key={provider.name} className="gap-0 p-4 opacity-60">
              <div className="flex items-center gap-3">
                <ProviderIcon icon={provider.icon} accent={provider.accent} />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium">{provider.name}</p>
                  <p className="truncate text-xs text-muted-foreground">
                    {provider.models > 0
                      ? t("aiProviders.modelCount", { count: provider.models })
                      : t("aiProviders.customSubtitle")}
                  </p>
                </div>
                <Plus className="size-4 shrink-0 text-muted-foreground" />
              </div>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
}
