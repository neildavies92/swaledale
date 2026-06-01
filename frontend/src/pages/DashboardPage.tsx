import { Link } from 'react-router-dom';
import { money } from '../api/format';
import { useSummary } from '../api/hooks';
import { ErrorState, LoadingState } from '../components/State';
import { MetricCard } from '../components/MetricCard';
import { PageHeader } from '../components/PageHeader';

export function DashboardPage() {
  const summary = useSummary();

  if (summary.isLoading) {
    return <LoadingState label="Loading household dashboard" />;
  }
  if (summary.isError) {
    return <ErrorState message={summary.error.message} />;
  }

  const data = summary.data!;

  return (
    <section>
      <PageHeader title="Household dashboard" eyebrow={data.snapshot.name} />
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard label="Personal income" value={money(data.totalPersonalIncome)} detail="Neil + Katie personal tabs" />
        <MetricCard label="Committed" value={money(data.totalCommitted)} detail="Bills plus savings" />
        <MetricCard label="Remaining" value={money(data.totalRemaining)} detail="After committed spending" />
        <MetricCard label="Joint leftover" value={money(data.jointAccount.leftover)} detail="Contributions minus shared bills" />
      </div>

      <div className="mt-6 grid gap-5 xl:grid-cols-[1fr_24rem]">
        <div className="rounded-md border border-ink/10 bg-white shadow-panel">
          <div className="border-b border-ink/10 px-4 py-3">
            <h3 className="font-semibold">Members</h3>
          </div>
          <div className="divide-y divide-ink/10">
            {data.members.map((member) => (
              <Link key={member.member.id} to={`/members/${member.member.id}`} className="grid gap-3 px-4 py-4 transition hover:bg-ink/5 sm:grid-cols-5">
                <div className="sm:col-span-2">
                  <p className="font-semibold">{member.member.name}</p>
                  <p className="text-sm text-ink/55">{money(member.income)} income</p>
                </div>
                <SummaryValue label="Bills" value={member.billsTotal} />
                <SummaryValue label="Savings" value={member.savingsTotal} />
                <SummaryValue label="Remaining" value={member.remaining} />
              </Link>
            ))}
          </div>
        </div>

        <div className="rounded-md border border-ink/10 bg-white p-4 shadow-panel">
          <h3 className="font-semibold">Goals</h3>
          <div className="mt-3 space-y-3">
            {data.goals.slice(0, 4).map((goal) => (
              <div key={goal.id}>
                <div className="flex justify-between gap-3 text-sm">
                  <span className="font-medium">{goal.owner}: {goal.name}</span>
                  <span className="tabular-nums text-ink/60">{money(goal.target)}</span>
                </div>
                <div className="mt-2 h-2 rounded-full bg-ink/10">
                  <div className="h-2 rounded-full bg-moss" style={{ width: `${Math.min((goal.current / goal.target) * 100, 100)}%` }} />
                </div>
              </div>
            ))}
          </div>
          <Link className="mt-4 inline-flex text-sm font-semibold text-moss hover:text-ink" to="/goals">
            View all goals
          </Link>
        </div>
      </div>
    </section>
  );
}

function SummaryValue({ label, value }: { label: string; value: number }) {
  return (
    <div>
      <p className="text-xs font-medium uppercase tracking-normal text-ink/45">{label}</p>
      <p className="font-semibold tabular-nums">{money(value)}</p>
    </div>
  );
}

