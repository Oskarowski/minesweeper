/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./public/**/*.{html,js}"],
    theme: {
        extend: {},
    },
    plugins: [],
    safelist: [
        "bg-blue-600",
        "bg-blue-700",
        "bg-red-600",
        "bg-red-700",
        "text-white",
        "rounded",
        "hover:bg-blue-700",
        "hover:bg-red-700",
    ],
};

// npx tailwindcss -i ./main.css -o ./dist/tailwind.css --watch
