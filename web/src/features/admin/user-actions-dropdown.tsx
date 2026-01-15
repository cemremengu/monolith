import { Ban, CheckCircle, MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import { useTranslation } from "react-i18next";

import type { User } from "@/types/api";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

type UserActionsDropdownProps = {
  user: User;
  onEdit: (user: User) => void;
  onDelete: (id: string) => void;
  onDisable: (id: string) => void;
  onEnable: (id: string) => void;
  isDeleting: boolean;
  isDisabling: boolean;
  isEnabling: boolean;
};

export function UserActionsDropdown({
  user,
  onEdit,
  onDelete,
  onDisable,
  onEnable,
  isDeleting,
  isDisabling,
  isEnabling,
}: UserActionsDropdownProps) {
  const { t } = useTranslation();
  const isActive = user.status === "active";
  const isPending = user.status === "pending";

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
          <MoreHorizontal className="h-4 w-4" />
          <span className="sr-only">{t("admin.users.actions.openMenu")}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {isPending ? (
          <DropdownMenuItem
            variant="destructive"
            onClick={() => onDelete(user.id)}
            disabled={isDeleting}
          >
            <Trash2 className="mr-2 h-4 w-4" />
            {t("admin.users.actions.delete")}
          </DropdownMenuItem>
        ) : (
          <>
            <DropdownMenuItem onClick={() => onEdit(user)}>
              <Pencil className="mr-2 h-4 w-4" />
              {t("admin.users.actions.edit")}
            </DropdownMenuItem>
            {isActive ? (
              <DropdownMenuItem
                onClick={() => onDisable(user.id)}
                disabled={isDisabling}
              >
                <Ban className="mr-2 h-4 w-4" />
                {t("admin.users.actions.disable")}
              </DropdownMenuItem>
            ) : (
              <DropdownMenuItem
                onClick={() => onEnable(user.id)}
                disabled={isEnabling}
              >
                <CheckCircle className="mr-2 h-4 w-4" />
                {t("admin.users.actions.enable")}
              </DropdownMenuItem>
            )}
            <DropdownMenuSeparator />
            <DropdownMenuItem
              variant="destructive"
              onClick={() => onDelete(user.id)}
              disabled={isDeleting}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              {t("admin.users.actions.delete")}
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
