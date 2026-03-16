import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        canvas: '#eef4fb',
        ink: '#0f172a',
        muted: '#64748b',
        accent: '#0ea5e9',
        accentSoft: '#e0f2fe',
        accentStrong: '#0369a1',
        panel: '#ffffff',
        line: '#d9e2ec'
      },
      boxShadow: {
        shell: '0 24px 60px -28px rgba(15, 23, 42, 0.28)'
      },
      fontFamily: {
        sans: ['"Avenir Next"', '"Segoe UI"', 'ui-sans-serif', 'system-ui', 'sans-serif']
      }
    }
  },
  plugins: []
};

export default config;
