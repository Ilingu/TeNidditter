<script lang="ts">
	import ProfilePicture from "$lib/client/components/layout/ProfilePicture.svelte";
	import Feeds from "$lib/client/components/pages/teddit/Feeds.svelte";
	import { FormatNumbers } from "$lib/utils";

	export let data: import("./$types").PageData; // server data, ssr
</script>

<main
	class="max-w-[1500px] m-auto flex xl:flex-row flex-col-reverse justify-center xl:items-start items-center xl:gap-y-0 gap-y-4 gap-x-8 px-2 py-5"
>
	<div class="max-w-[750px] md:mt-0 mt-2">
		<Feeds rawPosts={data.posts} />
	</div>

	<aside
		class="bg-[#2e2d2f] xl:max-w-[350px] max-w-[750px] min-w-[280px] w-full p-5 rounded-lg h-1/4"
	>
		<div class="flex justify-center avatar mb-1">
			<ProfilePicture size="big" />
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
