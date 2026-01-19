import { describe, it, expect, vi } from "vitest";
import { screen } from "@testing-library/react";

import { render } from "@/test/test-utils";
import { mockUser, mockAdminUser } from "@/test/mocks/handlers";

import { UserTable } from "./user-table";

describe("UserTable", () => {
  const defaultProps = {
    users: [],
    onEdit: vi.fn(),
    onDelete: vi.fn(),
    onDisable: vi.fn(),
    onEnable: vi.fn(),
    isDeleting: false,
    isDisabling: false,
    isEnabling: false,
  };

  it("should render empty state when no users", () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByText(/no users/i)).toBeInTheDocument();
  });

  it("should render user rows when users are provided", () => {
    render(<UserTable {...defaultProps} users={[mockUser, mockAdminUser]} />);

    expect(screen.getByText("testuser")).toBeInTheDocument();
    expect(screen.getByText("test@example.com")).toBeInTheDocument();
    expect(screen.getByText("admin")).toBeInTheDocument();
    expect(screen.getByText("admin@example.com")).toBeInTheDocument();
  });

  it("should render role badges correctly", () => {
    render(<UserTable {...defaultProps} users={[mockUser, mockAdminUser]} />);

    const userBadges = screen.getAllByText("User");
    const adminBadges = screen.getAllByText("Admin");

    expect(userBadges.length).toBeGreaterThanOrEqual(1);
    expect(adminBadges.length).toBeGreaterThanOrEqual(1);
  });

  it("should render status badges correctly", () => {
    const pendingUser = {
      ...mockUser,
      id: "3",
      status: "pending",
    };

    render(<UserTable {...defaultProps} users={[mockUser, pendingUser]} />);

    // The status is translated, so we check for the translation key result
    // In tests, i18n returns the key or translated value
    expect(
      screen.getAllByText(/active|pending/i).length,
    ).toBeGreaterThanOrEqual(1);
  });

  it("should render table headers", () => {
    render(<UserTable {...defaultProps} users={[mockUser]} />);

    // Check for column headers (translation keys may be returned)
    expect(screen.getByRole("table")).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: /username/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: /email/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: /role/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("columnheader", { name: /status/i }),
    ).toBeInTheDocument();
  });

  it("should render user creation date", () => {
    render(<UserTable {...defaultProps} users={[mockUser]} />);

    // The date is formatted by toLocaleDateString
    const dateCell = screen.getByText(/1\/1\/2024|2024/);
    expect(dateCell).toBeInTheDocument();
  });

  it("should render actions dropdown for each user", () => {
    render(<UserTable {...defaultProps} users={[mockUser, mockAdminUser]} />);

    // Should have actions buttons (the MoreHorizontal icon button)
    const actionButtons = screen.getAllByRole("button");
    expect(actionButtons.length).toBeGreaterThanOrEqual(2);
  });
});
