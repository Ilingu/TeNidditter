<script lang="ts">
	import type { NeetComment } from "$lib/types/interfaces";
	import { IsEmptyString, isValidUrl } from "$lib/utils";
	import { afterUpdate } from "svelte";
	import Neet from "./Neet.svelte";
	import type Hls from "hls.js";
	import PictureZoom from "./PictureZoom.svelte";

	export let neets: NeetComment[][];
	export let queryMoreCb = () => {};

	let feedDiv: HTMLDivElement;
	const initFeed = async () => {
		if (!feedDiv) return;

		feedDiv.querySelectorAll(".neet .neet-body a").forEach((a) => {
			if (a.classList.contains("not")) return;
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
</script>

<div class="flex flex-col gap-y-5" bind:this={feedDiv}>
	{#each neets as neetThread}
		{#if neetThread?.length > 0 && neetThread[0].createdAt > 0}
			<div>
				{#each neetThread as neet, i}
					<Neet
						{neet}
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
