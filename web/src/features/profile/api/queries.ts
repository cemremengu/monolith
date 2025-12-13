import {
  useMutation,
  useQueryClient,
  queryOptions,
  useSuspenseQuery,
} from "@tanstack/react-query";

import { accountApi } from "./index";

export const accountKeys = {
  all: ["account"] as const,
  profile: () => [...accountKeys.all, "profile"] as const,
  sessions: () => [...accountKeys.all, "sessions"] as const,
};

export const profileQueryOptions = queryOptions({
  queryKey: accountKeys.profile(),
  queryFn: accountApi.profile,
});

export const useProfile = () => {
  return useSuspenseQuery(profileQueryOptions);
};

export const useUpdatePreferences = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.updatePreferences,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: accountKeys.profile() });
    },
  });
};

export const useRevokeSession = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.revokeSession,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: accountKeys.sessions() });
    },
  });
};

export const useRevokeAllOtherSessions = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: accountApi.revokeAllOtherSessions,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: accountKeys.sessions() });
    },
  });
};
