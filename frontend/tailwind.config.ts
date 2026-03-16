import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        canvas: '#F4F1EB',
        panel: '#F8F6F2',
        shell: '#ECE6DB',
        line: '#D8D0C5',
        lineMuted: '#C9C0B4',
        ink: '#1F2430',
        rail: '#2E3138',
        muted: '#6F675E',
        dusty: '#8D857C',
        accent: '#2F3440',
        accentSoft: '#D7CEC1',
        accentStrong: '#A88B52',
        highlight: '#BCA98B',
        successMuted: '#7D8A74',
        warningMuted: '#B4914E',
        errorMuted: '#B96F62',
        infoMuted: '#7A8CA3'
      },
      boxShadow: {
        shell: '0 28px 80px -40px rgba(31, 36, 48, 0.35)',
        float: '0 20px 48px -26px rgba(31, 36, 48, 0.22)'
      },
      fontFamily: {
        sans: ['"Avenir Next"', '"Segoe UI"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        display: ['"Iowan Old Style"', '"Baskerville"', '"Times New Roman"', 'serif']
      }
    }
  },
  plugins: []
};

export default config;
