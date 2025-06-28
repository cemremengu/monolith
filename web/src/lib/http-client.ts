import { API_BASE, getHeaders } from "@/api/config";

interface RequestOptions extends RequestInit {
  url: string;
}

class HttpClient {
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;

  private async refreshToken(): Promise<void> {
    if (this.isRefreshing) {
      return this.refreshPromise!;
    }

    this.isRefreshing = true;
    this.refreshPromise = fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: getHeaders(),
      credentials: "include",
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error("Failed to refresh token");
        }
      })
      .finally(() => {
        this.isRefreshing = false;
        this.refreshPromise = null;
      });

    return this.refreshPromise;
  }

  async request<T>(options: RequestOptions): Promise<T> {
    const { url, ...fetchOptions } = options;

    const makeRequest = async (): Promise<Response> => {
      return fetch(url, {
        headers: getHeaders(),
        credentials: "include",
        ...fetchOptions,
      });
    };

    let response = await makeRequest();

    // If we get a 401 try to refresh
    if (response.status === 401) {
      try {
        await this.refreshToken();
        // Retry the original request
        response = await makeRequest();
      } catch (refreshError) {
        console.error("Token refresh failed:", refreshError);
        throw new Error("Authentication required");
      }
    }

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(
        errorText || `Request failed with status ${response.status}`,
      );
    }

    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
      return response.json();
    }

    return response.text() as T;
  }

  get<T>(
    url: string,
    options?: Omit<RequestOptions, "url" | "method">,
  ): Promise<T> {
    return this.request<T>({ url, method: "GET", ...options });
  }

  post<T>(
    url: string,
    data?: unknown,
    options?: Omit<RequestOptions, "url" | "method">,
  ): Promise<T> {
    return this.request<T>({
      url,
      method: "POST",
      body: JSON.stringify(data),
      ...options,
    });
  }

  put<T>(
    url: string,
    data?: unknown,
    options?: Omit<RequestOptions, "url" | "method">,
  ): Promise<T> {
    return this.request<T>({
      url,
      method: "PUT",
      body: JSON.stringify(data),
      ...options,
    });
  }

  patch<T>(
    url: string,
    data?: unknown,
    options?: Omit<RequestOptions, "url" | "method">,
  ): Promise<T> {
    return this.request<T>({
      url,
      method: "PATCH",
      body: JSON.stringify(data),
      ...options,
    });
  }

  delete<T>(
    url: string,
    options?: Omit<RequestOptions, "url" | "method">,
  ): Promise<T> {
    return this.request<T>({ url, method: "DELETE", ...options });
  }
}

export const httpClient = new HttpClient();
