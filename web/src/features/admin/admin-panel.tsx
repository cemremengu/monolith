import { useTranslation } from "react-i18next";
import { Users } from "lucide-react";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { UsersPage } from "./users-page";

export function AdminPanel() {
  const { t } = useTranslation();

  return (
    <div className="px-6 py-3.5">
      <div className="mx-auto max-w-7xl space-y-6">
        <div>
          <h1 className="text-3xl font-bold">{t("admin.title")}</h1>
          <p className="text-muted-foreground">{t("admin.subtitle")}</p>
        </div>

        <Tabs defaultValue="users" className="space-y-6">
          <TabsList className="grid w-full max-w-md grid-cols-1">
            <TabsTrigger value="users" className="flex items-center gap-2">
              <Users className="h-4 w-4" />
              {t("admin.users.title")}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="users" className="space-y-6">
            <UsersPage />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
