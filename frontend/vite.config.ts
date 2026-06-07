import { sveltekit } from '@sveltejs/kit/vite';
import { svelteTesting } from '@testing-library/svelte/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit(), svelteTesting()],
	server: {
		port: 5173
	},
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		globals: true,
		css: false,
		setupFiles: ['./vitest-setup.ts'],
		server: {
			deps: {
				inline: [/katex/]
			}
		}
	}
});
