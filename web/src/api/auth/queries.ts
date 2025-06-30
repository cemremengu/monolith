import {
  useMutation,
  useQueryClient,
  queryOptions,
} from "@tanstack/react-query";
import { authApi } from "./index";

export const authKeys = {
  all: ["auth"] as const,
  profile: () => [...authKeys.all, "profile"] as const,
};

export const profileQueryOptions = queryOptions({
  queryKey: authKeys.profile(),
  queryFn: authApi.profile,
});

export const useUpdatePreferences = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authApi.updatePreferences,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.profile() });
    },
  });
};
