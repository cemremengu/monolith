import { Suspense } from "react";

import { Spinner } from "@/components/ui/spinner";
import { UsersPage } from "@/features/admin/users-page";
import { AIProvidersPage } from "@/features/ai-providers/ai-providers-page";
import { LicensePage } from "@/features/license/license-page";

import { PreferencesSection } from "./preferences-section";

export type SettingsSection = "preferences" | "users" | "providers" | "license";

export function SettingsPage({ section = "preferences" }: { section?: SettingsSection }) {
  return (
    <div className="mx-auto max-w-5xl px-8 py-8">
      <Suspense
        fallback={
          <div className="flex justify-center py-8">
            <Spinner className="size-5" />
          </div>
        }
      >
        {section === "users" ? (
          <UsersPage />
        ) : section === "providers" ? (
          <AIProvidersPage />
        ) : section === "license" ? (
          <LicensePage />
        ) : (
          <PreferencesSection />
        )}
      </Suspense>
    </div>
  );
}
