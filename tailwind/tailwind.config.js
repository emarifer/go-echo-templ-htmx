const { fontFamily } = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["../views/**/*.{templ,go}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Kanit', ...fontFamily.sans]
      }
    },
  },
  plugins: [
    require("daisyui")
  ],
  daisyui: {
    themes: ["dark"]
  }
}

