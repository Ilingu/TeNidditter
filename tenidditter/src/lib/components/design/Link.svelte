<script lang="ts">
	import { dispatchRouting } from "$lib/utils/routing";

	export let href: string;
	export let classStyle = "";
</script>

<a {href} class={classStyle} on:click={dispatchRouting}><slot /></a>
