import { render, screen } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
// @ts-ignore - Svelte component import
import Countdown from './Countdown.svelte';

describe('Countdown Component', () => {
	it('renders correctly with default values', () => {
		// Mock date to a specific point
		const futureDate = new Date(Date.now() + 120000).toISOString(); // 2 minutes from now
		render(Countdown, { deadline: futureDate });

		// It should display something like ⏱ 02:00
		const span = screen.getByText(/⏱/);
		expect(span).toBeInTheDocument();
	});

	it('fires onExpire when deadline is reached', () => {
		vi.useFakeTimers();
		const onExpire = vi.fn();
		const pastDate = new Date(Date.now() - 1000).toISOString(); // Already expired
		
		render(Countdown, { deadline: pastDate, onExpire });

		// Fast-forward 1 second so setInterval runs
		vi.advanceTimersByTime(1000);

		expect(onExpire).toHaveBeenCalled();
		vi.useRealTimers();
	});
});
