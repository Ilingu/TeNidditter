<script lang="ts">
	import { enhance } from "$app/forms";
	import { page } from "$app/stores";
	import Link from "$lib/client/components/design/Link.svelte";
	import Loader from "$lib/client/components/design/Loader.svelte";
	import type { DBSubtedditsShape } from "$lib/client/types/teddit";
	import api from "$lib/shared/api";
	import { FormatUsername, IsEmptyString } from "$lib/shared/utils";
	import { afterUpdate } from "svelte";

	let IsSubExist: boolean | null;
	$: IsSubExist = $page.form; // contains the result of the form submit

	let Input: HTMLInputElement;

	let subname = ""; // input subname
	let loading = false; // loading state

	let execAfterUpdate: Function[] = []; // little trick to execute certains function after the rerender of the component tree
	afterUpdate(() => {
		execAfterUpdate.forEach((fn) => fn());
		execAfterUpdate = [];
	});

	let timeout: NodeJS.Timer;
	const handleDebounce = () => {
		clearTimeout(timeout);
		timeout = setTimeout(FetchDBSubteddit, 500); // if user stopped typing after 500ms, query db
		if (IsSubExist !== null && !IsEmptyString(subname)) IsSubExist = null;
	};

	let MatchedSubs: DBSubtedditsShape[] = [];
	const FetchDBSubteddit = async () => {
		const SecureSubname = encodeURI(FormatUsername(subname));
		if (IsEmptyString(SecureSubname) || SecureSubname.length < 3) return;

		loading = true; // loading state true

		// Fetch
		const { success, data: matchedSubs } = await api.get("/teddit/r/search", {
			query: { q: SecureSubname }
		});

		// Reset
		execAfterUpdate.push(() => Input?.focus());
		loading = false;

		// Display
		if (!success || typeof matchedSubs !== "object" || matchedSubs.length <= 0) return;
		MatchedSubs = matchedSubs;
		document.addEventListener("click", WaitClick);
	};

	const WaitClick = (ev: MouseEvent) => {
		const elemType = (ev.target as HTMLElement).tagName;
		if (elemType === "SELECT" || elemType === "OPTION") return;

		MatchedSubs = [];
		document.removeEventListener("click", WaitClick);
	};
</script>

<main
	class="page-content max-w-[1150px] m-auto flex flex-col justify-center items-center h-screen px-2 py-5 gap-y-2"
>
	<header>
		<h1 class="text-white font-bold text-xl">🔎 Search a teddit subteddit</h1>
	</header>
	<form
		method="POST"
		class="bg-[#2e2d2f] lg:w-1/2 sm:w-3/4 w-full p-5 rounded-lg flex flex-col items-center relative"
		use:enhance={() => {
			loading = true;

			return async ({ update }) => {
				update({ reset: true });
				loading = false;
			};
		}}
	>
		<div class="form-control w-full max-w-xs">
			<label class="label" for="searchSubtedditInput">
				<span class="label-text">What is the subteddit?</span>
			</label>
			<div class="input-group">
				<input
					id="searchSubtedditInput"
					disabled={loading}
					name="subteddit"
					type="text"
					bind:this={Input}
					bind:value={subname}
					on:input={handleDebounce}
					placeholder="e.g: Szinek"
					class="input input-primary input-bordered w-full max-w-xs"
				/>
				<button class="btn btn-primary btn-square" type="submit">
					{#if loading}
						<Loader w={24} h={24} />
					{:else}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
							/></svg
						>
					{/if}
				</button>
			</div>
			{#if MatchedSubs.length > 0}
				<select
					class="select select-primary select-bordered text-center w-full max-w-xs p-1"
					multiple
				>
					{#each MatchedSubs as subteddit}
						<option
							on:dblclick={({ currentTarget: { value } }) => (subname = value)}
							value={subteddit.subname}
							class="text-white font-bold text-xl my-2">{subteddit.subname}</option
						>
					{/each}
				</select>
			{/if}
		</div>

		{#if IsSubExist !== null}
			<p class={`mt-2 text-lg font-semibold ${IsSubExist ? "text-success" : "text-error"}`}>
				This subteddit {IsSubExist ? "" : "doesn't"} exist
			</p>
			{#if IsSubExist}
				<Link href={`/teddit/r/${subname}`}
					><i class="fa-solid fa-arrow-up-right-from-square" /> {subname}'s Profile</Link
				>
			{/if}
		{/if}
	</form>

	{#if false}
		<!-- Why do I do that? simply because if I don't my "a" tag style below won't be include in this file and I want them for the style of "Link" component which is an... "a" tag! -->
		<a href="/" class="hidden">Easter Eggs</a>
	{/if}
</main>

<style scoped>
	a {
		font-style: italic;
		color: #92bddf;
		transition: all;
	}
	a:hover {
		text-decoration: underline wavy;
		color: #5296dd;
	}
</style>
