/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        critical: "#ef4444",
        warning: "#f59e0b",
        info: "#3b82f6",
        success: "#22c55e",
      },
    },
  },
  plugins: [],
};
