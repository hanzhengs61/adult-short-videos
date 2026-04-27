/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js}'],
  theme: {
    extend: {
      colors: {
        bg: {
          DEFAULT: '#0d0d0f',
          surface: '#161618',
          card: '#1c1c1f',
          hover: '#222226',
        },
        border: { DEFAULT: '#2a2a2d', light: '#3a3a3e' },
        primary: {
          DEFAULT: '#FF2080',
          hover: '#FF44A0',
          light: '#FF88C8',
        },
        text: {
          primary: '#ffffff',
          secondary: '#a0a0a8',
          muted: '#606068',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      backgroundImage: {
        'gradient-primary': 'linear-gradient(135deg, #8B2FC9 0%, #FF2080 100%)',
        'gradient-card': 'linear-gradient(180deg, transparent 40%, rgba(0,0,0,0.95) 100%)',
      },
      boxShadow: {
        'card-hover': '0 8px 32px rgba(0,0,0,0.6), 0 0 0 1px rgba(255,32,128,0.25)',
        'glow': '0 0 24px rgba(255,32,128,0.45)',
      },
      animation: {
        'fade-in': 'fadeIn 0.2s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: { from: { opacity: '0' }, to: { opacity: '1' } },
        slideUp: { from: { opacity: '0', transform: 'translateY(12px)' }, to: { opacity: '1', transform: 'translateY(0)' } },
      },
    },
  },
  plugins: [],
}
