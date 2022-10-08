<script lang="ts">
	import type { TedditPost } from "$lib/types/interfaces";

	import { fade } from "svelte/transition";
	import { afterUpdate, onMount } from "svelte";
	import Feed from "./Feed.svelte";
	import Loader from "$lib/components/design/Loader.svelte";

	export let loading = false;
	export let rawPosts: TedditPost[];

	type HandlePostParams = { afterId?: string; appendResult?: boolean };
	export let queryMorePostHandler: ({ afterId, appendResult }: HandlePostParams) => Promise<void>;

	let InfiniteScrollObserver: IntersectionObserver;
	onMount(() => {
		InfiniteScrollObserver = new IntersectionObserver(QueryMorePost);
	});

	afterUpdate(() => {
		InfiniteScrollObserver.disconnect();
		InfiniteScrollObserver.observe(document.getElementById(rawPosts.at(-5)?.id || "")!);
	});

	const QueryMorePost: IntersectionObserverCallback = async ([lastPost]) => {
		if (lastPost.isIntersecting) {
			InfiniteScrollObserver.unobserve(document.getElementById(rawPosts.at(-5)?.id || "")!);
			queryMorePostHandler({ afterId: rawPosts.at(-1)?.id, appendResult: true });
		}
	};
</script>

<div id="PostContainer" class="w-full relative" data-loading={loading}>
	{#each rawPosts as postData}
		<Feed post={postData} />
		<div class="my-5" />
	{/each}
	<div class="flex flex-col items-center">
		<button
			class="btn gap-x-2"
			on:click={() => queryMorePostHandler({ afterId: rawPosts.at(-1)?.id, appendResult: true })}
			>Next <i class="fa-solid fa-arrow-right icon" /></button
		>
	</div>

	{#if loading}
		<div class="flex justify-center">
			<div
				transition:fade
				class="absolute rounded-md -top-2 z-10 w-[105%] h-[102%] bg-[#343939] bg-opacity-30"
			/>
			<div class="fixed z-20 top-1/2 left-1/2 h-14 w-14 -translate-x-1/2 -translate-y-1/2">
				<Loader />
			</div>
		</div>
	{/if}
</div>

<style>
	:global(#PostContainer[data-loading="true"] .post) {
		opacity: 0.5;
	}
</style>
