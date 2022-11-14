<script lang="ts">
	import Feeds from "$lib/components/pages/nitter/Feeds.svelte";

	export let data: import("./$types").PageServerData;

	const onNeetRemovedFromList = (neetId: string) => {
		data.savedNeets = data.savedNeets?.filter(({ id }) => id !== neetId) ?? null;
	};
</script>

<main class="grid place-items-center mt-5">
	{#if data.savedNeets && (data.savedNeets || []).length > 0}
		<div class="max-w-[750px]">
			<Feeds neets={data.savedNeets.map((thread) => [thread])} {onNeetRemovedFromList} />
		</div>
	{:else}
		<div class="h-[calc(100vh-92px)] mx-4 flex flex-col justify-center">
			<div
				class="md:h-auto h-2/6 md:block flex flex-col justify-center gap-y-5 p-10 bg-neutral rounded-2xl"
			>
				<h1 class="md:text-2xl text-xl font-bold text-accent">
					Nothing saved in this list 🙀... yet 😏
				</h1>
				<p>To add a neet here go to your feed and pin some neets!</p>
			</div>
		</div>
	{/if}
</main>
