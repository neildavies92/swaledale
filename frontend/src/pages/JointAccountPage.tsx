import type { ReactNode } from 'react';
import { money } from '../api/format';
import { useJointAccount, useUpdateJointAccountItem } from '../api/hooks';
import { EditableMoneyRow } from '../components/EditableMoneyRow';
import { MetricCard } from '../components/MetricCard';
import { PageHeader } from '../components/PageHeader';
import { ErrorState, LoadingState } from '../components/State';

export function JointAccountPage() {
  const joint = useJointAccount();
  const updateItem = useUpdateJointAccountItem();

  if (joint.isLoading) {
    return <LoadingState label="Loading joint account" />;
  }
  if (joint.isError) {
    return <ErrorState message={joint.error.message} />;
  }

  const data = joint.data!;

  return (
    <section>
      <PageHeader title="Joint account" eyebrow="Shared household bills" />
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard label="Shared bills" value={money(data.total)} />
        <MetricCard label="Each per month" value={money(data.perPerson)} />
        <MetricCard label="Contributions" value={money(data.contributions.reduce((sum, contribution) => sum + contribution.amount, 0))} />
        <MetricCard label="Leftover" value={money(data.leftover)} />
      </div>

      <div className="mt-6 grid gap-5 xl:grid-cols-[1fr_22rem]">
        <div className="rounded-md border border-ink/10 bg-white shadow-panel">
          <div className="border-b border-ink/10 px-4 py-3">
            <h3 className="font-semibold">Shared bills</h3>
          </div>
          <div className="px-4">
            {data.items.map((item) => (
              <EditableMoneyRow
                key={item.id}
                label={item.label}
                amount={item.amount}
                sideValue={`${money(item.perPerson)} each`}
                disabled={updateItem.isPending}
                onSave={(input) => updateItem.mutateAsync({ itemId: item.id, input })}
              />
            ))}
          </div>
        </div>

        <div className="space-y-5">
          <Panel title="Wages">
            {data.wages.map((wage) => (
              <ReadonlyRow key={wage.id} label={wage.label} amount={wage.amount} />
            ))}
          </Panel>
          <Panel title="Contributions">
            {data.contributions.map((contribution) => (
              <ReadonlyRow key={contribution.memberId} label={contribution.member} amount={contribution.amount} />
            ))}
          </Panel>
        </div>
      </div>
    </section>
  );
}

function Panel({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="rounded-md border border-ink/10 bg-white p-4 shadow-panel">
      <h3 className="font-semibold">{title}</h3>
      <div className="mt-3 divide-y divide-ink/5">{children}</div>
    </div>
  );
}

function ReadonlyRow({ label, amount }: { label: string; amount: number }) {
  return (
    <div className="flex items-center justify-between gap-3 py-2">
      <span className="font-medium">{label}</span>
      <span className="font-semibold tabular-nums">{money(amount)}</span>
    </div>
  );
}
