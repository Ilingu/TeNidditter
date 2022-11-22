<script lang="ts">
	import Link from "$lib/client/components/design/Link.svelte";
	import type { NittosPreview } from "$lib/shared/types/nitter";
	import { IsEmptyString, isValidUrl } from "$lib/shared/utils";
	import { onMount } from "svelte";

	export let previewDatas: NittosPreview[];

	let feedDiv: HTMLDivElement;
	onMount(() => {
		if (!feedDiv) return;

		feedDiv.querySelectorAll(".nittosPreviewContainer a").forEach((a) => {
			if (a.classList.contains("not")) return;
			const href = a.getAttribute("href");

			if (href?.startsWith("http")) {
				if (isValidUrl(a.textContent || "")) a.textContent = new URL(a.textContent || "").host;
				a.innerHTML = `<i class="fa-solid fa-arrow-up-right-from-square icon"></i>` + a.innerHTML;

				a.setAttribute("target", "_blank");
				a.setAttribute("rel", "noopener noreferrer");
			} else a.setAttribute("href", "/nitter" + href);
		});
	});
</script>

<div bind:this={feedDiv} class="nittosPreviewContainer flex flex-col gap-y-2 md:mx-0 mx-2">
	{#each previewDatas as nittos}
		{#if !IsEmptyString(nittos.username)}
			<div
				class="grid gap-x-2 grid-cols-10 bg-neutral-focus p-2 rounded-lg hover:bg-neutral transition-all"
			>
				<div class="col-span-2 flex justify-center items-center">
					<img src={nittos.avatarUrl} alt="" class="w-fit h-fit rounded-full" />
				</div>
				<div class="col-span-8">
					<Link href={`/nitter/${nittos.username.slice(1)}`} className="not">{nittos.username}</Link
					>
					<p>{@html nittos.description}</p>
				</div>
			</div>
		{/if}
	{/each}
</div>

<style scoped>
	:global(.nittosPreviewContainer a) {
		color: #f3cc30;
		font-weight: bold;
		transition: 0.3s all;
	}
	:global(.nittosPreviewContainer a:hover) {
		text-decoration: underline wavy;
		color: #58c7f3;
	}
</style>
