<script lang="ts">
	import { onMount } from "svelte";
	import { InitWasm } from "$lib/services/wasm/wasm";

	import AlertProvider from "$lib/components/layout/AlertProvider.svelte";
	import Navbar from "$lib/components/layout/Navbar.svelte";

	import "../style/app.css";
	import { changeAppTheme } from "$lib/utils";
	import { afterNavigate } from "$app/navigation";
	import { AutoLogin } from "$lib/services/auth";

	afterNavigate((n) => {
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
		await InitWasm();
		AutoLogin();

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
