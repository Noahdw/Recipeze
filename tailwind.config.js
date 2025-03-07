/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      "./**/*.go",
      "./templates/**/*.html"
    ],
    theme: {
      extend: {},
    },
    plugins: [
      // If you need typography plugin, you'll need to install it
      require('@tailwindcss/typography'),
      function({ addUtilities }) {
        const newUtilities = {
          '.scrollbar-stable': {
            'overflow-y': 'scroll',
          },
        }
        addUtilities(newUtilities)
      }
    ],
  }