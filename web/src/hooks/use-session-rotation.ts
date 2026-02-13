import { useEffect, useRef, useCallback } from "react";

import { httpClient } from "@/lib/http-client";

export const cookieUtils = {
  getSessionExpiry: () => {
    const expiryCookie = document.cookie
      .split("; ")
      .find((row) => row.startsWith("session_expiry="));

    if (!expiryCookie) return 0;

    const expiresStr = expiryCookie.split("=")[1];
    return expiresStr ? parseInt(expiresStr, 10) : 0;
  },

  hasSessionExpiry: () => {
    return document.cookie
      .split("; ")
      .some((row) => row.startsWith("session_expiry="));
  },
};

export function useSessionRotation() {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastExpiryRef = useRef<number>(0);
  const scheduleRotationRef = useRef<(() => void) | undefined>(undefined);

  const rotateToken = useCallback(async () => {
    try {
      await httpClient.post("account/sessions/rotate", undefined);

      return { success: true };
    } catch (error) {
      if (error instanceof Error && error.message?.includes("401")) {
        return { success: false, unauthorized: true };
      }

      return { success: false, error };
    }
  }, []);

  const scheduleRotation = useCallback(() => {
    // Clear any existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }

    // Check if we can schedule
    if (cookieUtils.getSessionExpiry() === 0) {
      return;
    }

    const expires = cookieUtils.getSessionExpiry();
    lastExpiryRef.current = expires;

    // because this job is scheduled for every tab we have open that shares a session we try
    // to distribute the scheduling of the job. For now this can be between 1 and 20 seconds
    const randomDelay = Math.floor(Math.random() * 19) + 1;
    const expiresWithDistribution = expires - randomDelay;

    // nextRun is when the job should be scheduled for in ms. setTimeout ms has a max value of 2147483647.
    const nextRun = Math.min(
      expiresWithDistribution * 1000 - Date.now(),
      2147483647,
    );

    timeoutRef.current = setTimeout(async () => {
      const currentExpiry = cookieUtils.getSessionExpiry();

      // Another tab already rotated the token
      if (currentExpiry > lastExpiryRef.current) {
        scheduleRotationRef.current?.();
        return;
      }

      const result = await rotateToken();

      if (result.success) {
        scheduleRotationRef.current?.();
      }
    }, nextRun);
  }, [rotateToken]);

  useEffect(() => {
    scheduleRotationRef.current = scheduleRotation;
  }, [scheduleRotation]);

  useEffect(() => {
    scheduleRotation();

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
        timeoutRef.current = null;
      }
    };
  }, [scheduleRotation]);

  return {
    scheduleRotation: scheduleRotation,
    cancelRotation: () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
        timeoutRef.current = null;
      }
    },
  };
}
