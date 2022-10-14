<script lang="ts">
	import Link from "$lib/components/design/Link.svelte";
	import { humanElapsedTime, IsEmptyString } from "$lib/utils";

	export let subreddit: string = "";
	export let author: string = "";
	export let created: number;
</script>

{#if !IsEmptyString(subreddit)}
	<Link href={`/teddit/r/${subreddit}`} classStyle="text-white hover:underline font-bold"
		>r/{subreddit}</Link
	> <span class="mx-1">•</span>
{/if}
{#if !IsEmptyString(author)}
	<Link href={`/teddit/u/${author}`} classStyle="text-gray-300 hover:underline">u/{author}</Link>
	<span class="mx-1">•</span>
{/if}
<p class="text-gray-400">{humanElapsedTime(created * 1000, Date.now())}</p>
