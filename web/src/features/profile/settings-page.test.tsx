import { screen, waitFor } from "@testing-library/react";
import { Suspense } from "react";
import { describe, it, expect } from "vitest";

import { render } from "@/test/test-utils";

import { SettingsPage } from "./settings-page";

function renderWithSuspense(ui: React.ReactElement) {
  return render(<Suspense fallback={<div>Loading...</div>}>{ui}</Suspense>);
}

describe("SettingsPage", () => {
  it("should render the preferences heading by default", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByRole("heading", { level: 1 })).toBeInTheDocument();
    });
  });

  it("should render a row for each preference", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Theme")).toBeInTheDocument();
      expect(screen.getByText("Language")).toBeInTheDocument();
      expect(screen.getByText("Timezone")).toBeInTheDocument();
    });
  });

  it("should render the theme options", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByRole("button", { name: "system" })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "light" })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "dark" })).toBeInTheDocument();
    });
  });

  it("should render the current language", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("English")).toBeInTheDocument();
    });
  });
});
