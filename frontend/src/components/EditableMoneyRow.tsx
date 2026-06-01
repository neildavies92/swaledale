import { Check, Pencil, X } from 'lucide-react';
import { FormEvent, useState } from 'react';
import { money } from '../api/format';

type Props = {
  label: string;
  amount: number;
  sideValue?: string;
  disabled?: boolean;
  onSave: (input: { label: string; amount: number }) => Promise<unknown> | void;
};

export function EditableMoneyRow({ label, amount, sideValue, disabled, onSave }: Props) {
  const [editing, setEditing] = useState(false);
  const [draftLabel, setDraftLabel] = useState(label);
  const [draftAmount, setDraftAmount] = useState(String(amount.toFixed(2)));

  async function submit(event: FormEvent) {
    event.preventDefault();
    await onSave({ label: draftLabel.trim(), amount: Number(draftAmount) });
    setEditing(false);
  }

  if (editing) {
    return (
      <form className="grid grid-cols-[1fr_7rem_5rem] gap-2 py-2" onSubmit={submit}>
        <input className="field" value={draftLabel} onChange={(event) => setDraftLabel(event.target.value)} aria-label="Label" />
        <input
          className="field"
          type="number"
          step="0.01"
          min="0"
          value={draftAmount}
          onChange={(event) => setDraftAmount(event.target.value)}
          aria-label="Amount"
        />
        <div className="flex justify-end gap-1">
          <button className="icon-button" type="submit" disabled={disabled} aria-label="Save">
            <Check size={16} />
          </button>
          <button className="icon-button" type="button" onClick={() => setEditing(false)} aria-label="Cancel">
            <X size={16} />
          </button>
        </div>
      </form>
    );
  }

  return (
    <div className="grid grid-cols-[1fr_auto_2.25rem] items-center gap-3 border-b border-ink/5 py-3 last:border-b-0">
      <div>
        <p className="font-medium">{label}</p>
        {sideValue ? <p className="text-sm text-ink/55">{sideValue}</p> : null}
      </div>
      <p className="font-semibold tabular-nums">{money(amount)}</p>
      <button className="icon-button" type="button" onClick={() => setEditing(true)} aria-label={`Edit ${label}`}>
        <Pencil size={15} />
      </button>
    </div>
  );
}
