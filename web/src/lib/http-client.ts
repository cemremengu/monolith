import ky, { type KyInstance, type Options } from "ky";

import { useAuth } from "@/hooks/use-auth";

type HttpClientOptions = Pick<Options, "headers" | "retry" | "timeout">;

export class AuthError extends Error {
  constructor(message: string = "Authentication required") {
    super(message);
    this.name = "AuthError";
  }
}

// Use absolute URL to support both browser and Node.js test environments
function getApiBase(): string {
  if (typeof window !== "undefined" && window.location?.origin) {
    return `${window.location.origin}/api`;
  }
  return "/api";
}

class HttpClient {
  private _client: KyInstance | null = null;

  private get client(): KyInstance {
    if (!this._client) {
      this._client = ky.create({
        prefixUrl: getApiBase(),
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
                useAuth.getState().logout({ redirectToLogin: true });
              }
              return error;
            },
          ],
        },
        retry: 0,
      });
    }
    return this._client;
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
