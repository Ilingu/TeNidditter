<script lang="ts">
	import type { NeetComment } from "$lib/types/interfaces";
	import { isValidUrl } from "$lib/utils";
	import { afterUpdate } from "svelte";
	import Neet from "./Neet.svelte";
	import type Hls from "hls.js";
	import PictureZoom from "./PictureZoom.svelte";

	export let neets: NeetComment[][];

	let feedDiv: HTMLDivElement;
	const initFeed = async () => {
		console.log("Init");
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

	$: if (neets) execAfterUpdate.add(initFeed);

	let execAfterUpdate = new Set<Function>();
	afterUpdate(() => {
		console.log(execAfterUpdate);
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
