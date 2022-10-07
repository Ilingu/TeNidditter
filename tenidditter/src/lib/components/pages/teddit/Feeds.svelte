<script lang="ts">
	import type { TedditPost } from "$lib/types/interfaces";
	import Feed from "./Feed.svelte";

	export let rawPosts: TedditPost[];
</script>

<div class="w-full">
	{#each rawPosts as postData}
		<Feed post={postData} />
		<div class="my-5" />
	{/each}
</div>
