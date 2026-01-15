import { createFileRoute, redirect } from "@tanstack/react-router";

import { AdminPanel } from "@/features/admin/admin-panel";

export const Route = createFileRoute("/_authenticated/admin")({
  beforeLoad: ({ context }) => {
    if (!context.auth.user?.isAdmin) {
      throw redirect({
        to: "/",
      });
    }
  },
  component: AdminPanel,
});
