<script lang="ts">
	import { onMount } from "svelte";
	import AlertProvider from "$lib/services/AlertProvider.svelte";
	import Navbar from "$lib/components/layout/Navbar.svelte";

	import "../style/app.css";
	import { AutoLogin } from "$lib/stores/auth";

	onMount(() => {
		AutoLogin();

		let ScrollAnimationObserver = new IntersectionObserver(ScrollAnimation, { threshold: 1.0 });
		document
			.querySelectorAll(".scrollAnimate")
			.forEach((el) => ScrollAnimationObserver.observe(el));
	});

	const ScrollAnimation: IntersectionObserverCallback = (entries) => {
		for (const { isIntersecting, target } of entries) {
			if (isIntersecting) (target as HTMLElement).classList.add("animate");
			// else (target as HTMLElement).classList.remove("animate");
		}
	};
</script>

<Navbar />
<slot />
<AlertProvider />
