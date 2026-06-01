import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import { App } from './App';
import { DashboardPage } from './pages/DashboardPage';
import { MemberBudgetPage } from './pages/MemberBudgetPage';
import { JointAccountPage } from './pages/JointAccountPage';
import { GoalsPage } from './pages/GoalsPage';
import './styles.css';

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: 'members/:memberId', element: <MemberBudgetPage /> },
      { path: 'joint-account', element: <JointAccountPage /> },
      { path: 'goals', element: <GoalsPage /> },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </React.StrictMode>,
);

