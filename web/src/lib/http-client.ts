import axios from "axios";
import qs from "qs";
import { getAuthState } from "@/context/auth";

const API_BASE = "/api";

const httpClient = axios.create({
  baseURL: API_BASE,
  withCredentials: true,
  paramsSerializer: (params) => {
    return qs.stringify(params, { arrayFormat: "comma" });
  },
});

let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

const refreshToken = (): Promise<void> => {
  if (isRefreshing) {
    return refreshPromise!;
  }

  isRefreshing = true;
  refreshPromise = httpClient
    .post("/auth/refresh")
    .then(() => {})
    .finally(() => {
      isRefreshing = false;
      refreshPromise = null;
    });

  return refreshPromise;
};

httpClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        await refreshToken();
        return httpClient(originalRequest);
      } catch (refreshError) {
        console.error("Token refresh failed:", refreshError);
        getAuthState().setUnauthenticated();
        throw new Error("Authentication required");
      }
    }

    if (error.response?.data) {
      throw new Error(error.response.data.message || error.response.data);
    }

    throw error;
  },
);

export { httpClient };
