import { useTranslation } from "react-i18next";

import type { User } from "@/types/api";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

import { UserActionsDropdown } from "./user-actions-dropdown";

type UserTableProps = {
  users: User[];
  onEdit: (user: User) => void;
  onDelete: (id: string) => void;
  onDisable: (id: string) => void;
  onEnable: (id: string) => void;
  isDeleting: boolean;
  isDisabling: boolean;
  isEnabling: boolean;
};

export function UserTable({
  users,
  onEdit,
  onDelete,
  onDisable,
  onEnable,
  isDeleting,
  isDisabling,
  isEnabling,
}: UserTableProps) {
  const { t } = useTranslation();

  const getStatusVariant = (status: string) => {
    switch (status) {
      case "active":
        return "default";
      case "pending":
        return "secondary";
      default:
        return "outline";
    }
  };

  const getRoleBadge = (isAdmin: boolean) => {
    return isAdmin ? "Admin" : "User";
  };

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t("admin.users.table.username")}</TableHead>
            <TableHead>{t("admin.users.table.email")}</TableHead>
            <TableHead>{t("admin.users.table.role")}</TableHead>
            <TableHead>{t("admin.users.table.status")}</TableHead>
            <TableHead>{t("admin.users.table.createdAt")}</TableHead>
            <TableHead className="w-12.5"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {users.length === 0 ? (
            <TableRow>
              <TableCell colSpan={6} className="h-24 text-center">
                {t("admin.users.messages.noUsers")}
              </TableCell>
            </TableRow>
          ) : (
            users.map((user) => (
              <TableRow key={user.id}>
                <TableCell className="font-medium">{user.username}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  <Badge variant="outline">{getRoleBadge(user.isAdmin)}</Badge>
                </TableCell>
                <TableCell>
                  <Badge variant={getStatusVariant(user.status)}>
                    {t(`admin.users.status.${user.status}`)}
                  </Badge>
                </TableCell>
                <TableCell>
                  {new Date(user.createdAt).toLocaleDateString()}
                </TableCell>
                <TableCell>
                  <UserActionsDropdown
                    user={user}
                    onEdit={onEdit}
                    onDelete={onDelete}
                    onDisable={onDisable}
                    onEnable={onEnable}
                    isDeleting={isDeleting}
                    isDisabling={isDisabling}
                    isEnabling={isEnabling}
                  />
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}
