<script lang="ts">
	import type { TedditCommmentShape, TedditPost } from "$lib/types/interfaces";

	import { fade } from "svelte/transition";
	import { onMount } from "svelte";
	import Feed from "./Feed/Feed.svelte";
	import Loader from "$lib/components/design/Loader.svelte";
	import { isValidUrl } from "$lib/utils";
	import Comments from "./comments/Comments.svelte";

	export let loading = false;
	export let rawPosts: (TedditPost | TedditCommmentShape)[];

	onMount(() => {
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
