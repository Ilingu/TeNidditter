<script lang="ts">
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";

	export let data: import("./$types").PageData;
	console.log(data);
</script>

<!-- DO SCRAPPING FOR Query info about subreddit -->

<main class="max-w-[1500px] m-auto flex justify-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feeds rawPosts={data?.data || []} queryMorePostHandler={async () => {}} />
	</div>
	<!-- Subreddit's info -->
	<aside class="max-w-[350px] min-w-[280px] h-screen bg-white rounded-lg" />
</main>
