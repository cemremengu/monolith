import { useState } from "react";
import { useTranslation } from "react-i18next";
import { PlusIcon } from "lucide-react";

import { Button } from "@/components/ui/button";

import {
  useDeleteUser,
  useDisableUser,
  useEnableUser,
  useUsers,
} from "./api/queries";
import { UserTable } from "./user-table";
import { CreateUserDialog } from "./create-user-dialog";

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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-muted-foreground text-sm">
            {t("admin.users.description")}
          </p>
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <PlusIcon className="mr-2 h-4 w-4" />
          {t("admin.users.newUser")}
        </Button>
      </div>

      <UserTable
        users={users}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onDisable={handleDisable}
        onEnable={handleEnable}
        isDeleting={deleteUser.isPending}
        isDisabling={disableUser.isPending}
        isEnabling={enableUser.isPending}
      />

      <CreateUserDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
      />
    </div>
  );
}
