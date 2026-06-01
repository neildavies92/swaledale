import { Check, Pencil, X } from 'lucide-react';
import { FormEvent, useState } from 'react';
import type { Goal } from '../api/types';
import { money, percent } from '../api/format';
import { useGoals, useUpdateGoal } from '../api/hooks';
import { PageHeader } from '../components/PageHeader';
import { ErrorState, LoadingState } from '../components/State';

export function GoalsPage() {
  const goals = useGoals();
  const updateGoal = useUpdateGoal();

  if (goals.isLoading) {
    return <LoadingState label="Loading goals" />;
  }
  if (goals.isError) {
    return <ErrorState message={goals.error.message} />;
  }

  return (
    <section>
      <PageHeader title="Goals" eyebrow="Runway and emergency targets" />
      <div className="grid gap-4 lg:grid-cols-2">
        {goals.data!.map((goal) => (
          <GoalCard key={goal.id} goal={goal} pending={updateGoal.isPending} onSave={(input) => updateGoal.mutateAsync({ goalId: goal.id, input })} />
        ))}
      </div>
    </section>
  );
}

function GoalCard({ goal, pending, onSave }: { goal: Goal; pending: boolean; onSave: (input: Omit<Goal, 'id' | 'owner'>) => Promise<unknown> }) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState({ name: goal.name, target: String(goal.target.toFixed(2)), current: String(goal.current.toFixed(2)), notes: goal.notes });
  const progress = Math.min((goal.current / goal.target) * 100, 100);

  async function submit(event: FormEvent) {
    event.preventDefault();
    await onSave({ name: draft.name, target: Number(draft.target), current: Number(draft.current), notes: draft.notes });
    setEditing(false);
  }

  return (
    <div className="rounded-md border border-ink/10 bg-white p-4 shadow-panel">
      {editing ? (
        <form className="space-y-3" onSubmit={submit}>
          <input className="field" value={draft.name} onChange={(event) => setDraft({ ...draft, name: event.target.value })} aria-label="Goal name" />
          <div className="grid gap-3 sm:grid-cols-2">
            <input className="field" type="number" step="0.01" value={draft.target} onChange={(event) => setDraft({ ...draft, target: event.target.value })} aria-label="Target" />
            <input className="field" type="number" step="0.01" value={draft.current} onChange={(event) => setDraft({ ...draft, current: event.target.value })} aria-label="Current" />
          </div>
          <textarea className="field min-h-20" value={draft.notes} onChange={(event) => setDraft({ ...draft, notes: event.target.value })} aria-label="Notes" />
          <div className="flex justify-end gap-2">
            <button className="icon-button" type="submit" disabled={pending} aria-label="Save">
              <Check size={16} />
            </button>
            <button className="icon-button" type="button" onClick={() => setEditing(false)} aria-label="Cancel">
              <X size={16} />
            </button>
          </div>
        </form>
      ) : (
        <>
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="text-sm font-medium text-moss">{goal.owner}</p>
              <h3 className="font-semibold">{goal.name}</h3>
            </div>
            <button className="icon-button" type="button" onClick={() => setEditing(true)} aria-label={`Edit ${goal.name}`}>
              <Pencil size={15} />
            </button>
          </div>
          <div className="mt-4 flex items-baseline justify-between gap-3">
            <p className="text-2xl font-semibold">{money(goal.current)}</p>
            <p className="text-sm font-medium text-ink/55">of {money(goal.target)}</p>
          </div>
          <div className="mt-3 h-2 rounded-full bg-ink/10">
            <div className="h-2 rounded-full bg-moss" style={{ width: `${progress}%` }} />
          </div>
          <p className="mt-3 text-sm text-ink/60">{percent(goal.current, goal.target)} complete</p>
          {goal.notes ? <p className="mt-3 text-sm text-ink/60">{goal.notes}</p> : null}
        </>
      )}
    </div>
  );
}

