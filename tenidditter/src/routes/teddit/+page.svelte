<!-- Feed -->
<script lang="ts">
	import Tabs from "$lib/components/design/Tabs.svelte";
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";
	import { QueryHomePost } from "$lib/services/teddit";
	import { FeedTypeEnum } from "$lib/types/enums";

	export let data: import("./$types").PageData;
	console.log(data);

	let FeedDisplayType = FeedTypeEnum.Hot;

	const ChangeFeedType = (active: number) => {
		FeedDisplayType = active;
		HandleQueryingPost();
	};

	const HandleQueryingPost = async (afterId?: string) => {
		const { success, data: newPosts } = await QueryHomePost(FeedDisplayType, afterId);
		if (success && newPosts && newPosts?.length > 0)
			data.data = [...(data.data || []), ...newPosts];
	};
</script>

<main class="max-w-[1000px] m-auto flex justify-center px-2 py-5">
	<div class="max-w-[750px]">
		{#if data.type === "home_feed"}
			<Tabs
				elems={[
					`<i class="fa-brands fa-hotjar icon"></i> Hot`,
					`<i class="fa-solid fa-star icon"></i> New`,
					`<i class="fa-solid fa-bolt icon"></i> Top`,
					`<i class="fa-solid fa-chart-line icon"></i> Rising`,
					`<i class="fa-solid fa-comment icon"></i> Controversial`
				]}
				active={FeedDisplayType}
				cb={ChangeFeedType}
			/>
		{/if}
		<Feeds rawPosts={data?.data || []} queryMorePostHandler={HandleQueryingPost} />
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
</style>
