<script>
	import Footer from "$lib/components/layout/Footer.svelte";
	import Background from "$lib/components/pages/main/Background.svelte";
	import Hero from "$lib/components/pages/main/Hero.svelte";
	import HeroBot from "$lib/components/pages/main/HeroBot.svelte";
	import Props from "$lib/components/pages/main/Props.svelte";
</script>

<Background />
<main class="w-screen mb-20">
	<!-- HERO -->
	<Hero />
	<HeroBot />
	<Props />
</main>
<Footer />
