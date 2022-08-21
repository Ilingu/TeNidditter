<script context="module">
	import Background from "$lib/components/pages/main/Background.svelte";
	import Hero from "$lib/components/pages/main/Hero.svelte";

	export const prerender = true;
</script>

<Background />
<main class="w-screen h-[500vh]">
	<!-- HERO -->
	<Hero />
</main>
