import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useMe } from '../api/hooks';
import { LoadingState } from './State';

export function AuthShell({ title, eyebrow, children }: { title: string; eyebrow: string; children: ReactNode }) {
  const me = useMe();

  if (me.isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-paper px-4">
        <LoadingState label="Loading" />
      </div>
    );
  }
  if (me.data) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-paper px-4 text-ink">
      <div className="w-full max-w-md">
        <div className="mb-6 text-center">
          <p className="text-sm font-medium text-moss">Swaledale</p>
          <h1 className="text-2xl font-semibold">{title}</h1>
          <p className="mt-1 text-sm text-ink/60">{eyebrow}</p>
        </div>
        <div className="rounded-md border border-ink/10 bg-white p-6 shadow-panel">{children}</div>
      </div>
    </div>
  );
}
