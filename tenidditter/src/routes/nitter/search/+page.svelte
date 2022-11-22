<script lang="ts">
	import { page } from "$app/stores";
	import api from "$lib/api";
	import Loader from "$lib/client/components/design/Loader.svelte";
	import Feeds from "$lib/client/components/pages/nitter/Feeds.svelte";
	import NittosPreview from "$lib/client/components/pages/nitter/NittosPreview.svelte";
	import { IsEmptyString, pushAlert, TrimSpecialChars } from "$lib/utils";
	import { onMount } from "svelte";

	let dataType = ($page.url.searchParams.get("type") as "tweets" | "users") ?? "tweets";
	export let data: any; // import("./$types").PageData;

	let activeTab: "tweets" | "users" = dataType || "tweets";
	let query = $page.url.searchParams.get("q") ?? "";

	let loading = false;

	onMount(() => UpdateURLSearch());
	const UpdateURLSearch = () => {
		if (IsEmptyString(activeTab) || IsEmptyString(query)) return;

		window.history.pushState(
			null,
			"",
			location.pathname + `?q=${encodeURIComponent(query)}&type=${activeTab}`
		);
	};

	const SearchWrap = async () => {
		loading = true;
		await Search();
		loading = false;
	};

	const Search = async (limit = 5) => {
		if (IsEmptyString(activeTab) || IsEmptyString(query)) return;
		if (activeTab === "users") query = TrimSpecialChars(query).replaceAll(" ", ""); // trim special char

		const { success, data: searchResult } = await api.get("/nitter/search", {
			query: { q: query, type: activeTab, limit }
		});
		if (!success || typeof searchResult !== "object")
			return pushAlert("Nitter returned nothing", "error", 4000);

		dataType = activeTab;
		data = { searchResult };
		UpdateURLSearch();
	};

	let lastLimit = 5;
	const queryMore = () => {
		Search(lastLimit + 1);
		lastLimit++;
	};
</script>

<main class="grid place-items-center mt-5">
	<div class="sm:w-[500px] flex flex-col items-center">
		<form on:submit|preventDefault={SearchWrap} class="w-full flex flex-col items-center">
			<div class="form-control bg-neutral w-full">
				<div class="input-group w-full">
					<input
						type="text"
						name="query"
						bind:value={query}
						placeholder={`Search ${activeTab}`}
						class="input input-bordered input-accent w-full bg-neutral"
					/>
					<button class="btn btn-square" type="submit" disabled={loading}>
						{#if loading}
							<Loader h={20} w={20} />
						{:else}
							<i class="fa-solid fa-magnifying-glass" />
						{/if}
					</button>
				</div>
			</div>
		</form>
		<div class="w-full h-8 rounded bg-neutral mt-1 grid grid-cols-2 text-center">
			<button
				class={`gap-x-2  ${activeTab === "tweets" ? "border-b-2 border-accent" : ""}`}
				on:click={() => {
					activeTab = "tweets";
					SearchWrap();
				}}><i class="fa-brands fa-twitter" /> Tweets</button
			>
			<button
				class={`gap-x-2 ${activeTab === "users" ? "border-b-2 border-accent" : ""}`}
				on:click={() => {
					activeTab = "users";
					SearchWrap();
				}}
			>
				<i class="fa-solid fa-user" /> Users</button
			>
		</div>
		<div class="mt-2">
			{#if data}
				{#if dataType === "tweets"}
					<Feeds neets={data.searchResult} queryMoreCb={queryMore} />
				{:else}
					<NittosPreview previewDatas={data.searchResult} />
				{/if}
			{/if}
		</div>
	</div>
</main>
