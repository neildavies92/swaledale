import { FormEvent, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useLogin } from '../api/hooks';
import { AuthShell } from '../components/AuthShell';

export function LoginPage() {
  const login = useLogin();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  async function submit(event: FormEvent) {
    event.preventDefault();
    await login.mutateAsync({ email, password });
    navigate('/', { replace: true });
  }

  return (
    <AuthShell title="Sign in" eyebrow="Household finance, without the spreadsheet">
      <form className="space-y-4" onSubmit={submit}>
        <label className="block space-y-1.5">
          <span className="text-sm font-medium">Email</span>
          <input
            className="field"
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={(event) => setEmail(event.target.value)}
          />
        </label>
        <label className="block space-y-1.5">
          <span className="text-sm font-medium">Password</span>
          <input
            className="field"
            type="password"
            autoComplete="current-password"
            required
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </label>
        {login.isError ? <p className="text-sm font-medium text-clay">{login.error.message}</p> : null}
        <button className="button-primary" type="submit" disabled={login.isPending}>
          {login.isPending ? 'Signing in…' : 'Sign in'}
        </button>
      </form>
      <p className="mt-4 text-center text-sm text-ink/60">
        New here?{' '}
        <Link className="font-medium text-moss hover:underline" to="/register">
          Create an account
        </Link>
      </p>
    </AuthShell>
  );
}
