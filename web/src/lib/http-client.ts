import { useAuth } from "@/store/auth";
import ky, { type KyInstance } from "ky";

const API_BASE = "/api";

export class AuthError extends Error {
  constructor(message: string = "Authentication required") {
    super(message);
    this.name = "AuthError";
  }
}

class HttpClient {
  private client: KyInstance;
  private refreshPromise: Promise<void> | null = null;

  constructor() {
    this.client = ky.create({
      prefixUrl: API_BASE,
      credentials: "include",
      hooks: {
        beforeRequest: [
          (request) => {
            if (!(request.body instanceof FormData)) {
              request.headers.set("Content-Type", "application/json");
            }
          },
        ],
        beforeRetry: [
          async () => {
            try {
              await this.refreshToken();
              return;
            } catch (refreshError) {
              console.error("Token refresh failed:", refreshError);
              useAuth.getState().logout();
            }
          },
        ],
      },
      retry: {
        limit: 1,
        statusCodes: [401],
        methods: ["get", "post", "put", "patch", "delete"],
      },
    });
  }

  refreshToken(): Promise<void> {
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = this.client
      .post("auth/refresh", { retry: 0 })
      .then((response) => {
        if (!response.ok) {
          throw new AuthError("Failed to refresh token");
        }
      })
      .finally(() => {
        this.refreshPromise = null;
      });

    return this.refreshPromise;
  }

  get<T>(
    url: string,
    params?: Record<
      string,
      string | number | boolean | (string | number | boolean)[]
    >,
  ): Promise<T> {
    const p = new URLSearchParams();
    for (const [key, value] of Object.entries(params || {})) {
      p.set(key, String(value));
    }

    return this.client.get(url, { searchParams: p }).json<T>();
  }

  post<T>(url: string, data?: unknown): Promise<T> {
    const options = data instanceof FormData ? { body: data } : { json: data };

    return this.client.post(url, options).json<T>();
  }

  put<T>(url: string, data?: unknown): Promise<T> {
    const options = data instanceof FormData ? { body: data } : { json: data };

    return this.client.put(url, options).json<T>();
  }

  patch<T>(url: string, data?: unknown): Promise<T> {
    const options = data instanceof FormData ? { body: data } : { json: data };

    return this.client.patch(url, options).json<T>();
  }

  delete<T>(url: string): Promise<T> {
    return this.client.delete(url).json<T>();
  }
}

export const httpClient = new HttpClient();
