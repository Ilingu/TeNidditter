<script lang="ts">
	import Feed from "$lib/components/pages/teddit/Feed/Feed.svelte";
	import { page } from "$app/stores";
	import { setContext } from "svelte";
	import Comments from "$lib/components/pages/teddit/comments/Comments.svelte";

	export let data: import("./$types").PageData;
	console.log(data);

	setContext("COMMEMTS_CTX", data.comments);
	const PostUps = (): number => {
		const ups = data.metadata.post_ups;
		if (ups.endsWith("k")) {
			const [num] = ups.split("k");
			const safeNum = Math.round(parseFloat(num));
			return parseInt(`${safeNum}000`);
		}
		return parseInt(ups);
	};
</script>

<main class="max-w-[1500px] m-auto flex flex-col items-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feed
			post={{
				id: $page.params.postid,
				title: data.metadata.post_title,
				author: data.metadata.post_author,
				created: data.metadata.post_created,

				ups: PostUps(),
				num_comments: data.metadata.post_nb_comments,

				subreddit: $page.params.subteddit,

				is_self_link: true,
				is_video: false,
				stickied: false,

				url: "",
				domain: "",
				permalink: ""
			}}
		/>

		<section class="mt-5">
			{#each data.comments as co, i}
				<Comments idxCtx={i} comment={co[0]} open recursive />
				<div class="my-10" />
			{/each}
		</section>
	</div>
</main>
