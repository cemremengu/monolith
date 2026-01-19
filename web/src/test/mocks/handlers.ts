import { http, HttpResponse } from "msw";

import type { User } from "@/types/api";

const BASE_URL = "http://localhost:3000";

export const mockUser: User = {
  id: "1",
  username: "testuser",
  email: "test@example.com",
  name: "Test User",
  isAdmin: false,
  language: "en-US",
  theme: "system",
  timezone: "UTC",
  status: "active",
  isDisabled: false,
  createdAt: "2024-01-01T00:00:00Z",
  updatedAt: "2024-01-01T00:00:00Z",
};

export const mockAdminUser: User = {
  ...mockUser,
  id: "2",
  username: "admin",
  email: "admin@example.com",
  name: "Admin User",
  isAdmin: true,
};

export const mockUsers: User[] = [mockUser, mockAdminUser];

export const handlers = [
  // Auth login
  http.post(`${BASE_URL}/api/login`, async ({ request }) => {
    const body = (await request.json()) as { login: string; password: string };

    if (body.login === "testuser" && body.password === "password123") {
      return HttpResponse.json({ message: "Login successful" });
    }

    return HttpResponse.json({ error: "Invalid credentials" }, { status: 401 });
  }),

  // Auth logout
  http.post(`${BASE_URL}/api/logout`, () => {
    return HttpResponse.json({ message: "Logged out" });
  }),

  // Account profile
  http.get(`${BASE_URL}/api/account/profile`, () => {
    return HttpResponse.json(mockUser);
  }),

  // Update preferences
  http.patch(`${BASE_URL}/api/account/preferences`, async ({ request }) => {
    const body = (await request.json()) as Partial<User>;
    return HttpResponse.json({ ...mockUser, ...body });
  }),

  // User list
  http.get(`${BASE_URL}/api/users`, () => {
    return HttpResponse.json(mockUsers);
  }),

  // Create user
  http.post(`${BASE_URL}/api/users`, async ({ request }) => {
    const body = (await request.json()) as Partial<User>;
    const newUser: User = {
      ...mockUser,
      id: "new-user-id",
      ...body,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    return HttpResponse.json(newUser, { status: 201 });
  }),

  // Invite users
  http.post(`${BASE_URL}/api/users/invite`, async ({ request }) => {
    const body = (await request.json()) as {
      emails: string[];
      isAdmin: boolean;
    };
    const invitedUsers = body.emails.map((email, index) => ({
      ...mockUser,
      id: `invited-${index}`,
      email,
      status: "pending",
    }));
    return HttpResponse.json({ success: invitedUsers, failed: [] });
  }),

  // Delete user
  http.delete(`${BASE_URL}/api/users/:id`, () => {
    return HttpResponse.json({ message: "User deleted" });
  }),

  // Enable user
  http.post(`${BASE_URL}/api/users/:id/enable`, () => {
    return HttpResponse.json({ message: "User enabled" });
  }),

  // Disable user
  http.post(`${BASE_URL}/api/users/:id/disable`, () => {
    return HttpResponse.json({ message: "User disabled" });
  }),
];
