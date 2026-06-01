import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from './client';
import type { GoalInput, MoneyLabelInput } from './types';

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

