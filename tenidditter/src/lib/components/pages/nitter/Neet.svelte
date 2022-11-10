<script lang="ts">
	import { HashDatas } from "$lib/services/wasm/encryption";
	import type { NeetComment } from "$lib/types/interfaces";
	import { FormatUsername } from "$lib/utils";

	export let neet: NeetComment;
</script>

<div
	id={FormatUsername(HashDatas(`${neet.creator}${neet.createdAt}`))}
	class="w-[500px] bg-primary-focus rounded-md"
>
	<div class="pp"><img src={neet.creator.avatarUrl} alt="User Profile Pic" /></div>
	<header>
		<!-- {#if neet}
			
		{/if} -->
	</header>

	<div class="body">{@html neet.content}</div>
	{#if neet.quote}
		<svelte:self neet={neet.quote} />
	{/if}
	<div class="attachments" />
	<div class="stats" />
</div>

<style scoped>
</style>
