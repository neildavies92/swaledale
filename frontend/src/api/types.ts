export type Money = number;

export type Household = {
  id: number;
  name: string;
  currency: string;
};

export type Snapshot = {
  id: number;
  name: string;
  month: string;
  sourceUrl: string;
};

export type Member = {
  id: number;
  name: string;
};

export type BudgetItem = {
  id: number;
  memberId: number;
  categoryId: number;
  category: string;
  kind: 'bill' | 'saving';
  label: string;
  amount: Money;
};

export type AllocationRule = {
  id: number;
  memberId: number;
  label: string;
  percent: number;
  amount: Money;
};

export type MemberBudget = {
  member: Member;
  income: Money;
  items: BudgetItem[];
  allocations: AllocationRule[];
  billsTotal: Money;
  savingsTotal: Money;
  committedTotal: Money;
  remaining: Money;
  spendableIncome: Money;
};

export type Goal = {
  id: number;
  name: string;
  owner: string;
  target: Money;
  current: Money;
  notes: string;
};

export type JointAccountItem = {
  id: number;
  label: string;
  amount: Money;
  perPerson: Money;
};

export type JointContribution = {
  memberId: number;
  member: string;
  amount: Money;
};

export type IncomeEntry = {
  id: number;
  memberId: number;
  label: string;
  amount: Money;
};

export type JointAccount = {
  wages: IncomeEntry[];
  items: JointAccountItem[];
  contributions: JointContribution[];
  total: Money;
  perPerson: Money;
  leftover: Money;
};

export type Summary = {
  household: Household;
  snapshot: Snapshot;
  members: MemberBudget[];
  goals: Goal[];
  jointAccount: JointAccount;
  totalPersonalIncome: Money;
  totalCommitted: Money;
  totalRemaining: Money;
  totalSavings: Money;
};

export type MoneyLabelInput = {
  label: string;
  amount: Money;
};

export type GoalInput = {
  name: string;
  target: Money;
  current: Money;
  notes: string;
};

