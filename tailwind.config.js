/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{tmpl,html}", "./static/**/*.js"],
  theme: {
    extend: {
      gridTemplateColumns: {
        '1fr-2fr': '1fr 2fr',
        '2fr-3fr': '2fr 3fr'
      },
      animation: {
        "beat-once": "beat 0.2s linear 1",
        "beacon": "beacon 1s linear infinite"
      },
      keyframes: {
        beat: {
          "0%, 100%": { transform: "scale(1)" },
          "50%": { transform: "scale(1.01)" },
        },
        beacon: {
          "0%, 100%": { transform: "scale(1)" },
          "50%": { transform: "scale(1.03)" }
        }
      },
    },
  },
  plugins: [],
};
