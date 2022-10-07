<script lang="ts">
	import { onMount } from "svelte";
	import { AutoLogin } from "$lib/stores/auth";
	import { InitWasm } from "$lib/wasm";

	import AlertProvider from "$lib/AlertProvider.svelte";
	import Navbar from "$lib/components/layout/Navbar.svelte";

	import "../style/app.css";
	import { changeAppTheme } from "$lib/utils";
	import { afterNavigate } from "$app/navigation";

	afterNavigate((n) => {
		const path = n.to?.url.pathname;
		if (!path) return;

		if (path.includes("nitter")) changeAppTheme("nitter");
		else if (path.includes("teddit")) changeAppTheme("teddit");
		else changeAppTheme("tenidditter");
	});

	onMount(async () => {
		await InitWasm();
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
