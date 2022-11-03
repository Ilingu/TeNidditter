<!-- Feed -->
<script lang="ts">
	import AlertChild from "$lib/components/design/Alert/AlertChild.svelte";
	import Link from "$lib/components/design/Link.svelte";
	import Tabs from "$lib/components/design/Tabs.svelte";
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";
	import { QueryHomePost } from "$lib/services/teddit";
	import { FeedTypeEnum } from "$lib/types/enums";
	import { afterUpdate, onMount } from "svelte";

	export let data: import("./$types").PageData;

	let FeedDisplayType = FeedTypeEnum.Hot;
	let loading = false;

	let minusFifthPostId = data.data?.at(-5)?.id ?? "";
	$: minusFifthPostId = data.data?.at(-5)?.id ?? "";

	let lastPostId = data.data?.at(-1)?.id ?? "";
	$: lastPostId = data.data?.at(-1)?.id ?? "";

	let InfiniteScrollObserver: IntersectionObserver;
	onMount(() => {
		if (data.type === "home_feed") {
			InfiniteScrollObserver = new IntersectionObserver(ObserverHandler);
			InfiniteScrollObserver.observe(document.getElementById(minusFifthPostId)!);
		}
	});

	let execAfterUpdate: Function[] = [];
	afterUpdate(() => {
		execAfterUpdate.forEach((fn) => fn());
		execAfterUpdate = [];
	});

	const ChangeFeedType = (active: number) => {
		if (active === FeedDisplayType) return;
		FeedDisplayType = active;
		HandleQueryingPost({ appendResult: false });
	};

	// Querying post observer
	type HandlePostParams = { afterId?: string; appendResult?: boolean };
	const HandleQueryingPost = async ({ afterId, appendResult }: HandlePostParams) => {
		if (data.type !== "home_feed") return;
		loading = !appendResult;
		const { success, data: newPosts } = await QueryHomePost(FeedDisplayType, afterId);
		if (success && newPosts && newPosts?.length > 0) {
			const newRender = appendResult ? [...(data.data || []), ...newPosts] : newPosts;

			execAfterUpdate.push(() =>
				InfiniteScrollObserver.observe(document.getElementById(newRender?.at(-5)?.id ?? "")!)
			);
			data.data = newRender;
		}
		loading = false;
	};

	const ObserverHandler: IntersectionObserverCallback = async ([lastPost]) => {
		if (!lastPost.isIntersecting) return;
		QueryMorePost();
	};

	const QueryMorePost = async () => {
		try {
			InfiniteScrollObserver.unobserve(document.getElementById(minusFifthPostId)!);
			HandleQueryingPost({
				afterId: lastPostId,
				appendResult: true
			});
		} catch (err) {}
	};
</script>

<main class="max-w-[1000px] m-auto flex justify-center px-2 py-5">
	<div class="max-w-[750px] md:mt-0 mt-2">
		{#if data.type === "home_feed"}
			<Tabs
				elems={[
					`<i class="fa-brands fa-hotjar icon"></i> <span>Hot</span>`,
					`<i class="fa-solid fa-star icon"></i> <span>New</span>`,
					`<i class="fa-solid fa-bolt icon"></i> <span>Top</span>`,
					`<i class="fa-solid fa-chart-line icon"></i> <span>Rising</span>`,
					`<i class="fa-solid fa-comment icon"></i> <span>Controversial</span>`
				]}
				active={FeedDisplayType}
				cb={ChangeFeedType}
			/>
		{/if}
		<Feeds rawPosts={data?.data || []} {loading} />
		{#if data.type === "user_feed"}
			<div class="text-center">
				<AlertChild type="info"
					>This is the end of your feed, to see more content: <Link href="/teddit?type=home_feed">
						<button class="btn btn-warning">Click Here</button>
					</Link>
				</AlertChild>
			</div>
		{:else}
			<div class="flex flex-col items-center">
				<button class="btn gap-x-2" on:click={QueryMorePost}
					>Next <i class="fa-solid fa-arrow-right icon" /></button
				>
			</div>
		{/if}
	</div>
</main>

<style scoped>
	main {
		min-height: calc(100vh - 64px);
	}
	:global(.tabs-component p) {
		font-size: 1.1rem;
		border: none;
		color: #ddd;
		transition: 0.2s all;
		border: none;
		border-top: solid 2px #484848;

		margin-bottom: 35px;
		column-gap: 5px;
	}
	:global(.tabs-component p.tab-active) {
		font-weight: bold;
		color: #fff;
		border-radius: 0px;
		border-top: solid 2px #ff4500;
		background: transparent;
	}
	:global(.tabs-component p.tab-active:hover) {
		color: #fff;
	}
	@media (max-width: 640px) {
		:global(.tabs-component p) {
			font-size: 0.8rem;
			column-gap: 2px;
		}
	}
	@media (max-width: 470px) {
		:global(.tabs-component p span) {
			display: none;
		}
		:global(.tabs-component i) {
			font-size: 1.1rem;
		}
	}
</style>
