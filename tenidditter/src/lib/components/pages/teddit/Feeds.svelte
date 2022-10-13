<script lang="ts">
	import type { TedditCommmentShape, TedditPost } from "$lib/types/interfaces";

	import { fade } from "svelte/transition";
	import { afterUpdate, onMount } from "svelte";
	import Feed from "./Feed/Feed.svelte";
	import Loader from "$lib/components/design/Loader.svelte";
	import { isValidUrl } from "$lib/utils";
	import Comments from "./Comments.svelte";

	export let loading = false;
	export let rawPosts: (TedditPost | TedditCommmentShape)[];

	type HandlePostParams = { afterId?: string; appendResult?: boolean };
	export let queryMorePostHandler: ({
		afterId,
		appendResult
	}: HandlePostParams) => Promise<void> = async () => {};

	let InfiniteScrollObserver: IntersectionObserver;
	onMount(() => {
		InfiniteScrollObserver = new IntersectionObserver(ObserverHandler);

		document.querySelectorAll(".md a").forEach((a) => {
			const href = a.getAttribute("href");
			if (href?.startsWith("/r/")) a.setAttribute("href", `/teddit${href}`);
			else {
				if (isValidUrl(a.textContent || "")) a.textContent = new URL(a.textContent || "").host;
				a.innerHTML = `<i class="fa-solid fa-arrow-up-right-from-square"></i>` + a.innerHTML;

				a.setAttribute("target", "_blank");
				a.setAttribute("rel", "noopener noreferrer");
			}
		});
	});

	const ObserverId = (from: number) =>
		new Promise<number>((res, rej) => {
			let i = from;

			const ChooseAgain = () => {
				if (i >= 25) return rej();
				if (Object.hasOwn(rawPosts.at(-i) || {}, "id")) return res(i);

				i++;
				ChooseAgain();
			};
			ChooseAgain();
		});

	afterUpdate(async () => {
		try {
			InfiniteScrollObserver.disconnect();
			InfiniteScrollObserver.observe(
				document.getElementById((rawPosts.at(-(await ObserverId(5))) as any)?.id || "")!
			);
		} catch (err) {}
	});

	const ObserverHandler: IntersectionObserverCallback = async ([lastPost]) => {
		if (!lastPost.isIntersecting) return;
		QueryMorePost();
	};

	const QueryMorePost = async () => {
		try {
			InfiniteScrollObserver.unobserve(
				document.getElementById((rawPosts.at(-(await ObserverId(5))) as any)?.id || "")!
			);
			queryMorePostHandler({
				afterId: (rawPosts.at(-(await ObserverId(1))) as any)?.id,
				appendResult: true
			});
		} catch (err) {}
	};
</script>

<div id="PostContainer" class="w-full relative" data-loading={loading}>
	{#each rawPosts as postData}
		{#if postData.type === "t1"}
			<Comments comment={postData} />
		{:else}
			<Feed post={postData} />
		{/if}
		<div class="my-5" />
	{/each}
	<div class="flex flex-col items-center">
		<button class="btn gap-x-2" on:click={QueryMorePost}
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
