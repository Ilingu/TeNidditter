<script lang="ts">
	import type { TedditPost } from "$lib/types/interfaces";

	export let post: TedditPost;
</script>

<div
	id={post.id}
	class="post w-[750px] p-2 bg-primary-content min-h-[128px] rounded-lg ring-1 ring-[#686868] "
>
	<score class="post-score" />
	<header class="post-header" />
	<div class="post-title" />
	<div class="post-body" />
	<div class="post-footer" />
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
		grid-template-rows: 0.225fr 0.4fr 1fr 1fr 0.225fr;
		grid-template-columns: 50px 1fr;
	}
	:global(.post.nobody) {
		grid-template-rows: 1fr 1fr 0fr 0fr 1fr;
	}

	.post-score {
		grid-area: post-score;
	}
	.post-header {
		grid-area: post-header;
	}
	.post-title {
		grid-area: post-title;
	}
	.post-body {
		grid-area: post-body;
	}
	.post-footer {
		grid-area: post-footer;
	}
</style>
