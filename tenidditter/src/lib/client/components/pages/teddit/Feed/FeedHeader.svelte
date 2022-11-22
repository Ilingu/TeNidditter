<script lang="ts">
	import { humanElapsedTime } from "$lib/client/ClientUtils";
	import Link from "$lib/client/components/design/Link.svelte";
	import { IsEmptyString } from "$lib/shared/utils";

	export let subreddit: string = "";
	export let author: string = "";
	export let created: number;
</script>

{#if !IsEmptyString(subreddit)}
	<Link href={`/teddit/r/${subreddit}`}
		><p class="text-white hover:underline font-bold">r/{subreddit}</p>
	</Link> <span class="mx-1">•</span>
{/if}
{#if !IsEmptyString(author)}
	<Link href={`/teddit/u/${author}`}><p class="text-gray-300 hover:underline">u/{author}</p></Link>
	<span class="mx-1">•</span>
{/if}
<p class="text-gray-400">{humanElapsedTime(created * 1000, Date.now())}</p>
