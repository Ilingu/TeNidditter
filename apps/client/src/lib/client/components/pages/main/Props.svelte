<script lang="ts">
	import CheckBox from "$lib/client/components/pages/main/CheckBox.svelte";

	const PropsText = [
		"<strong>No</strong> JavaScript or ads",
		"All requests go through the backend",
		"Client <strong>never talks</strong> to Twitter/Reddit... (Proxy)",
		"Prevents your <strong>IP/Fingerprint tracking</strong>",
		"<strong>Lightweight && Fast</strong> (SvelteKit + Golang backend)",
		"Mobile support (responsive design)",
		`MIT && <a href="https://github.com/Ilingu/TeNidditter" target="_blank" rel="noopener noreferrer"><strong>Open Source</strong></a>`
	];
</script>

<section id="Props" class="w-screen flex justify-center">
	<ul class="sm:text-justify sm:text-xl text-lg flex flex-col gap-y-4">
		{#each PropsText as text}
			<li class="scrollAnimate flex items-center gap-x-2 sm:leading-normal leading-5">
				<CheckBox />
				<div>
					{@html text}
				</div>
			</li>
		{/each}
	</ul>
</section>
