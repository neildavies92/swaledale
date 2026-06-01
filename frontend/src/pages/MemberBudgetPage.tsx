import { useParams } from 'react-router-dom';
import { money } from '../api/format';
import { useMemberBudget, useUpdateBudgetItem } from '../api/hooks';
import { EditableMoneyRow } from '../components/EditableMoneyRow';
import { MetricCard } from '../components/MetricCard';
import { PageHeader } from '../components/PageHeader';
import { ErrorState, LoadingState } from '../components/State';

export function MemberBudgetPage() {
  const memberId = Number(useParams().memberId);
  const budget = useMemberBudget(memberId);
  const updateItem = useUpdateBudgetItem(memberId);

  if (budget.isLoading) {
    return <LoadingState label="Loading member budget" />;
  }
  if (budget.isError) {
    return <ErrorState message={budget.error.message} />;
  }

  const data = budget.data!;
  const bills = data.items.filter((item) => item.kind === 'bill');
  const savings = data.items.filter((item) => item.kind === 'saving');

  return (
    <section>
      <PageHeader title={`${data.member.name}'s budget`} eyebrow="Member budget" />
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard label="Income" value={money(data.income)} />
        <MetricCard label="Bills" value={money(data.billsTotal)} />
        <MetricCard label="Savings" value={money(data.savingsTotal)} />
        <MetricCard label="Remaining" value={money(data.remaining)} />
      </div>

      <div className="mt-6 grid gap-5 xl:grid-cols-2">
        <BudgetPanel
          title="Bills"
          items={bills}
          pending={updateItem.isPending}
          onSave={(itemId, input) => updateItem.mutateAsync({ itemId, input })}
        />
        <BudgetPanel
          title="Savings / Investments"
          items={savings}
          pending={updateItem.isPending}
          onSave={(itemId, input) => updateItem.mutateAsync({ itemId, input })}
        />
      </div>

      <div className="mt-5 rounded-md border border-ink/10 bg-white p-4 shadow-panel">
        <h3 className="font-semibold">Allocation rules</h3>
        <div className="mt-3 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
          {data.allocations.map((allocation) => (
            <div key={allocation.id} className="rounded-md bg-paper p-3">
              <p className="text-sm font-medium text-ink/60">{allocation.label}</p>
              <p className="mt-1 text-lg font-semibold">{money(allocation.amount)}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function BudgetPanel({
  title,
  items,
  pending,
  onSave,
}: {
  title: string;
  items: { id: number; label: string; amount: number }[];
  pending: boolean;
  onSave: (itemId: number, input: { label: string; amount: number }) => Promise<unknown>;
}) {
  return (
    <div className="rounded-md border border-ink/10 bg-white shadow-panel">
      <div className="border-b border-ink/10 px-4 py-3">
        <h3 className="font-semibold">{title}</h3>
      </div>
      <div className="px-4">
        {items.map((item) => (
          <EditableMoneyRow key={item.id} label={item.label} amount={item.amount} disabled={pending} onSave={(input) => onSave(item.id, input)} />
        ))}
      </div>
    </div>
  );
}

