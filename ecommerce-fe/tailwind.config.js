/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      maxWidth: {
        "screen-max": "1440px",
      },
      colors: {
        "dark-gray": "#303030",
        "ocean-green": "#4EA674",
        cyprus: "#023337",
        "surf-crest": "#C1E6BA",
        "aqua-spring": "#EAF8E7",
        primary: "#6467F2",
        success: "#21C45D",
        pending: "#F0D411",
        error: "#EF4343",
        grey: "#7C7C7C",
        "ocean-blue": "#0569B5",
      },
      spacing: {
        4: "4px",
        6: "6px",
        8: "8px",
        12: "12px",
        16: "16px",
        20: "20px",
        24: "24px",
        32: "32px",
        40: "40px",
        48: "48px",
        56: "56px",
        64: "64px",
      },
      fontFamily: {
        lato: ["Lato", "sans-serif"],
      },
      fontSize: {
        body: "16px",
        button: "14px",
        title: "22px",
        dashboard: "18px",
        h2: "32px",
        h1: "57px",
      },
      lineHeight: {
        body: "1.4",
      },
      letterSpacing: {
        header: "-0.25%",
        title: "-2%",
        button: "0",
      },
      boxShadow: {
        "ambient-1": "0 2px 4px rgba(0, 0, 0, 0.05)",
        "ambient-3": "0 4px 8px rgba(0, 0, 0, 0.1)",
        "ambient-6": "0 8px 16px rgba(0, 0, 0, 0.15)",
      },
      borderRadius: {
        default: "8px",
      },
    },
  },
  plugins: [],
};
