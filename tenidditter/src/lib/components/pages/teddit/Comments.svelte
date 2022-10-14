<script lang="ts">
	import Link from "$lib/components/design/Link.svelte";
	import type { TedditCommmentShape } from "$lib/types/interfaces";
	import FeedHeader from "./Feed/FeedHeader.svelte";

	export let comment: TedditCommmentShape;

	const FormattedTedditUrl =
		"/teddit" +
		comment.permalink
			.replace(/\/+$/, "")
			.replace("comments/", "")
			.split("/")
			.slice(0, -2)
			.join("/");
</script>

<div class="flex  gap-x-2">
	<div
		class="min-w-[32px] w-8 h-8 bg-teddit rounded text-black font-bold font-nerd flex justify-center items-center"
		title="⬆ Ups"
	>
		{comment.ups}
	</div>

	<details class="translate-y-1">
		<summary>
			<div class="inline-block">
				<div class="flex text-sm">
					<Link href={FormattedTedditUrl}>
						<p class="truncate-word text-white hover:underline" title={comment.link_title}>
							{comment.link_title}
						</p>
					</Link> •
					<FeedHeader
						author={comment.link_author}
						created={comment.created}
						subreddit={comment.subreddit}
					/>
				</div>
			</div>
		</summary>
		<div class="bg-primary-content color-white rounded-lg p-5">{@html comment.body_html}</div>
	</details>
</div>

<style scoped>
	summary > * {
		margin: -2px;
	}
	summary::marker {
		color: rgb(255 69 0);
	}

	.truncate-word {
		max-width: 15rem;
		display: -webkit-box;
		-webkit-line-clamp: 1;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
