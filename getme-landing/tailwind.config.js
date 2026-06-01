/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-mona)', 'Inter', 'system-ui', 'sans-serif'],
        mono: ['var(--font-jetbrains)', 'JetBrains Mono', 'monospace'],
        display: ['var(--font-mona)', 'Syne', 'sans-serif'],
      },
      colors: {
        blue: {
          950: '#020b18',
          900: '#040e20',
          850: '#061428',
          800: '#0a1f3d',
          700: '#102a52',
          600: '#1a3f78',
          500: '#2056a8',
          400: '#3477d4',
          300: '#5b9ee8',
          200: '#93c0f4',
          100: '#c8dff9',
          50: '#e8f2fd',
        },
        cyan: {
          400: '#22d3ee',
          300: '#67e8f9',
        }
      },
      boxShadow: {
        'glow-sm': '0 0 20px rgba(52,119,212,0.2)',
        'glow-md': '0 0 40px rgba(52,119,212,0.25)',
        'glow-lg': '0 0 80px rgba(32,86,168,0.4)',
      },
      animation: {
        'pulse-slow': 'pulse 2s ease-in-out infinite',
        'float': 'float 6s ease-in-out infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' },
        }
      }
    },
  },
  plugins: [],
}
