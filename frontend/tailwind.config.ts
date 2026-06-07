import type { Config } from 'tailwindcss';
import typography from '@tailwindcss/typography';

// Palet warna mengikuti UI/UX Guideline pada PRD (Rose/Crimson akademis).
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
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
					soft: '#FDF2F4',
					muted: '#F9FAFB'
				},
				ink: {
					heading: '#1F2937',
					body: '#374151',
					caption: '#6B7280'
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
			}
		}
	},
	plugins: [typography]
} satisfies Config;
