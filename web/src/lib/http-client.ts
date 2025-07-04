import { useAuth } from "@/store/auth";
import ky, { type KyInstance, type Options } from "ky";

const API_BASE = "/api";

type HttpClientOptions = Pick<Options, "headers" | "retry" | "timeout">;

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
    searchParams?: Record<
      string,
      string | number | boolean | (string | number | boolean)[]
    >,
    options?: HttpClientOptions,
  ): Promise<T> {
    const p = new URLSearchParams();
    for (const [key, value] of Object.entries(searchParams || {})) {
      p.set(key, String(value));
    }

    return this.client.get(url, { searchParams: p, ...options }).json<T>();
  }

  post<T>(
    url: string,
    data?: unknown,
    options?: HttpClientOptions,
  ): Promise<T> {
    const requestOptions =
      data instanceof FormData ? { body: data } : { json: data };

    return this.client.post(url, { ...requestOptions, ...options }).json<T>();
  }

  put<T>(url: string, data?: unknown, options?: HttpClientOptions): Promise<T> {
    const requestOptions =
      data instanceof FormData ? { body: data } : { json: data };

    return this.client.put(url, { ...requestOptions, ...options }).json<T>();
  }

  patch<T>(
    url: string,
    data?: unknown,
    options?: HttpClientOptions,
  ): Promise<T> {
    const requestOptions =
      data instanceof FormData ? { body: data } : { json: data };

    return this.client.patch(url, { ...requestOptions, ...options }).json<T>();
  }

  delete<T>(url: string, options?: HttpClientOptions): Promise<T> {
    return this.client.delete(url, options).json<T>();
  }
}

export const httpClient = new HttpClient();
