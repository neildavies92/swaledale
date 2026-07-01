import type {
  Goal,
  GoalInput,
  JointAccount,
  LoginInput,
  Member,
  MemberBudget,
  MoneyLabelInput,
  RegisterInput,
  Summary,
  User,
} from './types';

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

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
    throw new ApiError(response.status, body.error ?? 'Request failed');
  }
  return response.json() as Promise<T>;
}

export const api = {
  me: () => request<User>('/api/auth/me'),
  register: (input: RegisterInput) =>
    request<User>('/api/auth/register', { method: 'POST', body: JSON.stringify(input) }),
  login: (input: LoginInput) =>
    request<User>('/api/auth/login', { method: 'POST', body: JSON.stringify(input) }),
  logout: () => request<{ status: string }>('/api/auth/logout', { method: 'POST' }),
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
