<script lang="ts">
	import type { TedditCommmentShape, TedditPost } from "$lib/types/interfaces";

	import { fade } from "svelte/transition";
	import { onMount } from "svelte";
	import Feed from "./Feed/Feed.svelte";
	import Loader from "$lib/components/design/Loader.svelte";
	import { EscapeHTML, isValidUrl } from "$lib/utils";
	import Comments from "./comments/Comments.svelte";

	export let loading = false;
	export let rawPosts: ((TedditPost & { blur?: boolean }) | TedditCommmentShape)[];

	onMount(async () => {
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

		const codeBlock = document.querySelectorAll("#PostContainer .md pre code");
		if (codeBlock.length <= 0) return;

		// code Highlighing (maybe webworker if too much)
		const hljs = (await import("highlight.js")).default;
		codeBlock.forEach((el) => {
			const safeHtml = EscapeHTML(el.innerHTML);
			el.innerHTML = safeHtml; // XSS protection

			hljs.highlightElement(el as HTMLElement);
		});
	});
</script>

<div id="PostContainer" class="w-full grid place-items-center relative" data-loading={loading}>
	{#each rawPosts as postData}
		{#if postData.type === "t1"}
			<Comments comment={postData} />
		{:else}
			<Feed post={postData} blur={postData.blur} />
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
