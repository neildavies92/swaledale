import { Home, PiggyBank, Target, UsersRound } from 'lucide-react';
import type { ReactNode } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { useMembers } from './api/hooks';

export function App() {
  const members = useMembers();

  return (
    <div className="min-h-screen bg-paper text-ink">
      <aside className="fixed inset-y-0 left-0 hidden w-64 border-r border-ink/10 bg-white/75 px-4 py-5 backdrop-blur lg:block">
        <div className="mb-8">
          <p className="text-sm font-medium text-moss">Swaledale</p>
          <h1 className="text-2xl font-semibold tracking-normal">Household finance</h1>
        </div>
        <nav className="space-y-1">
          <NavItem to="/" icon={<Home size={18} />} label="Dashboard" />
          {members.data?.map((member) => (
            <NavItem key={member.id} to={`/members/${member.id}`} icon={<UsersRound size={18} />} label={member.name} />
          ))}
          <NavItem to="/joint-account" icon={<PiggyBank size={18} />} label="Joint account" />
          <NavItem to="/goals" icon={<Target size={18} />} label="Goals" />
        </nav>
      </aside>
      <main className="lg:pl-64">
        <div className="mx-auto max-w-7xl px-4 py-5 sm:px-6 lg:px-8">
          <div className="mb-5 flex items-center justify-between border-b border-ink/10 pb-4 lg:hidden">
            <div>
              <p className="text-sm font-medium text-moss">Swaledale</p>
              <h1 className="text-xl font-semibold">Household finance</h1>
            </div>
          </div>
          <Outlet />
        </div>
      </main>
    </div>
  );
}

function NavItem({ to, icon, label }: { to: string; icon: ReactNode; label: string }) {
  return (
    <NavLink
      to={to}
      end={to === '/'}
      className={({ isActive }) =>
        [
          'flex h-10 items-center gap-3 rounded-md px-3 text-sm font-medium transition',
          isActive ? 'bg-mint text-ink' : 'text-ink/70 hover:bg-ink/5 hover:text-ink',
        ].join(' ')
      }
    >
      {icon}
      <span>{label}</span>
    </NavLink>
  );
}
