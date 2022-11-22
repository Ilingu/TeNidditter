<script lang="ts">
	export let elems: string[];
	export let active: number;

	export let cb: (active: number) => void;
</script>

<div class="tabs-component tabs justify-center tabs-boxed tab-lg">
	{#each elems as elem, i}
		<p class={`tab tab-bordered  ${i === active ? "tab-active" : ""}`} on:click={() => cb(i)}>
			{@html elem}
		</p>
	{/each}
</div>
