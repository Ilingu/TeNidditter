<script lang="ts">
	import type { TedditPost } from "$lib/types/interfaces";

	import Link from "$lib/components/design/Link.svelte";
	import { FormatNumbers, humanElapsedTime } from "$lib/utils";

	export let post: TedditPost;

	const FormattedTedditUrl = post.permalink
		.replace(/\/+$/, "")
		?.replace("https://teddit.net", "/teddit")
		.replace("comments/", "")
		.split("/")
		.slice(0, -1)
		.join("/");
</script>

<div
	id={post.id}
	class={`post w-[750px] px-1 pt-4 pb-1 transition-all bg-primary-content hover:bg-[rgba(45,49,49,0.71)] min-h-[128px] rounded-lg ring-1 ring-[#686868] ${
		post.stickied === true ? "stickied" : ""
	}`}
>
	<score class="post-score text-sm text-teddit text-center"
		><i class="fa-solid fa-arrow-up mr-1" /> <br />{FormatNumbers(post.ups)}</score
	>
	<header class="post-header flex text-sm">
		<Link href={`/teddit/r/${post.subreddit}`} classStyle="text-white hover:underline font-bold"
			>r/{post.subreddit}</Link
		> <span class="mx-1">•</span>
		<Link href={`/teddit/u/${post.author}`} classStyle="text-gray-300 hover:underline"
			>u/{post.author}</Link
		> <span class="mx-1">•</span>
		<p class="text-gray-400">{humanElapsedTime(post.created * 1000, Date.now())}</p>
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
		<Link href={FormattedTedditUrl} classStyle="hover:underline hover:text-white transition-all"
			>{post.title}</Link
		>
	</div>
	<div class="post-body flex flex-col items-center">
		{#if post.selftext_html}
			{@html post.selftext_html}
		{/if}

		{#if post.is_video && post.media}
			<video
				class="my-5 max-h-[512px]"
				width={post.media?.reddit_video?.width}
				src={post.media?.reddit_video?.fallback_url}
				controls
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
		{FormatNumbers(post.num_comments)} comments
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
		max-height: 300px;
		overflow: hidden;
		position: relative;
		-webkit-mask-image: linear-gradient(180deg, black, 75%, transparent);
		mask-image: linear-gradient(180deg, black, 75%, transparent);
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
	:global(.md code) {
		background: #666;
		border-left: 3px solid #ff4500;
		color: #fff;
		page-break-inside: avoid;
		font-family: monospace;
		font-size: 15px;
		line-height: 1.6;
		margin-bottom: 1.6em;
		max-width: 100%;
		overflow: auto;
		padding: 1em 1.5em;
		display: block;
		word-wrap: break-word;
	}
</style>
