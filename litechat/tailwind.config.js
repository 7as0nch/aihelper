/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: '#0066ff', // Zhihu blue-ish
                secondary: '#f6f6f6',
            }
        },
    },
    plugins: [],
}
