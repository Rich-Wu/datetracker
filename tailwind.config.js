/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/**/*.{tmpl,html}'],
  theme: {
    extend: {
      animation: {
        'beat-once': 'beat 0.1s linear 1'
      },
      keyframes: {
        beat: {
          '0%, 100%': { transform: 'scale(1)' },
          '50%': { transform: 'scale(1.05)' }
        }
      }
    },
  },
  plugins: [],
}

