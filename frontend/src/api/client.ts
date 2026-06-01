import type { Goal, GoalInput, JointAccount, Member, MemberBudget, MoneyLabelInput, Summary } from './types';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  });
  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(body.error ?? 'Request failed');
  }
  return response.json() as Promise<T>;
}

export const api = {
  summary: () => request<Summary>('/api/summary'),
  members: () => request<Member[]>('/api/members'),
  memberBudget: (memberId: number) => request<MemberBudget>(`/api/members/${memberId}/budget`),
  updateBudgetItem: (memberId: number, itemId: number, input: MoneyLabelInput) =>
    request<MemberBudget>(`/api/members/${memberId}/budget-items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    }),
  jointAccount: () => request<JointAccount>('/api/joint-account'),
  updateJointAccountItem: (itemId: number, input: MoneyLabelInput) =>
    request<JointAccount>(`/api/joint-account/items/${itemId}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    }),
  goals: () => request<Goal[]>('/api/goals'),
  updateGoal: (goalId: number, input: GoalInput) =>
    request<Goal[]>(`/api/goals/${goalId}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    }),
};

