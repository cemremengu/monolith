import { describe, it, expect, vi, beforeEach } from "vitest";
import { screen, waitFor, act } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";

import { render } from "@/test/test-utils";
import { useAuth } from "@/hooks/use-auth";

import { LanguageSwitcher } from "./language-switcher";

describe("LanguageSwitcher", () => {
  beforeEach(() => {
    // Reset auth store before each test
    act(() => {
      useAuth.setState({
        user: null,
        isLoggedIn: false,
      });
    });
  });

  it("should render current language", () => {
    render(<LanguageSwitcher />);

    expect(screen.getByRole("button")).toBeInTheDocument();
    expect(screen.getByText(/english/i)).toBeInTheDocument();
  });

  it("should render with provided value", () => {
    render(<LanguageSwitcher value="tr-TR" />);

    expect(screen.getByText(/türkçe/i)).toBeInTheDocument();
  });

  it("should open dropdown when clicked", async () => {
    const user = userEvent.setup();
    render(<LanguageSwitcher />);

    const trigger = screen.getByRole("button");
    await user.click(trigger);

    await waitFor(() => {
      expect(screen.getByRole("menu")).toBeInTheDocument();
    });
  });

  it("should show available languages in dropdown", async () => {
    const user = userEvent.setup();
    render(<LanguageSwitcher />);

    const trigger = screen.getByRole("button");
    await user.click(trigger);

    await waitFor(() => {
      expect(
        screen.getByRole("menuitem", { name: /english/i }),
      ).toBeInTheDocument();
      expect(
        screen.getByRole("menuitem", { name: /türkçe/i }),
      ).toBeInTheDocument();
    });
  });

  it("should call onChange when language is selected", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<LanguageSwitcher onChange={onChange} />);

    const trigger = screen.getByRole("button");
    await user.click(trigger);

    await waitFor(() => {
      expect(
        screen.getByRole("menuitem", { name: /türkçe/i }),
      ).toBeInTheDocument();
    });

    await user.click(screen.getByRole("menuitem", { name: /türkçe/i }));

    expect(onChange).toHaveBeenCalledWith("tr-TR");
  });

  it("should display flag emoji for current language", () => {
    render(<LanguageSwitcher value="en-US" />);

    // The English flag emoji should be present
    expect(screen.getByText("🇺🇸")).toBeInTheDocument();
  });

  it("should display Turkish flag when Turkish is selected", () => {
    render(<LanguageSwitcher value="tr-TR" />);

    expect(screen.getByText("🇹🇷")).toBeInTheDocument();
  });
});
