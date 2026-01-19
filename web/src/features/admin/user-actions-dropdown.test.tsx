import { describe, it, expect, vi } from "vitest";
import { userEvent } from "@testing-library/user-event";
import { screen, waitFor } from "@testing-library/react";

import { mockUser } from "@/test/mocks/handlers";
import { render } from "@/test/test-utils";

import { UserActionsDropdown } from "./user-actions-dropdown";

describe("UserActionsDropdown", () => {
  const defaultProps = {
    user: mockUser,
    onEdit: vi.fn(),
    onDelete: vi.fn(),
    onDisable: vi.fn(),
    onEnable: vi.fn(),
    isDeleting: false,
    isDisabling: false,
    isEnabling: false,
  };

  it("should render dropdown trigger button", () => {
    render(<UserActionsDropdown {...defaultProps} />);

    expect(screen.getByRole("button")).toBeInTheDocument();
  });

  it("should open dropdown menu when clicked", async () => {
    const user = userEvent.setup();
    render(<UserActionsDropdown {...defaultProps} />);

    const trigger = screen.getByRole("button");
    await user.click(trigger);

    await waitFor(() => {
      expect(screen.getByRole("menu")).toBeInTheDocument();
    });
  });

  describe("Active user actions", () => {
    it("should show edit, disable, and delete options for active user", async () => {
      const user = userEvent.setup();
      render(<UserActionsDropdown {...defaultProps} />);

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/edit/i)).toBeInTheDocument();
        expect(screen.getByText(/disable/i)).toBeInTheDocument();
        expect(screen.getByText(/delete/i)).toBeInTheDocument();
      });
    });

    it("should call onEdit when edit is clicked", async () => {
      const user = userEvent.setup();
      const onEdit = vi.fn();
      render(<UserActionsDropdown {...defaultProps} onEdit={onEdit} />);

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/edit/i)).toBeInTheDocument();
      });

      await user.click(screen.getByText(/edit/i));

      expect(onEdit).toHaveBeenCalledWith(mockUser);
    });

    it("should call onDisable when disable is clicked", async () => {
      const user = userEvent.setup();
      const onDisable = vi.fn();
      render(<UserActionsDropdown {...defaultProps} onDisable={onDisable} />);

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/disable/i)).toBeInTheDocument();
      });

      await user.click(screen.getByText(/disable/i));

      expect(onDisable).toHaveBeenCalledWith(mockUser.id);
    });

    it("should call onDelete when delete is clicked", async () => {
      const user = userEvent.setup();
      const onDelete = vi.fn();
      render(<UserActionsDropdown {...defaultProps} onDelete={onDelete} />);

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/delete/i)).toBeInTheDocument();
      });

      await user.click(screen.getByText(/delete/i));

      expect(onDelete).toHaveBeenCalledWith(mockUser.id);
    });
  });

  describe("Disabled user actions", () => {
    const disabledUser = { ...mockUser, status: "disabled" };

    it("should show enable option for disabled user", async () => {
      const user = userEvent.setup();
      render(<UserActionsDropdown {...defaultProps} user={disabledUser} />);

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/enable/i)).toBeInTheDocument();
        expect(screen.queryByText(/disable/i)).not.toBeInTheDocument();
      });
    });

    it("should call onEnable when enable is clicked", async () => {
      const user = userEvent.setup();
      const onEnable = vi.fn();
      render(
        <UserActionsDropdown
          {...defaultProps}
          user={disabledUser}
          onEnable={onEnable}
        />,
      );

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/enable/i)).toBeInTheDocument();
      });

      await user.click(screen.getByText(/enable/i));

      expect(onEnable).toHaveBeenCalledWith(disabledUser.id);
    });
  });

  describe("Pending user actions", () => {
    const pendingUser = { ...mockUser, status: "pending" };

    it("should only show delete option for pending user", async () => {
      const user = userEvent.setup();
      render(<UserActionsDropdown {...defaultProps} user={pendingUser} />);

      const trigger = screen.getByRole("button");
      await user.click(trigger);

      await waitFor(() => {
        expect(screen.getByText(/delete/i)).toBeInTheDocument();
        expect(screen.queryByText(/edit/i)).not.toBeInTheDocument();
        expect(screen.queryByText(/disable/i)).not.toBeInTheDocument();
        expect(screen.queryByText(/enable/i)).not.toBeInTheDocument();
      });
    });
  });
});
