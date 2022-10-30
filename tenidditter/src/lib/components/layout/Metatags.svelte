<script lang="ts">
	import { onMount } from "svelte";

	export let title = "TeNidditter";
	export let description = "🔒 TeNidditter - Secure Nitter and Teddit Feed Generator";
	export let icon = "/favicon.ico";
	export let routePath = false;

	onMount(() => {
		if (routePath) title += location.pathname;
	});
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<link rel="icon" href={icon} />
</svelte:head>
