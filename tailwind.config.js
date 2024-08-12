/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./public/**/*.{html,js}"],
    theme: {
        extend: {},
    },
    plugins: [],
};

// npx tailwindcss -i ./main.css -o ./dist/tailwind.css --watch
