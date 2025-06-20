import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from './api';
import type { CreateUserRequest } from '@/types';

// Query keys
export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (filters: string) => [...userKeys.lists(), { filters }] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
};

// Query hooks
export const useUsers = () => {
  return useQuery({
    queryKey: userKeys.lists(),
    queryFn: api.users.getAll,
  });
};

export const useUser = (id: string) => {
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: () => api.users.getById(id),
    enabled: !!id,
  });
};

// Mutation hooks
export const useCreateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: api.users.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
};

export const useUpdateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateUserRequest }) =>
      api.users.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
      queryClient.invalidateQueries({ queryKey: userKeys.detail(variables.id) });
    },
  });
};

export const useDeleteUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: api.users.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
};