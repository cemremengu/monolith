import { useAuth } from "@/hooks/use-auth";
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

  constructor() {
    this.client = ky.create({
      prefixUrl: API_BASE,
      credentials: "include",
      hooks: {
        beforeRequest: [
          (request) => {
            if (request.body === null) {
              request.headers.set("Content-Type", "application/json");
            }
          },
        ],
        beforeError: [
          (error) => {
            if (
              error.response?.status === 401 &&
              useAuth.getState().isLoggedIn
            ) {
              useAuth.getState().logout();
            }
            return error;
          },
        ],
      },
      retry: 0,
    });
  }

  get<T>(
    url: string,
    searchParams?: Record<
      string,
      string | number | boolean | (string | number | boolean)[]
    >,
    options?: HttpClientOptions,
  ): Promise<T> {
    const hasSearchParams =
      searchParams && Object.keys(searchParams).length > 0;

    if (!hasSearchParams) {
      return this.client.get(url, options).json<T>();
    }

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
