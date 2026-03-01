/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主题色
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          200: '#bae6fd',
          300: '#7dd3fc',
          400: '#38bdf8',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
          800: '#075985',
          900: '#0c4a6e',
          950: '#082f49',
        },
        // 背景色
        background: 'var(--color-background)',
        'background-alt': 'var(--color-background-alt)',
        // 前景色
        foreground: 'var(--color-foreground)',
        'foreground-alt': 'var(--color-foreground-alt)',
        // 边框色
        border: 'var(--color-border)',
        // 静音色
        muted: 'var(--color-muted)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Menlo', 'monospace'],
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: '65ch',
            color: 'var(--color-foreground)',
            a: {
              color: 'var(--color-primary-500)',
              '&:hover': {
                color: 'var(--color-primary-600)',
              },
            },
            code: {
              color: 'var(--color-primary-500)',
              backgroundColor: 'var(--color-muted)',
              padding: '0.25rem 0.375rem',
              borderRadius: '0.25rem',
              fontWeight: '400',
            },
            'code::before': {
              content: '""',
            },
            'code::after': {
              content: '""',
            },
            pre: {
              backgroundColor: 'var(--color-background-alt)',
            },
            blockquote: {
              color: 'var(--color-foreground-alt)',
              borderLeftColor: 'var(--color-primary-500)',
            },
            hr: {
              borderColor: 'var(--color-border)',
            },
            strong: {
              color: 'var(--color-foreground)',
            },
            'ul li::marker': {
              color: 'var(--color-primary-500)',
            },
            'ol li::marker': {
              color: 'var(--color-primary-500)',
            },
          },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}