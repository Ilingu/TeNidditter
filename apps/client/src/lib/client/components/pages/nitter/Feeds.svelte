<script lang="ts">
	import AuthStore from "$lib/client/stores/auth";
	import { afterUpdate, onDestroy, onMount } from "svelte";
	import Neet from "./Neet.svelte";
	import type Hls from "hls.js";
	import PictureZoom from "./PictureZoom.svelte";
	import { fade } from "svelte/transition";
	import { page } from "$app/stores";
	import api from "$lib/shared/api";
	import type { NeetComment } from "$lib/shared/types/nitter";
	import {
		IsEmptyString,
		isValidUrl,
		MakeBearerToken,
		removeDuplicates,
		TrimNonDigitsChars
	} from "$lib/shared/utils";
	import { pushAlert } from "$lib/client/ClientUtils";
	import SSEClient from "$lib/client/sse";
	import type { ExternalLinksDatas } from "$lib/client/types/nitter";

	export let neets: NeetComment[][];
	export let queryMoreCb = () => {};
	export let onNeetRemovedFromList = (_neetId: string) => {};

	/* externalLinks streams */
	let client: SSEClient<"/nitter/stream-in-external-links">;
	onMount(async () => {
		client = new SSEClient("/nitter/stream-in-external-links");

		const connected = await client.connect();
		if (!connected) return;

		client.on("message", (externalLink) => {
			if (!externalLink) return;
			externalLinks.push(externalLink);
		});
		client.on("close", StreamInExternalLinks);
	});
	onDestroy(() => {
		client && client.close();
	});

	let externalLinks: ExternalLinksDatas[] = [];
	const StreamInExternalLinks = () => {
		if (externalLinks.length <= 0) return;
		const extLinksNeetId = externalLinks.map(({ neetId }) => neetId);
		neets = neets.map((thread) =>
			thread.map((neet) => {
				if (extLinksNeetId.includes(neet.id)) {
					const externalLinkData = externalLinks.find(({ neetId }) => neetId === neet.id);
					if (!externalLinkData) return neet;

					return { ...neet, externalLinkMetatags: externalLinkData };
				} else return neet;
			})
		);
	};

	/* feed init-> videos processing, link tag cutsomization, neet observer (to querMore)... */
	let feedDiv: HTMLDivElement;
	const initFeed = async () => {
		if (!feedDiv) return;

		feedDiv.querySelectorAll(".neet .neet-body a").forEach((a) => {
			if (a.classList.contains("not")) return;
			if (a.innerHTML.includes("fa-solid fa-arrow-up-right-from-square icon")) return;
			const href = a.getAttribute("href");

			if (href?.startsWith("http")) {
				if (isValidUrl(a.textContent || "")) a.textContent = new URL(a.textContent || "").host;
				a.innerHTML = `<i class="fa-solid fa-arrow-up-right-from-square icon"></i>` + a.innerHTML;

				a.setAttribute("target", "_blank");
				a.setAttribute("rel", "noopener noreferrer");
			} else a.setAttribute("href", "/nitter" + href);
		});

		const NeetVideos = feedDiv.querySelectorAll(".neet video");
		if (NeetVideos.length <= 0) return;

		HlsInstance = (await import("hls.js")).default;
		NeetVideos.forEach((vid) => playVideo(vid as HTMLMediaElement));
	};
	const observeQueryMore = () => {
		const lastCommentId = [5, 4, 3, 2, 1].find(
			(i) => !IsEmptyString((neets.at(-i) ?? [{ id: "" }])[0]?.id)
		);
		if (!lastCommentId) return;

		const latest5Comment = document.getElementById(`neet-${neets.at(-lastCommentId)![0].id}`);
		if (!latest5Comment) return;

		const commentOberser: IntersectionObserverCallback = ([comment]) => {
			if (!comment.isIntersecting) return;
			InfiniteScrollObserver.unobserve(latest5Comment);
			queryMoreCb();
		};

		const InfiniteScrollObserver = new IntersectionObserver(commentOberser);
		InfiniteScrollObserver.observe(latest5Comment);
	};

	$: if (neets) execAfterUpdate.add(initFeed).add(observeQueryMore);

	let execAfterUpdate = new Set<Function>();
	afterUpdate(() => {
		execAfterUpdate.forEach((fn) => fn());
		execAfterUpdate.clear();
	});

	let HlsInstance: typeof Hls;
	const playVideo = (vid: HTMLMediaElement) => {
		if (HlsInstance.isSupported()) {
			var hls = new HlsInstance();
			hls.loadSource(vid.getAttribute("src") ?? "");
			hls.attachMedia(vid);
			hls.on(HlsInstance.Events.MANIFEST_PARSED, function () {
				hls.loadLevel = hls.levels.length - 1;
				hls.startLoad();
			});
		} else if (vid.canPlayType("application/vnd.apple.mpegurl")) {
			vid.addEventListener("canplay", function () {});
		}
	};

	/* user nitter lists */

	let neetIdToAdd = "";
	const setNeetIdToAdd = (neetId: string) => {
		neetIdToAdd = TrimNonDigitsChars(neetId);

		// This part is only when user is already on his list page (e.g: /nitter/list/[id])
		// On this page all the neets passed in this components are for sure neets from his list so we add it to the "neetAddedInList" variable when the users clicks on one of these neets
		const pageListId = TrimNonDigitsChars($page.params?.listId);
		if (!IsEmptyString(pageListId) && !isNaN(parseInt(pageListId)))
			neetAddedInList = removeDuplicates([...neetAddedInList, parseInt(pageListId)]);
	};

	let neetAddedInList: number[] = []; // this variable store in which listId this neet (neetIdToAdd) is saved
	const AddNeetToList = async (listId: number) => {
		if (IsEmptyString(neetIdToAdd)) return pushAlert("Please select a neet first", "warning");
		if (isNaN(listId)) return pushAlert("Invalid List", "warning");

		if (neetAddedInList.includes(listId)) {
			const { success: deleted } = await api.delete("/tedinitter/nitter/list/%s/removeNeet/%s", {
				headers: MakeBearerToken($AuthStore.JwtToken ?? ""),
				params: [`${listId}`, neetIdToAdd]
			});
			if (!deleted) return pushAlert("Couldn't remove this neet from this list", "error", 6000);
			neetAddedInList = neetAddedInList.filter((listIdAdded) => listIdAdded !== listId);
			onNeetRemovedFromList(neetIdToAdd);
		} else {
			const neetDatas = Array<NeetComment>()
				.concat(...neets)
				.find(({ id }) => TrimNonDigitsChars(id) === TrimNonDigitsChars(neetIdToAdd));
			if (!neetDatas) return pushAlert("Couldn't find this neet", "error");

			const { success: saved } = await api.post("/tedinitter/nitter/list/%s/saveNeet", {
				headers: MakeBearerToken($AuthStore.JwtToken ?? ""),
				params: [`${listId}`],
				body: neetDatas
			});
			if (!saved)
				return pushAlert(
					"Couldn't add this neet to this list (probably because it's already in this list)",
					"error",
					10000
				);
			neetAddedInList = removeDuplicates([...neetAddedInList, listId]);
		}
	};
</script>

<div class="flex flex-col gap-y-5" bind:this={feedDiv}>
	{#each neets as neetThread}
		{#if neetThread?.length > 0 && neetThread[0].createdAt > 0}
			<div>
				{#each neetThread as neet, i}
					<Neet
						{neet}
						{setNeetIdToAdd}
						thread={[
							neetThread.length > 1,
							i === 0 ? "first" : i === neetThread.length - 1 ? "last" : undefined
						]}
					/>
				{/each}
			</div>
		{/if}
	{/each}
</div>
<PictureZoom />

<input type="checkbox" id="modal-add-to-list" class="modal-toggle" />
<div class="modal">
	<div class="modal-box">
		<h3 class="font-bold text-lg">Add this tweet to your list</h3>
		<p class="py-4">Here all your lists:</p>
		<div class="max-h-[300px] overflow-y-auto flex flex-col items-center">
			{#if $AuthStore?.Lists && $AuthStore.Lists.length > 0 && !IsEmptyString(neetIdToAdd)}
				<div class="w-[400px] flex flex-col gap-y-3 mt-2 cursor-pointer">
					{#each $AuthStore.Lists as list, i}
						<div
							in:fade={{ delay: 15 * i }}
							on:click={() => AddNeetToList(list.list_id)}
							class="flex gap-x-2 group items-center h-14 bg-neutral rounded-md w-full p-2"
						>
							<span
								class="ml-2 from-accent to-secondary bg-gradient-to-br bg-clip-text text-clip text-transparent text-bold text-xl capitalize flex-1"
								>{list.title}</span
							>

							{#key neetAddedInList}
								{#if neetAddedInList.includes(list.list_id)}
									<span class="text-success mr-2 text-lg" transition:fade
										><i class="fa-solid fa-circle-check" /></span
									>
								{/if}
							{/key}
						</div>
					{/each}
				</div>
			{:else}
				<h2 class="text-xl mt-4">Nothing to show! Create a list first.</h2>
			{/if}
		</div>
		<div class="modal-action">
			<label
				for="modal-add-to-list"
				class="btn"
				on:click={() => {
					neetIdToAdd = "";
					neetAddedInList = [];
				}}>Done!</label
			>
		</div>
	</div>
</div>
