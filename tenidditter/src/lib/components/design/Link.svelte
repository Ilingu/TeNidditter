<script lang="ts">
	import { dispatchRouting } from "$lib/routing";

	export let href: string;
	export let classStyle = "";
</script>

<a data-sveltekit-prefetch {href} class={classStyle} on:click={dispatchRouting}><slot /></a>
