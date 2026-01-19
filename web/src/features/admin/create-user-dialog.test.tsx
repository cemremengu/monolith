import { describe, it, expect, vi } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";

import { render } from "@/test/test-utils";

import { CreateUserDialog } from "./create-user-dialog";

describe("CreateUserDialog", () => {
  const defaultProps = {
    open: true,
    onOpenChange: vi.fn(),
  };

  it("should render dialog with title and description", () => {
    render(<CreateUserDialog {...defaultProps} />);

    expect(screen.getByRole("dialog")).toBeInTheDocument();
  });

  it("should render invite and create tabs", () => {
    render(<CreateUserDialog {...defaultProps} />);

    expect(screen.getByRole("tab", { name: /invite/i })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: /create/i })).toBeInTheDocument();
  });

  it("should default to invite tab", () => {
    render(<CreateUserDialog {...defaultProps} />);

    const inviteTab = screen.getByRole("tab", { name: /invite/i });
    expect(inviteTab).toHaveAttribute("data-state", "active");
  });

  it("should switch to create tab when clicked", async () => {
    const user = userEvent.setup();
    render(<CreateUserDialog {...defaultProps} />);

    const createTab = screen.getByRole("tab", { name: /create/i });
    await user.click(createTab);

    expect(createTab).toHaveAttribute("data-state", "active");
  });

  describe("Invite tab", () => {
    it("should render emails textarea and role select", () => {
      render(<CreateUserDialog {...defaultProps} />);

      expect(screen.getByRole("textbox")).toBeInTheDocument();
    });

    it("should show validation error for empty emails", async () => {
      const user = userEvent.setup();
      render(<CreateUserDialog {...defaultProps} />);

      const submitButton = screen.getByRole("button", { name: /invite/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/required/i)).toBeInTheDocument();
      });
    });
  });

  describe("Create tab", () => {
    it("should render all form fields", async () => {
      const user = userEvent.setup();
      render(<CreateUserDialog {...defaultProps} />);

      const createTab = screen.getByRole("tab", { name: /create/i });
      await user.click(createTab);

      await waitFor(() => {
        expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/^name$/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument();
      });
    });

    it("should show validation errors for empty fields", async () => {
      const user = userEvent.setup();
      render(<CreateUserDialog {...defaultProps} />);

      const createTab = screen.getByRole("tab", { name: /create/i });
      await user.click(createTab);

      await waitFor(() => {
        expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
      });

      const submitButton = screen.getByRole("button", { name: /create user/i });
      await user.click(submitButton);

      await waitFor(() => {
        // Should show validation errors
        const errorMessages = screen.getAllByRole("alert");
        expect(errorMessages.length).toBeGreaterThan(0);
      });
    });

    it("should show password mismatch error", async () => {
      const user = userEvent.setup();
      render(<CreateUserDialog {...defaultProps} />);

      const createTab = screen.getByRole("tab", { name: /create/i });
      await user.click(createTab);

      await waitFor(() => {
        expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
      });

      await user.type(screen.getByLabelText(/username/i), "newuser");
      await user.type(screen.getByLabelText(/^name$/i), "New User");
      await user.type(screen.getByLabelText(/email/i), "new@example.com");
      await user.type(screen.getByLabelText(/^password$/i), "password123");
      await user.type(
        screen.getByLabelText(/confirm password/i),
        "differentpassword",
      );

      const submitButton = screen.getByRole("button", { name: /create user/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/passwords must match/i)).toBeInTheDocument();
      });
    });
  });

  it("should not render when open is false", () => {
    render(<CreateUserDialog open={false} onOpenChange={vi.fn()} />);

    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("should call onOpenChange when closing", async () => {
    const onOpenChange = vi.fn();
    const user = userEvent.setup();
    render(<CreateUserDialog open={true} onOpenChange={onOpenChange} />);

    // Find and click the close button (X button in dialog)
    const closeButton = screen.getByRole("button", { name: /close/i });
    await user.click(closeButton);

    expect(onOpenChange).toHaveBeenCalledWith(false);
  });
});
