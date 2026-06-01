import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        ink: '#17211c',
        moss: '#55745f',
        mint: '#dce9df',
        paper: '#f7f4ef',
        clay: '#b86445',
      },
      boxShadow: {
        panel: '0 10px 30px rgba(23, 33, 28, 0.08)',
      },
    },
  },
  plugins: [],
} satisfies Config;

