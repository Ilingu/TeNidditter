<script lang="ts">
	import Link from "$lib/client/components/design/Link.svelte";
	import FeedHeader from "../Feed/FeedHeader.svelte";
	import { page } from "$app/stores";
	import { getContext } from "svelte";
	import { onMount } from "svelte";
	import type { TedditCommmentShape } from "$lib/shared/types/teddit";
	import { isMobile } from "$lib/client/ClientUtils";
	import { IsEmptyString } from "$lib/shared/utils";

	type Comments = TedditCommmentShape & { id?: number; parentId?: number };
	export let comment: Comments;

	export let open = false;

	export let recursive = false;
	export let idxCtx: number = 0;

	const comments = ((getContext("COMMEMTS_CTX") as Comments[][]) || [""])[idxCtx];
	let childrenComments: Comments[] = [];
	onMount(() => {
		if (!recursive) return; // to the potential contributor, I fucked up this components but it works as following: if the "recursive" props is set to true your datas must have the "id" and "parentId" field
		childrenComments = comments?.filter(({ parentId }) => parentId === comment.id) ?? [];
		if (comment.id && comment.id >= 5) open = false;
	});

	const FormattedTedditUrl =
		comment?.permalink &&
		"/teddit" +
			comment?.permalink
				.replace(/\/+$/, "")
				.replace("comments/", "")
				.split("/")
				.slice(0, -2)
				.join("/");
</script>

<div class="flex justify-center gap-x-2 md:w-[750px] w-[92.5vw] relative">
	<div
		class="littleBar lg:block hidden absolute bottom-0 left-3 translate-x-0.5 rounded-2xl w-1 bg-accent"
	/>
	<div
		class="min-w-[32px] lg:flex hidden min-w-8 h-8 bg-teddit rounded font-fancy text-lg  justify-center items-center"
		title="⬆ Ups"
	>
		{`${comment?.ups}`.replace(/\D/g, "") || "none"}
	</div>

	<details class={`translate-y-1 w-full ${isMobile() && !open ? "ml-5" : ""}`} {open}>
		<summary>
			<div class="inline-block">
				<div class="flex flex-wrap text-sm">
					<p class="mr-2 lg:hidden text-teddit font-fancy -translate-y-0.5" title="⬆ Ups">
						⬆ {`${comment?.ups}`.replace(/\D/g, "") || "none"}
					</p>
					{#if !IsEmptyString(comment?.link_title)}
						<Link href={FormattedTedditUrl || ""}>
							<p class="truncate-word text-white hover:underline" title={comment?.link_title}>
								{comment?.link_title || ""}
							</p>
						</Link> •
					{/if}
					<FeedHeader
						author={comment?.link_author}
						created={comment?.created || Date.now() / 1000}
						subreddit={comment?.subreddit || $page.params.subteddit}
					/>
				</div>
			</div>
		</summary>
		<div
			class={`bg-primary-content rounded-lg px-5 py-2 text-justify ${
				comment?.id && comment.id % 2 === 0 ? "bg-[#333] text-white" : "bg-transparent"
			}`}
		>
			{@html comment?.body_html}
		</div>
		{#each childrenComments as nextComment}
			<div class="my-2" />
			<svelte:self {idxCtx} comment={nextComment} open={!isMobile()} recursive />
		{/each}
	</details>
</div>

<style scoped>
	summary > * {
		margin: -2px;
	}
	summary::marker {
		color: rgb(255, 69, 0);
	}

	.littleBar {
		height: calc(100% - 32px);
	}

	.truncate-word {
		max-width: 15rem;
		display: -webkit-box;
		-webkit-line-clamp: 1;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
