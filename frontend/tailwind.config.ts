import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./app/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: "#f0f5ff",
          100: "#e0ebff",
          200: "#b8d4ff",
          300: "#85b8ff",
          400: "#4d96ff",
          500: "#2563eb",
          600: "#1d4ed8",
          700: "#1e3a8a",
          800: "#1e3a5f",
          900: "#1a365d",
        },
      },
    },
  },
  plugins: [],
};

export default config;
