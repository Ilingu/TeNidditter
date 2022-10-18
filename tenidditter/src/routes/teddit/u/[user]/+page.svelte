<script lang="ts">
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";
	import { GenerateRandomPPUrl } from "$lib/services/ProfilesPictures";
	import { FormatNumbers } from "$lib/utils";

	export let data: import("./$types").PageData;
	console.log(data);
</script>

<main class="max-w-[1150px] m-auto flex justify-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feeds rawPosts={data.posts} />
	</div>

	<aside class="bg-[#2e2d2f] w-full p-5 rounded-lg h-1/4">
		<div class="flex justify-center avatar mb-1">
			<div class="w-24 mask mask-hexagon">
				<img src={GenerateRandomPPUrl(96)} alt="fake profile pictures" />
			</div>
		</div>
		<h1 class="text-3xl mb-4 text-center text-teddit font-semibold tracking-wide">
			{data.username}
		</h1>
		<ul>
			<li>
				<i class="fa-solid fa-right-to-bracket" /> Created:
				<span class="font-bold text-white"
					>{new Date(data.created * 1000).toLocaleDateString()}</span
				>
			</li>
			<li>
				<i class="fa-solid fa-comment" /> Comment Karma:
				<span class="font-bold text-white">{FormatNumbers(data.comment_karma)}</span>
			</li>
			<li>
				<i class={`fa-solid ${data.verified ? "fa-check" : "fa-xmark"}`} /> Verified:
				<span class={`font-bold ${data.verified ? "text-success" : "text-warning"}`}
					>{data.verified}</span
				>
			</li>
		</ul>
	</aside>
</main>
