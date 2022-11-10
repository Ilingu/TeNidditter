<script lang="ts">
	import type { NeetComment } from "$lib/types/interfaces";
	import { onMount } from "svelte";
	import Neet from "./Neet.svelte";

	export let neets: NeetComment[][];

	let feedDiv: HTMLDivElement;
	onMount(() => {
		if (!feedDiv) return;

		// feedDiv.querySelectorAll("");
	});
</script>

<div class="flex flex-col gap-y-5" bind:this={feedDiv}>
	{#each neets as neetThread}
		<div class="flex flex-col gap-y-2">
			{#each neetThread as neet}
				<Neet {neet} />
			{/each}
		</div>
	{/each}
</div>
