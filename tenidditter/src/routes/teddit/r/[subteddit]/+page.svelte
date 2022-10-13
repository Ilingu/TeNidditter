<script lang="ts">
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";
	import { page } from "$app/stores";
	import AuthStore from "$lib/stores/auth";

	export let data: import("./$types").PageData;
</script>

<main class="max-w-[1500px] m-auto flex justify-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feeds rawPosts={data?.Feed || []} queryMorePostHandler={async () => {}} />
	</div>
	<!-- Subreddit's info -->
	<aside
		class="max-w-[350px] min-w-[280px] px-3 min-h-screen bg-[rgba(45,49,49,0.71)] ring-1 ring-[#686868] flex flex-col items-center rounded-lg"
	>
		<h1 class="text-2xl font-bold	mt-2 font-fancy tracking-widest text-teddit">
			{$page.params.subteddit}
		</h1>

		<p class="mt-2">{data.Info?.description}</p>
		<p class="mt-2 underline text-nitter capitalize">{data.Info?.subs}</p>
		{#if $AuthStore.loggedIn}
			<button>Subscribe</button>
		{/if}

		<h1 class="text-xl text-center mt-20">Rules:</h1>
		{@html data.Info?.rules}
	</aside>
</main>

<style scoped>
	:global(ol) {
		list-style: inside decimal;
	}
</style>
