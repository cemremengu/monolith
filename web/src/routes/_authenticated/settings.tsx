import { createFileRoute, redirect } from "@tanstack/react-router";

import { profileQueryOptions } from "@/features/profile/api/queries";
import { SettingsPage, type SettingsSection } from "@/features/profile/settings-page";

type SettingsSearch = { section: SettingsSection };

const adminSections: SettingsSection[] = ["users", "providers", "license"];

function toSection(value: unknown): SettingsSection {
  if (value === "users" || value === "providers" || value === "license") return value;
  return "preferences";
}

export const Route = createFileRoute("/_authenticated/settings")({
  validateSearch: (search: Record<string, unknown>): SettingsSearch => ({
    section: toSection(search.section),
  }),
  beforeLoad: ({ context, search }) => {
    if (adminSections.includes(search.section) && !context.auth.user?.isAdmin) {
      throw redirect({ to: "/settings", search: { section: "preferences" } });
    }
  },
  loader: ({ context: { queryClient } }) => queryClient.ensureQueryData(profileQueryOptions),
  component: RouteComponent,
});

function RouteComponent() {
  const { section } = Route.useSearch();
  return <SettingsPage section={section} />;
}
