<script lang="ts">
	import api from "$lib/api";
	import { page } from "$app/stores";
	import Feeds from "$lib/components/pages/nitter/Feeds.svelte";
	import { FormatNumbers, isValidUrl } from "$lib/utils";

	export let data: import("./$types").PageData;

	type MoveEvent = MouseEvent & { currentTarget: EventTarget & HTMLImageElement };
	const AnimateBanner = (ev: MoveEvent) => {
		const xAxis = (innerWidth / 2 - ev.pageX) / 25,
			yAxis = (innerHeight / 2 - ev.pageY) / 25;
		ev.currentTarget.style.transform = `rotateY(${xAxis}deg) rotateX(${yAxis}deg)`;
	};

	let lastLimit = 3;
	const queryMore = async () => {
		const { success, data: Neets } = await api.get("/nitter/nittos/%s/neets", {
			params: [$page.params.nittos],
			query: { limit: lastLimit + 1 }
		});
		if (!success || typeof Neets !== "object") return;
		lastLimit++;
		data.userNeets = Neets;
	};
</script>

<main class="grid place-items-center mt-5">
	<div class="max-w-[900px] mb-5">
		<img
			src={data.userInfo.bannerUrl}
			alt="banner"
			class="banner rounded-md select-none"
			draggable="false"
			on:mousemove={AnimateBanner}
			on:mouseenter={(ev) => (ev.currentTarget.style.transition = "none")}
			on:mouseleave={(ev) => {
				ev.currentTarget.style.transition = "all 0.5s ease 0s";
				ev.currentTarget.style.transform = "rotateY(0deg) rotateX(0deg)";
			}}
		/>
	</div>
	<div class="max-w-[1250px] flex justify-center gap-x-10 z-0">
		<div class="userInfo w-80 h-fit bg-neutral rounded-lg p-3 flex flex-col gap-y-4 sticky top-20">
			<div class="flex justify-center">
				<img
					src={data.userInfo.avatarUrl}
					alt="User Avatar"
					class="max-w-[256px] ring ring-primary ring-offset-base-100 ring-offset-2 rounded"
				/>
			</div>
			<header>
				<h1 class="text-xl font-bold capitalize text-accent">{data.userInfo.username}</h1>
				<p class="bio break-words">{@html data.userInfo.bio}</p>
			</header>
			<ul class="text-gray-300">
				<li><i class="fa-solid fa-location-pin" /> {data.userInfo.location}</li>
				{#if isValidUrl(data.userInfo.website)}
					<li>
						<a href={data.userInfo.website} style="font-weight: normal;"
							><i class="fa-solid fa-link" /> {new URL(data.userInfo.website).host}</a
						>
					</li>
				{/if}
				<li><i class="fa-solid fa-calendar" /> {data.userInfo.joinDate}</li>
			</ul>
			<div class="flex text-sm text-center">
				<p>
					<span class="font-bold">Neets</span>
					<span class="text-sm">{FormatNumbers(data.userInfo.stats.tweets_counts)}</span>
				</p>
				<p>
					<span class="font-bold">Likes</span>
					<span class="text-sm">{FormatNumbers(data.userInfo.stats.likes_counts)}</span>
				</p>
				<p>
					<span class="font-bold">Followers</span>
					<span class="text-sm">{FormatNumbers(data.userInfo.stats.followers_counts)}</span>
				</p>
				<p>
					<span class="font-bold">Following</span>
					<span class="text-sm">{FormatNumbers(data.userInfo.stats.following_counts)}</span>
				</p>
			</div>
		</div>
		<div class="max-w-[750px]">
			{#if data.userNeets}
				<Feeds neets={data.userNeets} queryMoreCb={queryMore} />
			{/if}
		</div>
	</div>
</main>

<style scoped>
	:global(.bio a) {
		font-size: 12px;
	}

	:global(.userInfo a) {
		display: block;
		color: #f3cc30;
		font-weight: bold;
		transition: 0.3s all;
	}
	:global(.userInfo a:hover) {
		text-decoration: underline wavy;
		color: #58c7f3;
	}
</style>
