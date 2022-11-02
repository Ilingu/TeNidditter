<script lang="ts">
	import { dispatchRouting } from "$lib/routing";

	export let href: string;
	export let className = "";
</script>

<a data-sveltekit-prefetch {href} class={className} on:click={dispatchRouting}><slot /></a>
