<script lang="ts">
	import type { TedditPost } from "$lib/types/interfaces";

	import Link from "$lib/components/design/Link.svelte";
	import { FormatNumbers } from "$lib/utils";
	import FeedHeader from "./FeedHeader.svelte";
	import { onMount } from "svelte";

	export let post: TedditPost & { body_html?: string };
	export let blur = true;

	const FormattedTedditUrl = post.permalink
		.replace(/\/+$/, "")
		?.replace("https://teddit.net", "/teddit")
		.replace("comments/", "")
		.split("/")
		.slice(0, -1)
		.join("/");

	let PostElem: HTMLDivElement;
	onMount(async () => {
		if (!PostElem) return;
		const codeBlock = PostElem.querySelectorAll(".md pre code");
		if (codeBlock.length <= 0) return;

		// code Highlighing
		const hljs = (await import("highlight.js")).default;
		document
			.querySelectorAll(".md pre code")
			.forEach((el) => hljs.highlightElement(el as HTMLElement));
	});
</script>

<div
	id={post.id}
	bind:this={PostElem}
	class={`post ${
		blur ? "blurMask" : ""
	} md:w-[750px] w-[92.5vw] px-1 pt-4 pb-1 transition-all bg-primary-content hover:bg-light-dark min-h-[128px] rounded-lg ring-1 ring-[#686868] ${
		post.stickied === true ? "stickied" : ""
	}`}
>
	<score class="md:block hidden post-score text-sm text-teddit text-center"
		><i class="fa-solid fa-arrow-up mr-1" /> <br />{FormatNumbers(post.ups)}</score
	>
	<header class="post-header flex flex-wrap text-sm">
		<FeedHeader author={post.author} subreddit={post.subreddit} created={post.created} />
	</header>
	<div class="post-title mb-2 text-gray-50">
		{#if post.link_flair}
			<span class="badge badge-warning"
				>{@html post.link_flair.replaceAll(
					"/pics/flairs/",
					"https://teddit.pussthecat.org/pics/flairs/"
				)}</span
			>
		{/if}
		<Link href={FormattedTedditUrl}
			><p class="hover:underline hover:text-white transition-all">{post.title}</p></Link
		>
	</div>
	<div class="post-body flex flex-col items-center">
		{#if post.body_html}
			{@html post.body_html
				.replaceAll("/vids", "https://teddit.net/vids")
				.replaceAll("/pics", "https://teddit.net/pics")}
		{/if}

		{#if post.selftext_html}
			{@html post.selftext_html}
		{/if}

		{#if post.is_video && post.media}
			<video
				class="my-5 max-h-[512px] w-auto"
				width={post.media?.reddit_video?.width}
				src={post.media?.reddit_video?.fallback_url}
				controls
				preload="none"
			>
				<track kind="captions" /></video
			>
		{/if}

		{#if post.images && post.images.preview && post.is_self_link}
			<img
				src={"https://teddit.pussthecat.org" + post.images.preview}
				class="my-5 max-h-[512px] w-auto"
				alt="🖼"
			/>
		{/if}

		{#if !post.is_self_link}
			<div class="p-2 text-sm flex flex-col items-center bg-black rounded-md ring-1 ring-[#686868]">
				<a href={post.url} target="_blank" rel="noopener noreferrer">
					{#if post.images?.thumb}
						<img src={post.images?.thumb} class="my-5 w-full" alt="🔗" />
					{/if}
					<legend class="flex items-center gap-x-2 hover:underline"
						><i class="fa-solid fa-arrow-up-right-from-square" /> {post.domain}</legend
					>
				</a>
			</div>
		{/if}
	</div>
	<div class="post-footer text-sm flex items-end">
		<p class="md:hidden block post-score text-sm text-teddit text-center">
			{FormatNumbers(post.ups)}
		</p>
		<span class="mx-1">•</span>{FormatNumbers(post.num_comments)} comments
	</div>
</div>

<style scoped>
	.post {
		display: grid;

		grid-template-areas:
			"post-score post-header post-header"
			"post-score post-title post-title"
			"post-score post-body post-body"
			"post-score post-body post-body"
			"post-score post-footer post-footer";
		grid-template-rows: auto;
		/* grid-template-rows: 0.225fr 0.4fr minmax(0, 1fr) 1fr 0.225fr; */
		grid-template-columns: 45px 1fr;
	}
	@media (max-width: 768px) {
		.post {
			grid-template-areas:
				"post-header post-header"
				"post-title post-title"
				"post-body post-body"
				"post-body post-body"
				"post-footer post-footer";
		}
	}

	.post.stickied {
		border: 2px solid #ff4500;
		box-shadow: none;
	}

	.post-score {
		grid-area: post-score;
	}

	.post-header {
		grid-area: post-header;
	}

	:global(.post-title .emoji) {
		width: 1.25em;
		height: 1.25em;
		display: inline-block;
		background-size: contain;
		background-position: 50% 50%;
		background-repeat: no-repeat;
		vertical-align: middle;
	}

	.post-title {
		grid-area: post-title;
	}

	.post-body {
		grid-area: post-body;
	}
	.post-body > * {
		margin-right: 40px;
	}

	.post-footer {
		grid-area: post-footer;
	}

	:global(.post .md) {
		text-align: justify;
		width: 100%;
		overflow: hidden;
		position: relative;
	}
	:global(.post.blurMask .md) {
		-webkit-mask-image: linear-gradient(180deg, black, 75%, transparent);
		mask-image: linear-gradient(180deg, black, 75%, transparent);
		max-height: 300px;
	}

	:global(.md > *) {
		margin-bottom: 12.5px;
	}
	:global(.md a) {
		font-style: italic;
		color: #92bddf;
		transition: all;
	}
	:global(.md a:hover) {
		text-decoration: underline wavy;
		color: #5296dd;
	}

	:global(.md :not(pre) code) {
		background: #666;
		border-left: 3px solid #ff4500;
		color: #fff;
		font-family: monospace;
		font-size: 15px;
		overflow: auto;
		display: inline;
		word-wrap: break-word;
	}
</style>
