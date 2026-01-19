import { describe, it, expect, vi, beforeEach } from "vitest";
import { userEvent } from "@testing-library/user-event";
import { screen, waitFor, act } from "@testing-library/react";

import { render } from "@/test/test-utils";
import { useAuth } from "@/hooks/use-auth";

import { LoginForm } from "./login-form";

describe("LoginForm", () => {
  const mockOnSuccess = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    // Reset auth store
    act(() => {
      useAuth.setState({
        user: null,
        isLoggedIn: false,
      });
    });
  });

  it("should render login form with username and password fields", () => {
    render(<LoginForm onSuccess={mockOnSuccess} />);

    expect(screen.getByLabelText(/username or email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /login/i })).toBeInTheDocument();
  });

  it("should show validation error when username is empty", async () => {
    const user = userEvent.setup();
    render(<LoginForm onSuccess={mockOnSuccess} />);

    const submitButton = screen.getByRole("button", { name: /login/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(
        screen.getByText(/username or email is required/i),
      ).toBeInTheDocument();
    });
  });

  it("should show validation error when password is empty", async () => {
    const user = userEvent.setup();
    render(<LoginForm onSuccess={mockOnSuccess} />);

    const loginInput = screen.getByLabelText(/username or email/i);
    await user.type(loginInput, "testuser");

    const submitButton = screen.getByRole("button", { name: /login/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/password is required/i)).toBeInTheDocument();
    });
  });

  it("should call onSuccess on successful login", async () => {
    const user = userEvent.setup();
    render(<LoginForm onSuccess={mockOnSuccess} />);

    const loginInput = screen.getByLabelText(/username or email/i);
    const passwordInput = screen.getByLabelText(/password/i);

    await user.type(loginInput, "testuser");
    await user.type(passwordInput, "password123");

    const submitButton = screen.getByRole("button", { name: /login/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalled();
    });
  });

  it("should show error message on failed login", async () => {
    const user = userEvent.setup();
    render(<LoginForm onSuccess={mockOnSuccess} />);

    const loginInput = screen.getByLabelText(/username or email/i);
    const passwordInput = screen.getByLabelText(/password/i);

    await user.type(loginInput, "wronguser");
    await user.type(passwordInput, "wrongpassword");

    const submitButton = screen.getByRole("button", { name: /login/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(
        screen.getByText(/invalid username\/email or password/i),
      ).toBeInTheDocument();
    });

    expect(mockOnSuccess).not.toHaveBeenCalled();
  });

  it("should disable button while submitting", async () => {
    const user = userEvent.setup();
    render(<LoginForm onSuccess={mockOnSuccess} />);

    const loginInput = screen.getByLabelText(/username or email/i);
    const passwordInput = screen.getByLabelText(/password/i);

    await user.type(loginInput, "testuser");
    await user.type(passwordInput, "password123");

    const submitButton = screen.getByRole("button", { name: /login/i });

    // Before clicking, button should not be disabled
    expect(submitButton).not.toBeDisabled();

    await user.click(submitButton);

    // Wait for completion
    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalled();
    });

    // After completion, button should be enabled again
    expect(submitButton).not.toBeDisabled();
  });
});
