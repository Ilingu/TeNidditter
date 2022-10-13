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

<details>
	<summary>
		<Link href={FormattedTedditUrl}>
			<p class="text-white inline-block w-32 truncate translate-y-1 hover:underline">
				{comment.link_title}
			</p>
		</Link>
		<FeedHeader
			author={comment.link_author}
			created={comment.created}
			subreddit={comment.subreddit}
		/>
	</summary>
	<div>{@html comment.body_html}</div>
</details>

<style scoped>
	summary > * {
		margin: -2px;
	}
	summary::marker {
		color: rgb(255 69 0);
	}
</style>
