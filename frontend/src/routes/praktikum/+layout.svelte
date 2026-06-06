<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { hasToken, user } from '$lib/stores/auth';
	import AppShell from '$lib/components/AppShell.svelte';

	let { children } = $props();

	let isLogin = $derived($page.url.pathname === '/praktikum/login');
	let ready = $state(false);

	onMount(() => {
		if (!isLogin && !hasToken()) {
			goto('/praktikum/login');
			return;
		}
		ready = true;
	});
</script>

{#if isLogin}
	{@render children()}
{:else if ready}
	<AppShell>
		{@render children()}
	</AppShell>
{/if}

{#if !isLogin && ready && $user === null}
	<!-- token ada tapi user store kosong: tetap render shell -->
{/if}
