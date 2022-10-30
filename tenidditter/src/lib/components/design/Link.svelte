<script lang="ts">
	import { dispatchRouting } from "$lib/routing";

	export let href: string;
</script>

<a data-sveltekit-prefetch {href} on:click={dispatchRouting}><slot /></a>
