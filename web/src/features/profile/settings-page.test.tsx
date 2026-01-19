import { Suspense } from "react";
import { describe, it, expect } from "vitest";
import { screen, waitFor } from "@testing-library/react";

import { render } from "@/test/test-utils";

import { SettingsPage } from "./settings-page";

function renderWithSuspense(ui: React.ReactElement) {
  return render(<Suspense fallback={<div>Loading...</div>}>{ui}</Suspense>);
}

describe("SettingsPage", () => {
  it("should render settings title", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByRole("heading", { level: 1 })).toBeInTheDocument();
    });
  });

  it("should render language switcher component", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      // The language switcher shows the current language (English by default)
      expect(screen.getByText(/english/i)).toBeInTheDocument();
    });
  });

  it("should render preferences card with sections", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      // Look for the icons that indicate each section (Globe, Palette, Clock)
      // The settings page has a grid with 3 sections
      const buttons = screen.getAllByRole("button");
      // Should have at least the language, theme, and timezone selectors
      expect(buttons.length).toBeGreaterThanOrEqual(3);
    });
  });

  it("should render language flag emoji", async () => {
    renderWithSuspense(<SettingsPage />);

    await waitFor(() => {
      // The English flag emoji should be displayed
      expect(screen.getByText("🇺🇸")).toBeInTheDocument();
    });
  });
});
