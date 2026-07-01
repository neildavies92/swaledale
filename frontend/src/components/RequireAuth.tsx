import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useMe } from '../api/hooks';
import { ErrorState, LoadingState } from './State';

export function RequireAuth({ children }: { children: ReactNode }) {
  const me = useMe();

  if (me.isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-paper px-4">
        <LoadingState label="Checking your session" />
      </div>
    );
  }
  if (me.isError) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-paper px-4">
        <ErrorState message={me.error.message} />
      </div>
    );
  }
  if (!me.data) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}
