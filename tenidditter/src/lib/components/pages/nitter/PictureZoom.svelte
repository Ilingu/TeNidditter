<script lang="ts">
	import type { openImgArgs } from "$lib/types/interfaces";
	import { onMount } from "svelte";

	onMount(() => {
		document.addEventListener("openImageDrawer", handleOpenImg as EventListener);
	});
	const handleOpenImg = ({ detail: { urls, currIndex } }: CustomEvent<openImgArgs>) =>
		openImageDrawer(urls, currIndex);

	let imgDrawer: [show: boolean, imgUrl: string[], currentImg: number] = [false, [], 0];
	const openImageDrawer = (imgUrls: string[], activeId: number) => {
		imgDrawer = [true, imgUrls, activeId];
		document.body.style.overflow = "hidden";
	};
	const closeImageDrawer = () => {
		imgDrawer = [false, [], 0];
		document.body.style.overflow = "visible";
	};

	type AndInfoEvent = { currentTarget: EventTarget & HTMLImageElement };
	const ZoomIn = (ev: MouseEvent & AndInfoEvent) => {
		const x = ev.clientX - ev.currentTarget.offsetLeft,
			y = ev.clientY - ev.currentTarget.offsetTop;

		ev.currentTarget.style.transformOrigin = `${x}px ${y}px`;
		ev.currentTarget.style.transform = "scale(2.5)";
	};
	const ZoomOut = (ev: MouseEvent & AndInfoEvent) => {
		ev.currentTarget.style.transformOrigin = "center";
		ev.currentTarget.style.transform = "scale(1)";
	};

	const PrevPic = () => {
		let newId = imgDrawer[2] - 1;
		if (newId < 0) newId = imgDrawer[1].length - 1;
		imgDrawer[2] = newId;
	};
	const NextPic = () => {
		let newId = imgDrawer[2] + 1;
		if (newId > imgDrawer[1].length - 1) newId = 0;
		imgDrawer[2] = newId;
	};
</script>

{#if imgDrawer[0]}
	<div
		class="fixed w-screen h-screen bg-black bg-opacity-60 top-0 left-0 gap-x-5 z-50 flex justify-center items-center"
	>
		<button class="btn btn-secondary" on:click={PrevPic}>
			<i class="fa-solid fa-backward" />
		</button>
		<img
			on:mousemove={ZoomIn}
			on:mouseleave={ZoomOut}
			src={imgDrawer[1][imgDrawer[2]]}
			alt="Zoomed img"
		/>
		<button class="btn btn-secondary" on:click={NextPic}>
			<i class="fa-solid fa-forward" />
		</button>

		<button on:click={closeImageDrawer} class="btn btn-primary absolute top-[80px] left-2"
			><i class="fa-solid fa-xmark text-xl" /></button
		>
	</div>
{/if}
