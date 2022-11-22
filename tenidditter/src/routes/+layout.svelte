<script lang="ts">
	import { onMount } from "svelte";

	import AlertProvider from "$lib/client/components/layout/AlertProvider.svelte";
	import Navbar from "$lib/client/components/layout/Navbar.svelte";

	import "../style/app.css";
	import { afterNavigate } from "$app/navigation";
	import { changeAppTheme } from "$lib/client/ClientUtils";
	import { InitWasm } from "$lib/client/services/wasm/wasm";
	import { AutoLogin } from "$lib/client/services/auth";

	afterNavigate((n) => {
		// Change Theme for each new routes
		const path = n.to?.url.pathname;
		if (!path) return;

		if (path.includes("nitter")) changeAppTheme("nitter");
		else if (path.includes("teddit")) changeAppTheme("teddit");
		else changeAppTheme("tenidditter");

		// ReTrigger All Animation
		document.querySelectorAll(".scrollAnimate").forEach((el) => {
			ScrollAnimationObserver.observe(el);
		});
	});

	let ScrollAnimationObserver: IntersectionObserver;
	onMount(async () => {
		await InitWasm(); // inject wasm into dom
		AutoLogin(); // check the login state

		ScrollAnimationObserver = new IntersectionObserver(ScrollAnimation, {
			threshold: innerWidth > 1280 ? 1.0 : undefined
		});
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
