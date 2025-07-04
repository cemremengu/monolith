import { createFileRoute } from "@tanstack/react-router";
import { profileQueryOptions } from "@/features/profile/api/queries";
import { ProfilePage } from "@/features/profile/components/profile-page";

export const Route = createFileRoute("/_authenticated/profile")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(profileQueryOptions),
  component: ProfilePage,
});
