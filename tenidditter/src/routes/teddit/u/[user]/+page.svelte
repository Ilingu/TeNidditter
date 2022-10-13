<script lang="ts">
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";

	export let data: import("./$types").PageData;
	console.log(data);
</script>

<main class="max-w-[1150px] m-auto flex justify-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feeds rawPosts={data.posts} />
	</div>

	<aside class="bg-[#2e2d2f] w-full p-5 rounded-lg h-1/4">
		<img
			src={`https://libreddit.de/style/t5_7e51y/styles/${data.icon_img
				.split("_")
				.slice(-2)
				.join("_")}`}
			alt=""
		/>
		<h1 class="text-2xl">@ {data.username}</h1>
	</aside>
</main>
