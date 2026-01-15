import { useTranslation } from "react-i18next";
import { Users, UserCheck, UserX } from "lucide-react";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { UsersPage } from "./users-page";
import { useUsers } from "./api/queries";

export function AdminPanel() {
  const { t } = useTranslation();
  const { data: users } = useUsers({});

  const totalUsers = users?.length || 0;
  const activeUsers =
    users?.filter((user) => user.status === "active").length || 0;
  const disabledUsers =
    users?.filter((user) => user.status !== "active").length || 0;

  const stats = [
    {
      title: t("admin.stats.totalUsers"),
      value: totalUsers,
      description: t("admin.stats.totalUsersDesc"),
      icon: Users,
      color: "text-blue-600",
      bgColor: "bg-blue-50 dark:bg-blue-950/20",
    },
    {
      title: t("admin.stats.activeUsers"),
      value: activeUsers,
      description: t("admin.stats.activeUsersDesc"),
      icon: UserCheck,
      color: "text-green-600",
      bgColor: "bg-green-50 dark:bg-green-950/20",
    },
    {
      title: t("admin.stats.disabledUsers"),
      value: disabledUsers,
      description: t("admin.stats.disabledUsersDesc"),
      icon: UserX,
      color: "text-red-600",
      bgColor: "bg-red-50 dark:bg-red-950/20",
    },
  ];

  return (
    <div className="bg-muted/30 min-h-screen">
      <div className="bg-background border-b">
        <div className="mx-auto max-w-7xl px-6 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-4xl font-bold tracking-tight">
                {t("admin.title")}
              </h1>
              <p className="text-muted-foreground mt-2">
                {t("admin.subtitle")}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className="mx-auto max-w-7xl px-6 py-8">
        <div className="mb-8 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {stats.map((stat) => {
            const Icon = stat.icon;
            return (
              <Card key={stat.title} className="overflow-hidden">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    {stat.title}
                  </CardTitle>
                  <div className={`rounded-full p-2 ${stat.bgColor}`}>
                    <Icon className={`h-4 w-4 ${stat.color}`} />
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold">{stat.value}</div>
                  <p className="text-muted-foreground mt-1 text-xs">
                    {stat.description}
                  </p>
                </CardContent>
              </Card>
            );
          })}
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
