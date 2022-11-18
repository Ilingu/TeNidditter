<script lang="ts">
	import api from "$lib/api";
	import { page } from "$app/stores";
	import Feeds from "$lib/components/pages/nitter/Feeds.svelte";
	import { FormatNumbers, isValidUrl, MakeBearerToken, pushAlert } from "$lib/utils";
	import AuthStore from "$lib/stores/auth";

	export let data: import("./$types").PageData;

	let isSub = $AuthStore.Subs?.nitter?.includes($page.params.nittos);
	$: isSub = $AuthStore.Subs?.nitter?.includes($page.params.nittos);

	type MoveEvent = MouseEvent & { currentTarget: EventTarget & HTMLImageElement };
	const AnimateBanner = (ev: MoveEvent) => {
		const xAxis = (innerWidth / 2 - ev.pageX) / 25,
			yAxis = (innerHeight / 2 - ev.pageY) / 25;
		ev.currentTarget.style.transform = `rotateY(${xAxis}deg) rotateX(${yAxis}deg)`;
	};

	const TriggerSub = async () => {
		if (isSub) {
			isSub = false;
			const { success } = await api.delete("/tedinitter/nitter/unsub/%s", {
				params: [$page.params.nittos],
				headers: MakeBearerToken($AuthStore.JwtToken ?? "")
			});
			if (!success) {
				isSub = true;
				pushAlert("couldn't unsubscribe you", "error");
			}
		} else {
			isSub = true;
			const { success } = await api.post("/tedinitter/nitter/sub/%s", {
				params: [$page.params.nittos],
				headers: MakeBearerToken($AuthStore.JwtToken ?? "")
			});
			if (!success) {
				isSub = false;
				pushAlert("couldn't subscribe you", "error");
			}
		}
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
	<div
		class="max-w-[1250px] flex md:flex-row flex-col justify-center gap-y-2 md:items-start items-center gap-x-10 z-0"
	>
		<div
			class="userInfo w-80 h-fit bg-neutral rounded-lg p-3 flex flex-col gap-y-4 md:sticky top-20"
		>
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
			{#if $AuthStore.loggedIn}
				<button
					on:click={TriggerSub}
					class={`btn mt-3 gap-x-2 ${isSub ? "btn-secondary" : "btn-primary"}`}
					><i class={`fas ${isSub ? "fa-minus" : "fa-plus"}`} />
					{isSub ? "Unsubscribe" : "Subscribe"}</button
				>
			{/if}
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
