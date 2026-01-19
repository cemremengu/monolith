import { describe, it, expect, beforeEach } from "vitest";
import { waitFor } from "@testing-library/react";

import { mockUser } from "@/test/mocks/handlers";

import { useAuth } from "./use-auth";

describe("useAuth", () => {
  beforeEach(() => {
    // Clear localStorage first
    localStorage.clear();
    // Reset the store before each test
    useAuth.setState({
      user: null,
      isLoggedIn: false,
    });
  });

  describe("login", () => {
    it("should set isLoggedIn to true on successful login", async () => {
      const store = useAuth.getState();
      expect(store.isLoggedIn).toBe(false);

      await store.login({
        login: "testuser",
        password: "password123",
      });

      expect(useAuth.getState().isLoggedIn).toBe(true);
    });

    it("should throw error on failed login", async () => {
      const store = useAuth.getState();

      await expect(
        store.login({
          login: "wronguser",
          password: "wrongpassword",
        }),
      ).rejects.toThrow();

      expect(useAuth.getState().isLoggedIn).toBe(false);
    });
  });

  describe("logout", () => {
    it("should clear isLoggedIn and user state", async () => {
      // First login to set the state
      const store = useAuth.getState();

      await store.login({
        login: "testuser",
        password: "password123",
      });

      expect(useAuth.getState().isLoggedIn).toBe(true);

      // Then logout
      await useAuth.getState().logout();

      await waitFor(() => {
        expect(useAuth.getState().isLoggedIn).toBe(false);
        expect(useAuth.getState().user).toBeNull();
      });
    });

    it("should call window.location.replace with /login", async () => {
      // First login
      const store = useAuth.getState();

      await store.login({
        login: "testuser",
        password: "password123",
      });

      // Then logout
      await useAuth.getState().logout();

      await waitFor(() => {
        expect(window.location.replace).toHaveBeenCalledWith("/login");
      });
    });
  });

  describe("fetchUser", () => {
    it("should populate user on successful fetch", async () => {
      const store = useAuth.getState();

      await store.fetchUser();

      await waitFor(() => {
        expect(useAuth.getState().user).toEqual(mockUser);
      });
    });

    it("should keep user state after successful fetch when logged in", async () => {
      const store = useAuth.getState();

      // Login first
      await store.login({
        login: "testuser",
        password: "password123",
      });

      expect(useAuth.getState().isLoggedIn).toBe(true);

      await useAuth.getState().fetchUser();

      await waitFor(() => {
        expect(useAuth.getState().user).toEqual(mockUser);
      });
    });
  });
});
