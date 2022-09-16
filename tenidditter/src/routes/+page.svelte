<script>
	import Background from "$lib/components/pages/main/Background.svelte";
	import Hero from "$lib/components/pages/main/Hero.svelte";
	import HeroBot from "$lib/components/pages/main/HeroBot.svelte";
	import Props from "$lib/components/pages/main/props.svelte";
</script>

<Background />
<main class="w-screen h-[500vh]">
	<!-- HERO -->
	<Hero />
	<HeroBot />
	<Props />
</main>
