import type { Config } from 'tailwindcss';
import typography from '@tailwindcss/typography';

// Palet warna mengikuti UI/UX Guideline pada PRD (Rose/Crimson akademis).
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
				display: ['Inter', '"Plus Jakarta Sans"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
				mono: ['"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'monospace']
			},
			colors: {
				primary: {
					DEFAULT: '#8A1538',
					hover: '#B21F47',
					active: '#5C0E25',
					dark: '#5C0E25'
				},
				navbar: '#8A1538',
				sidebar: '#5C0E25',
				surface: {
					base: '#F8FAFC',
					soft: '#F1F5F9',
					muted: '#F8FAFC'
				},
				ink: {
					heading: '#1E293B',
					body: '#475569',
					caption: '#64748B'
				},
				state: {
					success: '#06D6A0',
					'success-bg': '#D1FAE5',
					warning: '#FFC300',
					'warning-bg': '#FEF3C7',
					error: '#EF4444',
					'error-bg': '#FEE2E2',
					info: '#48CAE4',
					'info-bg': '#DBEAFE'
				},
				'neo-maroon': '#5C0E25',
				'maroon-bg': '#FDF2F4',
				'fun-yellow': '#FFC300',
				'fun-blue': '#48CAE4',
				'fun-green': '#06D6A0',
				'fun-purple': '#9D4EDD',
				'fun-pink': '#F472B6'
			},
			borderRadius: {
				DEFAULT: '0.5rem'
			},
			transitionTimingFunction: {
				spring: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
				'smooth-out': 'cubic-bezier(0.16, 1, 0.3, 1)'
			},
			keyframes: {
				'fade-in-up': {
					'0%': { opacity: '0', transform: 'translateY(20px)' },
					'100%': { opacity: '1', transform: 'translateY(0)' }
				},
				'fade-in': {
					'0%': { opacity: '0' },
					'100%': { opacity: '1' }
				},
				float: {
					'0%, 100%': { transform: 'translateY(0)' },
					'50%': { transform: 'translateY(-10px)' }
				},
				'pulse-glow': {
					'0%, 100%': { opacity: '1', boxShadow: '0 0 0 0 rgba(208, 49, 83, 0.4)' },
					'50%': { opacity: '0.9', boxShadow: '0 0 20px 0 rgba(208, 49, 83, 0.6)' }
				},
				'gradient-shift': {
					'0%': { backgroundPosition: '0% 50%' },
					'50%': { backgroundPosition: '100% 50%' },
					'100%': { backgroundPosition: '0% 50%' }
				}
			},
			animation: {
				'fade-in-up': 'fade-in-up 0.6s ease-out forwards',
				'fade-in': 'fade-in 0.5s ease-out forwards',
				float: 'float 4s ease-in-out infinite',
				'pulse-glow': 'pulse-glow 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
				'gradient-shift': 'gradient-shift 8s ease infinite'
			}
		}
	},
	plugins: [typography]
} satisfies Config;
