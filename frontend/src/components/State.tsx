export function LoadingState({ label = 'Loading' }: { label?: string }) {
  return <div className="rounded-md border border-ink/10 bg-white p-6 text-sm text-ink/65 shadow-panel">{label}</div>;
}

export function ErrorState({ message }: { message: string }) {
  return <div className="rounded-md border border-clay/30 bg-clay/10 p-6 text-sm font-medium text-clay">{message}</div>;
}

