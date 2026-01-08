/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                cyber: {
                    bg: '#14141e',
                    card: 'rgba(30, 30, 40, 0.6)',
                    border: 'rgba(255, 255, 255, 0.1)',
                    primary: '#00f2ff', // Neon Cyan
                    secondary: '#bd00ff', // Neon Purple
                    text: '#e0e0e0',
                }
            },
            backdropBlur: {
                xs: '2px',
            }
        },
    },
    plugins: [],
}
