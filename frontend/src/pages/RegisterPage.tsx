import { FormEvent, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useRegister } from '../api/hooks';
import { AuthShell } from '../components/AuthShell';

export function RegisterPage() {
  const register = useRegister();
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [householdName, setHouseholdName] = useState('');

  async function submit(event: FormEvent) {
    event.preventDefault();
    await register.mutateAsync({ name, email, password, householdName });
    navigate('/', { replace: true });
  }

  return (
    <AuthShell title="Create your account" eyebrow="Set up your household budget in minutes">
      <form className="space-y-4" onSubmit={submit}>
        <label className="block space-y-1.5">
          <span className="text-sm font-medium">Your name</span>
          <input className="field" autoComplete="name" required value={name} onChange={(event) => setName(event.target.value)} />
        </label>
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
            autoComplete="new-password"
            required
            minLength={8}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
          <span className="block text-xs text-ink/55">At least 8 characters</span>
        </label>
        <label className="block space-y-1.5">
          <span className="text-sm font-medium">Household name (optional)</span>
          <input
            className="field"
            placeholder="e.g. Swaledale Household"
            value={householdName}
            onChange={(event) => setHouseholdName(event.target.value)}
          />
        </label>
        {register.isError ? <p className="text-sm font-medium text-clay">{register.error.message}</p> : null}
        <button className="button-primary" type="submit" disabled={register.isPending}>
          {register.isPending ? 'Creating account…' : 'Create account'}
        </button>
      </form>
      <p className="mt-4 text-center text-sm text-ink/60">
        Already have an account?{' '}
        <Link className="font-medium text-moss hover:underline" to="/login">
          Sign in
        </Link>
      </p>
    </AuthShell>
  );
}
