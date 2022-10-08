<script lang="ts">
	import type { TedditPost } from "$lib/types/interfaces";
	import { onMount } from "svelte";
	import Feed from "./Feed.svelte";

	export let rawPosts: TedditPost[];
	export let queryMorePostHandler: (afterId?: string) => void;

	let InfiniteScrollObserver: IntersectionObserver;
	onMount(() => {
		InfiniteScrollObserver = new IntersectionObserver(QueryMorePost);
		InfiniteScrollObserver.observe(document.getElementById(rawPosts.at(-5)?.id || "")!);
	});

	const QueryMorePost: IntersectionObserverCallback = async ([lastPost]) => {
		if (lastPost.isIntersecting) {
			InfiniteScrollObserver.unobserve(document.getElementById(rawPosts.at(-5)?.id || "")!);
			queryMorePostHandler(rawPosts.at(-1)?.id);
		}
	};
</script>

<div class="w-full">
	{#each rawPosts as postData}
		<Feed post={postData} />
		<div class="my-5" />
	{/each}
	<div class="flex flex-col items-center">
		<button class="btn gap-x-2">Next <i class="fa-solid fa-arrow-right icon" /></button>
		<!-- <button class="btn">Prev</button> -->
	</div>
</div>
