import type { Config } from 'tailwindcss';
import typography from '@tailwindcss/typography';

// Palet warna mengikuti UI/UX Guideline pada PRD (Rose/Crimson akademis).
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
				display: ['"Plus Jakarta Sans"', 'Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
				mono: ['"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'monospace']
			},
			colors: {
				primary: {
					DEFAULT: '#D03153',
					hover: '#E03A60',
					active: '#B02A47',
					dark: '#942240'
				},
				navbar: '#B02A47',
				sidebar: '#942240',
				surface: {
					base: '#FDFCFB',
					soft: '#FDF2F4',
					muted: '#FAF9F8'
				},
				ink: {
					heading: '#1C1917',
					body: '#44403C',
					caption: '#78716C'
				},
				state: {
					success: '#10B981',
					'success-bg': '#D1FAE5',
					warning: '#F59E0B',
					'warning-bg': '#FEF3C7',
					error: '#EF4444',
					'error-bg': '#FEE2E2',
					info: '#3B82F6',
					'info-bg': '#DBEAFE'
				}
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
