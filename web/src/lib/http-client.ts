import { API_BASE, getHeaders } from "@/api/config";
import { useAuth } from "@/store/auth";

interface RequestOptions extends RequestInit {
  url: string;
  params?: Record<string, string | number | boolean>;
  formData?: FormData;
}

class HttpClient {
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;

  private refreshToken(): Promise<void> {
    if (this.isRefreshing) {
      return this.refreshPromise!;
    }

    this.isRefreshing = true;
    this.refreshPromise = fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: getHeaders(),
      credentials: "include",
    })
      .then((response) => {
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

  private buildUrl(
    url: string,
    params?: Record<string, string | number | boolean>,
  ): string {
    if (!params) return url;

    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      searchParams.append(key, String(value));
    });

    return `${url}${url.includes("?") ? "&" : "?"}${searchParams.toString()}`;
  }

  async request<T>(options: RequestOptions): Promise<T> {
    const { url, params, formData, ...fetchOptions } = options;
    const finalUrl = this.buildUrl(url, params);

    const makeRequest = (): Promise<Response> => {
      const headers = formData ? {} : getHeaders();

      return fetch(finalUrl, {
        headers,
        credentials: "include",
        body: formData || fetchOptions.body,
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
        // Mark user as unauthenticated when token refresh fails
        useAuth.getState().setUnauthenticated();
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
    params?: Record<string, string | number | boolean>,
    options?: Omit<RequestOptions, "url" | "method" | "params">,
  ): Promise<T> {
    return this.request<T>({ url, method: "GET", params, ...options });
  }

  post<T>(
    url: string,
    data?: unknown,
    options?: Omit<RequestOptions, "url" | "method">,
  ): Promise<T> {
    const { formData, ...restOptions } = options || {};

    if (formData) {
      return this.request<T>({
        url,
        method: "POST",
        formData,
        ...restOptions,
      });
    }

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
    const { formData, ...restOptions } = options || {};

    if (formData) {
      return this.request<T>({
        url,
        method: "PUT",
        formData,
        ...restOptions,
      });
    }

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
    const { formData, ...restOptions } = options || {};

    if (formData) {
      return this.request<T>({
        url,
        method: "PATCH",
        formData,
        ...restOptions,
      });
    }

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
