import "@testing-library/jest-dom/vitest";
import { afterAll, afterEach, beforeAll, beforeEach } from "vitest";

import { server } from "./mocks/server";

// Set up base URL for relative URL resolution
const TEST_BASE_URL = "http://localhost:3000";

// Exported so tests can assert against the same mock instance without
// referencing it as an unbound method off `window.location`.
export const locationReplaceMock = vi.fn<(url: string | URL) => void>();

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
    replace: locationReplaceMock,
    assign: vi.fn<(url: string | URL) => void>(),
    reload: vi.fn<() => void>(),
  },
});

// Set document base URI for relative URL resolution
const baseElement = document.createElement("base");
baseElement.href = TEST_BASE_URL;
document.head.appendChild(baseElement);

// Mock window.matchMedia
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn<(query: string) => MediaQueryList>().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn<(listener: () => void) => void>(),
    removeListener: vi.fn<(listener: () => void) => void>(),
    addEventListener: vi.fn<(type: string, listener: () => void) => void>(),
    removeEventListener: vi.fn<(type: string, listener: () => void) => void>(),
    dispatchEvent: vi.fn<(event: Event) => boolean>(),
  })),
});

// Mock ResizeObserver
class ResizeObserverMock {
  observe = vi.fn<(target: Element) => void>();
  unobserve = vi.fn<(target: Element) => void>();
  disconnect = vi.fn<() => void>();
}
window.ResizeObserver = ResizeObserverMock;

// Mock IntersectionObserver
class IntersectionObserverMock {
  readonly root: Element | null = null;
  readonly rootMargin: string = "";
  readonly scrollMargin: string = "";
  readonly thresholds: ReadonlyArray<number> = [];
  observe = vi.fn<(target: Element) => void>();
  unobserve = vi.fn<(target: Element) => void>();
  disconnect = vi.fn<() => void>();
  takeRecords = vi.fn<() => IntersectionObserverEntry[]>().mockReturnValue([]);
}
window.IntersectionObserver = IntersectionObserverMock;

// Setup MSW server lifecycle
beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
beforeEach(() => {
  // Reset location.replace mock before each test
  locationReplaceMock.mockClear();
});
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
