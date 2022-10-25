<script lang="ts">
	import Feeds from "$lib/components/pages/teddit/Feeds.svelte";
	import { page } from "$app/stores";
	import AuthStore from "$lib/stores/auth";
	import { styles } from "$lib/services/style";
	import api from "$lib/api";
	import { MakeBearerToken, pushAlert } from "$lib/utils";

	export let data: import("./$types").PageData;

	let isSub = false;
	$: if ($AuthStore.loggedIn) isSub = !!$AuthStore.Subs?.teddit?.includes($page.params.subteddit);

	const TriggerSub = async () => {
		const copySub = isSub;
		isSub = !isSub; // so that the user don't wait 2s

		if (copySub) {
			const { success } = await api.delete("/tedinitter/teddit/unsub", {
				param: $page.params.subteddit,
				headers: MakeBearerToken($AuthStore.JwtToken || "")
			});
			if (!success) {
				pushAlert("Couldn't unsubscribe you, try again", "error");
				isSub = copySub;
			}
		} else {
			const { success } = await api.post("/tedinitter/teddit/sub", {
				param: $page.params.subteddit,
				headers: MakeBearerToken($AuthStore.JwtToken || "")
			});
			if (!success) {
				pushAlert("Couldn't subscribe you, try again", "error");
				isSub = copySub;
			}
		}
	};
</script>

<main class="max-w-[1500px] m-auto flex justify-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feeds rawPosts={data?.Feed || []} />
	</div>
	<!-- Subreddit's info -->
	<aside class="flex flex-col gap-y-5">
		<header
			class="max-w-[350px] min-w-[280px] px-3 py-5 bg-light-dark ring-1 ring-[#686868] flex flex-col items-center rounded-lg"
		>
			<h1
				class="text-3xl font-bold	mt-2  tracking-widest text-teddit"
				use:styles={{ fontFamily: Math.round(Math.random() * 2) === 0 ? "Silkscreen" : "Pacifico" }}
			>
				{$page.params.subteddit}
			</h1>

			<p class="mt-2 text-justify">{data.Info?.description}</p>
			<p class="mt-2 underline text-nitter capitalize">{data.Info?.subs}</p>
			{#if $AuthStore.loggedIn}
				<button
					on:click={TriggerSub}
					class={`btn mt-3 gap-x-2 ${isSub ? "btn-secondary" : "btn-primary"}`}
					><i class={`fas ${isSub ? "fa-minus" : "fa-plus"}`} />
					{isSub ? "Unsubscribe" : "Subscribe"}</button
				>
			{/if}
		</header>
		<details class="max-w-[350px] p-3 bg-light-dark ring-1 ring-[#686868] rounded-lg">
			<summary class="text-xl font-bold">Rules</summary>
			{@html data.Info?.rules}
		</details>
	</aside>
</main>

<style scoped>
	aside > header > h1 {
		font-family: var(--fontFamily);
	}

	details > summary::marker {
		color: rgb(255 69 0);
	}

	:global(ol) {
		list-style: inside decimal;
	}
</style>
