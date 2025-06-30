/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // Enhanced Social Pay Brand Colors (based on logo analysis)
        brand: {
          green: {
            50: '#f0fdf4',
            100: '#dcfce7',
            200: '#bbf7d0',
            300: '#86efac',
            400: '#4ade80',
            500: '#22c55e', // Primary brand green
            600: '#16a34a',
            700: '#15803d',
            800: '#166534',
            900: '#14532d',
            950: '#052e16',
          },
          gold: {
            50: '#fffbeb',
            100: '#fef3c7',
            200: '#fde68a',
            300: '#fcd34d',
            400: '#fbbf24', // Primary brand gold
            500: '#f59e0b',
            600: '#d97706',
            700: '#b45309',
            800: '#92400e',
            900: '#78350f',
            950: '#451a03',
          },
          // Additional accent colors for variety
          teal: {
            50: '#f0fdfa',
            100: '#ccfbf1',
            200: '#99f6e4',
            300: '#5eead4',
            400: '#2dd4bf',
            500: '#14b8a6',
            600: '#0d9488',
            700: '#0f766e',
            800: '#115e59',
            900: '#134e4a',
          },
          orange: {
            50: '#fff7ed',
            100: '#ffedd5',
            200: '#fed7aa',
            300: '#fdba74',
            400: '#fb923c',
            500: '#f97316',
            600: '#ea580c',
            700: '#c2410c',
            800: '#9a3412',
            900: '#7c2d12',
          }
        },
        // Maintain backward compatibility
        primary: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
        },
        secondary: {
          50: '#fffbeb',
          100: '#fef3c7',
          200: '#fde68a',
          300: '#fcd34d',
          400: '#fbbf24',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
          800: '#92400e',
          900: '#78350f',
        },
        // Additional utility colors
        success: {
          50: '#f0fdf4',
          500: '#22c55e',
          600: '#16a34a',
        },
        warning: {
          50: '#fffbeb',
          500: '#fbbf24',
          600: '#d97706',
        },
        danger: {
          50: '#fef2f2',
          500: '#ef4444',
          600: '#dc2626',
        },
        info: {
          50: '#f0f9ff',
          500: '#06b6d4',
          600: '#0891b2',
        }
      },
      backgroundImage: {
        // Enhanced gradients
        'brand-gradient': 'linear-gradient(135deg, #22c55e 0%, #fbbf24 100%)',
        'brand-gradient-reverse': 'linear-gradient(135deg, #fbbf24 0%, #22c55e 100%)',
        'brand-gradient-soft': 'linear-gradient(135deg, #22c55e 0%, #86efac 50%, #fbbf24 100%)',
        'hero-gradient': 'linear-gradient(135deg, #f0fdf4 0%, #ffffff 25%, #fffbeb 100%)',
        'card-gradient': 'linear-gradient(135deg, rgba(255,255,255,0.9) 0%, rgba(255,255,255,0.7) 100%)',
        // Animated gradients
        'animated-gradient': 'linear-gradient(-45deg, #22c55e, #16a34a, #fbbf24, #f59e0b)',
      },
      boxShadow: {
        // Enhanced shadows with brand colors
        'brand': '0 4px 14px 0 rgba(34, 197, 94, 0.25)',
        'brand-lg': '0 10px 25px -3px rgba(34, 197, 94, 0.25), 0 4px 6px -2px rgba(34, 197, 94, 0.05)',
        'brand-xl': '0 20px 50px -12px rgba(34, 197, 94, 0.30)',
        'gold': '0 4px 14px 0 rgba(251, 191, 36, 0.25)',
        'gold-lg': '0 10px 25px -3px rgba(251, 191, 36, 0.25), 0 4px 6px -2px rgba(251, 191, 36, 0.05)',
        'combined': '0 4px 20px rgba(34, 197, 94, 0.15), 0 8px 40px rgba(251, 191, 36, 0.10)',
        'glass': '0 8px 32px 0 rgba(31, 38, 135, 0.37)',
        'inner-brand': 'inset 0 2px 4px 0 rgba(34, 197, 94, 0.06)',
      },
      animation: {
        // Enhanced animations
        'blob': 'blob 7s infinite',
        'float': 'float 6s ease-in-out infinite',
        'pulse-green': 'pulse-green 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'slide-in-right': 'slideInFromRight 0.6s ease-out',
        'slide-in-left': 'slideInFromLeft 0.6s ease-out',
        'fade-in-up': 'fadeInUp 0.6s ease-out',
        'gradient': 'gradient 15s ease infinite',
        'bounce-subtle': 'bounceSubtle 2s infinite',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Inter', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        'xl': '0.75rem',
        '2xl': '1rem',
        '3xl': '1.5rem',
      },
      spacing: {
        '18': '4.5rem',
        '72': '18rem',
        '84': '21rem',
        '96': '24rem',
      },
      backdropBlur: {
        'xs': '2px',
      },
      scale: {
        '102': '1.02',
      }
    },
  },
  plugins: [],
}

