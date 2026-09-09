<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { hasToken, user } from '$lib/stores/auth';
	import AppShell from '$lib/components/AppShell.svelte';

	import { initRunner } from '$lib/stores/runner';

	let { children } = $props();

	// Login & register = halaman publik, lewati guard token.
	let isAuthPage = $derived(
		$page.url.pathname === '/praktikum/login' || $page.url.pathname === '/praktikum/register'
	);
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

{#if !isLogin && ready && $user === null}
	<!-- token ada tapi user store kosong: tetap render shell -->
{/if}
