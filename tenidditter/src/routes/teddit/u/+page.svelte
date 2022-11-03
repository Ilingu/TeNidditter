<script lang="ts">
	import { enhance } from "$app/forms";
	import { page } from "$app/stores";
	import Link from "$lib/components/design/Link.svelte";
	import Loader from "$lib/components/design/Loader.svelte";
	import { IsEmptyString } from "$lib/utils";

	let IsUserExist: boolean | null;
	$: IsUserExist = $page.form;

	let username = "";
	let loading = false;
</script>

<main
	class="page-content max-w-[1150px] m-auto flex flex-col justify-center items-center h-screen px-2 py-5 gap-y-2"
>
	<header>
		<h1 class="text-white font-bold text-xl">🔎 Search a teddit user</h1>
	</header>
	<form
		method="POST"
		class="bg-[#2e2d2f] lg:w-1/2 sm:w-3/4 w-full p-5 rounded-lg flex flex-col items-center"
		use:enhance={() => {
			loading = true;

			return async ({ update }) => {
				update({ reset: true });
				loading = false;
			};
		}}
	>
		<div class="form-control w-full max-w-xs">
			<label class="label" for="searchUserInput">
				<span class="label-text">What is they username?</span>
			</label>
			<div class="input-group">
				<input
					id="searchUserInput"
					name="username"
					type="text"
					disabled={loading}
					bind:value={username}
					on:input={() => {
						if (IsUserExist !== null && !IsEmptyString(username)) IsUserExist = null;
					}}
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
		</div>

		{#if IsUserExist !== null}
			<p class={`mt-2 text-lg font-semibold ${IsUserExist ? "text-success" : "text-error"}`}>
				This user {IsUserExist ? "" : "doesn't"} exist
			</p>
			{#if IsUserExist}
				<Link href={`/teddit/u/${username}`}
					><i class="fa-solid fa-arrow-up-right-from-square" /> {username}'s Profile</Link
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
