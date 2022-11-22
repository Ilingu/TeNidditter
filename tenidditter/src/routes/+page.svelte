<script>
	import Footer from "$lib/client/components/layout/Footer.svelte";
	import Background from "$lib/client/components/pages/main/Background.svelte";
	import Hero from "$lib/client/components/pages/main/Hero.svelte";
	import HeroBot from "$lib/client/components/pages/main/HeroBot.svelte";
	import Props from "$lib/client/components/pages/main/Props.svelte";
</script>

<Background />
<main class="w-screen mb-20">
	<!-- HERO -->
	<Hero />
	<HeroBot />
	<Props />
</main>
<Footer />
