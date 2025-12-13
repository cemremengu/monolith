import { createFileRoute } from "@tanstack/react-router";

import { profileQueryOptions } from "@/features/profile/api/queries";
import { SettingsPage } from "@/features/profile/settings-page";

export const Route = createFileRoute("/_authenticated/settings")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(profileQueryOptions),
  component: SettingsPage,
});
