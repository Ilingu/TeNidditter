<script lang="ts">
	import Feeds from "$lib/client/components/pages/teddit/Feeds.svelte";
	import { page } from "$app/stores";
	import AuthStore from "$lib/client/stores/auth";
	import { styles } from "$lib/services/style";

	export let data: import("./$types").PageData; // ssr

	let isSub = false; // whether the user is subscribe to this subreddit or not
	$: if ($AuthStore.loggedIn) isSub = !!$AuthStore.Subs?.teddit?.includes($page.params.subteddit);

	// is user sub it'll unsubscribe him and vise-versa
	const TriggerSub = async () => {
		const copySub = isSub;
		isSub = !isSub; // so that the user don't wait 2s

		const { success } = await $AuthStore?.user.action?.toggleTedditSubs(
			$page.params.subteddit,
			copySub,
			$AuthStore.JwtToken || ""
		);

		if (!success) isSub = copySub;
	};
</script>

<main
	class="max-w-[1500px] m-auto flex xl:flex-row flex-col-reverse justify-center xl:items-start items-center xl:gap-y-0 gap-y-4 gap-x-8 px-2 py-5"
>
	<div class="max-w-[750px] md:mt-0 mt-2">
		<Feeds rawPosts={data?.Feed || []} />
	</div>
	<!-- Subreddit's info -->
	<aside class="flex flex-col gap-y-5">
		<header
			class="xl:max-w-[350px] max-w-[750px] min-w-[280px] px-3 py-5 bg-light-dark ring-1 ring-[#686868] flex flex-col items-center rounded-lg"
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
		<details
			class="xl:max-w-[350px] max-w-[750px] p-3 bg-light-dark ring-1 ring-[#686868] rounded-lg"
		>
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
