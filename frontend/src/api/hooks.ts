import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api, ApiError } from './client';
import type { GoalInput, LoginInput, MoneyLabelInput, RegisterInput, User } from './types';

export function useMe() {
  return useQuery<User | null>({
    queryKey: ['me'],
    queryFn: async () => {
      try {
        return await api.me();
      } catch (error) {
        if (error instanceof ApiError && error.status === 401) {
          return null;
        }
        throw error;
      }
    },
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: LoginInput) => api.login(input),
    onSuccess: (user) => {
      queryClient.clear();
      queryClient.setQueryData(['me'], user);
    },
  });
}

export function useRegister() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: RegisterInput) => api.register(input),
    onSuccess: (user) => {
      queryClient.clear();
      queryClient.setQueryData(['me'], user);
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: api.logout,
    onSuccess: () => {
      queryClient.clear();
      queryClient.setQueryData(['me'], null);
    },
  });
}

export function useSummary() {
  return useQuery({ queryKey: ['summary'], queryFn: api.summary });
}

export function useMembers() {
  return useQuery({ queryKey: ['members'], queryFn: api.members });
}

export function useMemberBudget(memberId: number) {
  return useQuery({ queryKey: ['member-budget', memberId], queryFn: () => api.memberBudget(memberId), enabled: Number.isFinite(memberId) });
}

export function useUpdateBudgetItem(memberId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, input }: { itemId: number; input: MoneyLabelInput }) => api.updateBudgetItem(memberId, itemId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['member-budget', memberId] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
    },
  });
}

export function useJointAccount() {
  return useQuery({ queryKey: ['joint-account'], queryFn: api.jointAccount });
}

export function useUpdateJointAccountItem() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, input }: { itemId: number; input: MoneyLabelInput }) => api.updateJointAccountItem(itemId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['joint-account'] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
    },
  });
}

export function useGoals() {
  return useQuery({ queryKey: ['goals'], queryFn: api.goals });
}

export function useUpdateGoal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ goalId, input }: { goalId: number; input: GoalInput }) => api.updateGoal(goalId, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['goals'] });
      queryClient.invalidateQueries({ queryKey: ['summary'] });
    },
  });
}

