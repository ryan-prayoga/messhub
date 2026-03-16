import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        canvas: '#f8fafc',
        ink: '#0f172a',
        accent: '#f97316',
        accentSoft: '#fed7aa',
        panel: '#ffffff',
        line: '#e2e8f0'
      },
      boxShadow: {
        shell: '0 18px 60px -30px rgba(15, 23, 42, 0.35)'
      },
      fontFamily: {
        sans: ['"Avenir Next"', '"Segoe UI"', 'ui-sans-serif', 'system-ui', 'sans-serif']
      }
    }
  },
  plugins: []
};

export default config;
