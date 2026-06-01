import type { ReactNode } from 'react';

export function MetricCard({ label, value, detail }: { label: string; value: ReactNode; detail?: string }) {
  return (
    <div className="rounded-md border border-ink/10 bg-white p-4 shadow-panel">
      <p className="text-sm font-medium text-ink/55">{label}</p>
      <div className="mt-2 text-2xl font-semibold tracking-normal">{value}</div>
      {detail ? <p className="mt-2 text-sm text-ink/55">{detail}</p> : null}
    </div>
  );
}

