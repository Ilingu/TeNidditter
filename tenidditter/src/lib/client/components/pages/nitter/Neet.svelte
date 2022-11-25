<script lang="ts">
	import AuthStore from "$lib/client/stores/auth";
	import Link from "$lib/client/components/design/Link.svelte";
	import type { openImgArgs } from "$lib/client/types/zoomImg";
	import type { NeetComment } from "$lib/shared/types/nitter";
	import { IsEmptyString, TrimNonDigitsChars } from "$lib/shared/utils";
	import { FormatNumbers, humanElapsedTime } from "$lib/client/ClientUtils";

	export let neet: NeetComment;
	export let quoteMode = false;

	type threadType = [isThread: boolean, threadId: "first" | "last" | undefined];
	export let thread: threadType = [false, undefined];
	export let setNeetIdToAdd: (neetId: string) => void;

	const openImageDrawer = (urls: string[], currIndex: number) => {
		const alert = new CustomEvent("openImageDrawer", {
			detail: {
				urls,
				currIndex
			} as openImgArgs
		});
		document.dispatchEvent(alert);
	};
</script>

<div
	id={"neet-" + TrimNonDigitsChars(neet.id)}
	class={`neet relative break-words ${
		quoteMode ? "w-full text-primary-content" : "md:w-[500px] w-[95vw]"
	} gap-x-2 ${
		quoteMode ? "bg-secondary-focus" : "bg-primary-focus hover:bg-primary"
	} rounded-md p-2 transition-all ${
		thread[0]
			? thread[1] === "first"
				? "rounded-b-none"
				: thread[1] === "last"
				? "rounded-t-none"
				: "rounded-none"
			: ""
	}`}
>
	{#if thread[0]}
		<div class="absolute left-[28px] bottom-1 z-0 w-1 h-[calc(100%-65px)] bg-white rounded-full" />
	{/if}
	<div class="neet-pp">
		<img
			src={neet.creator.avatarUrl}
			alt="User Profile Pic"
			style="border-radius: 10% 65% / 50%;"
		/>
	</div>

	<header class="neet-header mb-2">
		{#if !IsEmptyString(neet.retweeted)}
			<p class="text-neutral-focus">
				<i class="fa-solid fa-retweet" />
				{neet.retweeted}
			</p>
		{/if}
		{#if neet.pinned}
			<p class="text-neutral-focus">
				<i class="fa-solid fa-thumbtack" /> pinned
			</p>
		{/if}

		<div class="flex flex-wrap gap-x-2">
			<div class="flex-1 text-accent font-bold">
				<Link href={`/nitter/${neet.creator.username.slice(1)}`}
					>{neet.creator.username}
					{#if thread[0] && thread[1] !== "first"}
						<span class="text-neutral-focus"><i class="fa-solid fa-reply" /> replied</span>
					{/if}</Link
				>
			</div>

			<div
				class="tooltip tooltip-secondary"
				data-tip={new Date(neet.createdAt * 1000).toLocaleString()}
			>
				{humanElapsedTime(neet.createdAt * 1000, Date.now())}
			</div>
		</div>
	</header>

	<div class="neet-body">
		<Link
			href={`/nitter/${neet.creator.username.slice(1)}/${TrimNonDigitsChars(neet.id)}`}
			className="not"
		>
			<p class="font-bold">{@html neet.content}</p>
		</Link>
		{#if neet?.externalLinkMetatags}
			<div class="grid gap-x-2 grid-cols-10 bg-neutral-focus p-2 rounded-lg">
				<div class="col-span-2 flex justify-center items-center">
					{#if !IsEmptyString(neet.externalLinkMetatags["og:image"])}
						<img src={neet.externalLinkMetatags["og:image"]} alt="" class="w-fit h-fit rounded" />
					{/if}
				</div>
				<div class="col-span-8">
					{#if !IsEmptyString(neet.externalLinkMetatags?.title)}
						<h1>{neet.externalLinkMetatags.title}</h1>
					{/if}
					{#if !IsEmptyString(neet.externalLinkMetatags?.description)}
						<p>{neet.externalLinkMetatags.description}</p>
					{/if}
				</div>
			</div>
		{/if}
		{#if neet.quote}
			<div>
				<svelte:self neet={neet.quote} quoteMode={true} {setNeetIdToAdd} />
			</div>
		{/if}
	</div>

	{#if Object.keys(neet.attachment ?? {}).length > 0}
		<div class="neet-attachments">
			{#if (neet.attachment?.images?.length ?? 0) > 0}
				<div
					class={`imgs place-items-center grid ${
						(neet.attachment?.images?.length ?? 0) > 1 ? "grid-cols-2" : "grid-cols-1"
					} gap-1`}
				>
					{#each neet.attachment?.images ?? [] as imgUrl, i}
						<img
							on:click={() => openImageDrawer(neet.attachment?.images ?? [], i)}
							src={imgUrl}
							class="rounded cursor-zoom-in max-h-[512px] w-auto"
							alt="🖼"
						/>
					{/each}
				</div>
			{/if}
			{#if (neet.attachment?.videos?.length ?? 0) === 1}
				<div class="vid">
					<video
						src={(neet.attachment?.videos ?? [])[0] ?? ""}
						controls
						preload="none"
						class="my-5 max-h-[512px] w-auto"
					>
						<track kind="captions" />
					</video>
				</div>
			{/if}
		</div>
	{/if}
	<div class="neet-stats flex gap-x-2 text-neutral mt-1">
		<p title="number of likes">
			<i class="fa-solid fa-heart" />
			{FormatNumbers(neet.stats.likes_counts)}
		</p>
		<p title="number of retweets">
			<i class="fa-solid fa-retweet" />
			{FormatNumbers(neet.stats.rt_counts)}
		</p>
		<p title="number of replies" class={neet.stats.play_counts !== undefined ? "" : "flex-1"}>
			<i class="fa-solid fa-reply" />
			{FormatNumbers(neet.stats.reply_counts)}
		</p>

		{#if neet.stats.play_counts !== undefined}
			<p title="number of plays" class="flex-1">
				<i class="fa-solid fa-play" />
				{FormatNumbers(neet.stats.play_counts)}
			</p>
		{/if}

		{#if $AuthStore.loggedIn}
			<label for="modal-add-to-list" on:click={() => setNeetIdToAdd(neet.id)}>
				<div class="tooltip tooltip-secondary" data-tip={"Add to list"}>
					<i class="fa-solid fa-bookmark" />
				</div></label
			>
		{/if}
	</div>
</div>

<style scoped>
	.neet {
		display: grid;

		grid-template-areas:
			"neet-pp neet-header neet-header"
			"neet-pp neet-body neet-body"
			"neet-pp neet-attachments neet-attachments"
			"neet-pp neet-stats neet-stats";
		grid-template-rows: auto;
		grid-template-columns: 45px 1fr;
	}

	.neet-pp {
		grid-area: neet-pp;
	}
	.neet-header {
		grid-area: neet-header;
	}
	.neet-body {
		grid-area: neet-body;
	}
	.neet-attachments {
		grid-area: neet-attachments;
	}
	.neet-stats {
		grid-area: neet-stats;
	}

	:global(.neet .neet-body a:not(.not)) {
		color: #f3cc30;
		font-weight: bold;
		transition: 0.3s all;
	}
	:global(.neet .neet-body a:not(.not):hover) {
		text-decoration: underline wavy;
		color: #58c7f3;
	}
</style>
