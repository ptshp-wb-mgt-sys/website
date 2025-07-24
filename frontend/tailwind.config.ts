import type { Config } from 'tailwindcss'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        aquamarine: {
          DEFAULT: '#6cf5bc',
          100: '#054228',
          200: '#098451',
          300: '#0ec579',
          400: '#29f09d',
          500: '#6cf5bc',
          600: '#89f7c9',
          700: '#a6f9d6',
          800: '#c4fbe4',
          900: '#e1fdf1'
        },
        'rich-black': {
          DEFAULT: '#001011',
          100: '#000303',
          200: '#000606',
          300: '#000909',
          400: '#000b0c',
          500: '#001011',
          600: '#006b72',
          700: '#00c7d5',
          800: '#39f2ff',
          900: '#9cf8ff'
        },
        gray: {
          DEFAULT: '#757780',
          100: '#181819',
          200: '#2f3033',
          300: '#47484c',
          400: '#5e5f66',
          500: '#757780',
          600: '#919299',
          700: '#acadb3',
          800: '#c8c9cc',
          900: '#e3e4e6'
        },
        seasalt: {
          DEFAULT: '#f8f8f8',
          100: '#313131',
          200: '#636363',
          300: '#949494',
          400: '#c6c6c6',
          500: '#f8f8f8',
          600: '#f9f9f9',
          700: '#fafafa',
          800: '#fcfcfc',
          900: '#fdfdfd'
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
} satisfies Config 
