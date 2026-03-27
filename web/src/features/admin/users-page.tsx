import { useState } from "react";
import { useTranslation } from "react-i18next";
import type { VariantProps } from "class-variance-authority";
import { PlusIcon } from "lucide-react";

import type { User } from "@/types/api";
import { DataTable, type ColumnDef } from "@/components/datatable";
import { Badge, badgeVariants } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

import {
  useDeleteUser,
  useDisableUser,
  useEnableUser,
  useUsers,
} from "./api/queries";
import { CreateUserDialog } from "./create-user-dialog";
import { UserActionsDropdown } from "./user-actions-dropdown";

export function UsersPage() {
  const { t } = useTranslation();
  const { data: users } = useUsers({});
  const deleteUser = useDeleteUser();
  const disableUser = useDisableUser();
  const enableUser = useEnableUser();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  const handleEdit = () => {
    // TODO: Implement edit functionality
  };

  const handleDelete = (id: string) => {
    deleteUser.mutate(id);
  };

  const handleDisable = (id: string) => {
    disableUser.mutate(id);
  };

  const handleEnable = (id: string) => {
    enableUser.mutate(id);
  };

  const getStatusVariant = (
    status: string,
  ): NonNullable<VariantProps<typeof badgeVariants>["variant"]> => {
    switch (status) {
      case "active":
        return "default";
      case "pending":
        return "secondary";
      default:
        return "outline";
    }
  };

  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "username",
      header: t("admin.users.table.username"),
      cell: ({ row }) => (
        <span className="font-medium">{row.original.username}</span>
      ),
    },
    {
      accessorKey: "email",
      header: t("admin.users.table.email"),
    },
    {
      id: "role",
      accessorFn: (user) =>
        user.isAdmin
          ? t("admin.users.roles.admin")
          : t("admin.users.roles.user"),
      header: t("admin.users.table.role"),
      cell: ({ row }) => (
        <Badge variant="outline">
          {row.original.isAdmin
            ? t("admin.users.roles.admin")
            : t("admin.users.roles.user")}
        </Badge>
      ),
    },
    {
      accessorKey: "status",
      header: t("admin.users.table.status"),
      cell: ({ row }) => (
        <Badge variant={getStatusVariant(row.original.status)}>
          {t(`admin.users.status.${row.original.status}`)}
        </Badge>
      ),
    },
    {
      accessorKey: "createdAt",
      header: t("admin.users.table.createdAt"),
      cell: ({ row }) => new Date(row.original.createdAt).toLocaleDateString(),
    },
    {
      id: "actions",
      cell: ({ row }) => (
        <div className="flex justify-end">
          <UserActionsDropdown
            user={row.original}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onDisable={handleDisable}
            onEnable={handleEnable}
            isDeleting={deleteUser.isPending}
            isDisabling={disableUser.isPending}
            isEnabling={enableUser.isPending}
          />
        </div>
      ),
      enableHiding: false,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-muted-foreground text-sm">
            {t("admin.users.description")}
          </p>
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <PlusIcon className="h-4 w-4" />
          {t("admin.users.newUser")}
        </Button>
      </div>

      <DataTable
        data={users}
        columns={columns}
        emptyText={t("admin.users.messages.noUsers")}
      />

      <CreateUserDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
      />
    </div>
  );
}
