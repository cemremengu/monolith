import "@testing-library/jest-dom/vitest";
import { afterAll, afterEach, beforeAll, beforeEach } from "vitest";

import { server } from "./mocks/server";

// Set up base URL for relative URL resolution
const TEST_BASE_URL = "http://localhost:3000";

// Mock window.location with proper origin for URL resolution
Object.defineProperty(window, "location", {
  writable: true,
  value: {
    href: TEST_BASE_URL,
    origin: TEST_BASE_URL,
    protocol: "http:",
    host: "localhost:3000",
    hostname: "localhost",
    port: "3000",
    pathname: "/",
    search: "",
    hash: "",
    replace: vi.fn(),
    assign: vi.fn(),
    reload: vi.fn(),
  },
});

// Set document base URI for relative URL resolution
const baseElement = document.createElement("base");
baseElement.href = TEST_BASE_URL;
document.head.appendChild(baseElement);

// Mock window.matchMedia
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock ResizeObserver
class ResizeObserverMock {
  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
}
window.ResizeObserver = ResizeObserverMock;

// Mock IntersectionObserver
class IntersectionObserverMock {
  readonly root: Element | null = null;
  readonly rootMargin: string = "";
  readonly thresholds: ReadonlyArray<number> = [];
  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
  takeRecords = vi.fn().mockReturnValue([]);
}
window.IntersectionObserver = IntersectionObserverMock;

// Setup MSW server lifecycle
beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
beforeEach(() => {
  // Reset location.replace mock before each test
  vi.mocked(window.location.replace).mockClear();
});
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
