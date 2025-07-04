import {
  useMutation,
  useQueryClient,
  queryOptions,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { usersApi } from "./index";
import type { CreateUserRequest } from "@/types/api";

export const userKeys = {
  all: ["users"] as const,
  lists: () => [...userKeys.all, "list"] as const,
  list: (filters: string) => [...userKeys.lists(), { filters }] as const,
  details: () => [...userKeys.all, "detail"] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
};

export const usersQueryOptions = (opts: {
  filterBy?: string;
  sortBy?: "name" | "email";
}) =>
  queryOptions({
    queryKey: userKeys.lists(),
    queryFn: () => usersApi.getAll(opts),
  });

export const userQueryOptions = (userId: string) =>
  queryOptions({
    queryKey: userKeys.detail(userId),
    queryFn: () => usersApi.getById(userId),
    enabled: !!userId,
  });

export const useUsers = (opts: {
  filterBy?: string;
  sortBy?: "name" | "email";
}) => {
  return useSuspenseQuery(usersQueryOptions(opts));
};

export const useCreateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: usersApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateUserRequest }) =>
      usersApi.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: userKeys.detail(variables.id),
      });
    },
  });
};

export const useDeleteUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: usersApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
};
