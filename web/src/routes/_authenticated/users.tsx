import { createFileRoute } from "@tanstack/react-router";
import { usersQueryOptions } from "@/features/users/api/queries";
import { UsersPage } from "@/features/users/components/users-page";

export const Route = createFileRoute("/_authenticated/users")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(usersQueryOptions({})),
  component: UsersPage,
});
