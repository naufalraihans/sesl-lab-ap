<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { hasToken, user } from '$lib/stores/auth';
	import AppShell from '$lib/components/AppShell.svelte';

	import { initRunner } from '$lib/stores/runner';

	let { children } = $props();

	// Halaman publik (tanpa token): login, register, lupa/reset password.
	const PUBLIC_PATHS = [
		'/praktikum/login',
		'/praktikum/register',
		'/praktikum/lupa-password',
		'/praktikum/reset-password'
	];
	let isAuthPage = $derived(PUBLIC_PATHS.includes($page.url.pathname));
	let ready = $state(false);

	onMount(() => {
		if (!isAuthPage && !hasToken()) {
			goto('/praktikum/login');
			return;
		}
		initRunner();
		ready = true;
	});
</script>

{#if isAuthPage}
	{@render children()}
{:else if ready}
	<AppShell>
		{@render children()}
	</AppShell>
{/if}

