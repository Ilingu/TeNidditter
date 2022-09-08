<script lang="ts">
	import { routingListener } from "$lib/utils/routing";

	import { isMobile } from "$lib/utils/utils";

	import { onMount } from "svelte";
	import { fade } from "svelte/transition";

	interface Image {
		x: number;
		y: number;
		z: number;

		duration: number;
		name: "reddit" | "twitter";
	}

	let ViewportHeight = 1080,
		ViewportWidth = 1920;
	let images: Image[] = [];
	let AnimationMode: "Move" | "Ping" | "stop" = "Move";
	let AnimationReq: number;

	onMount(() => {
		if (isMobile()) return;

		ViewportHeight = innerHeight;
		ViewportWidth = innerWidth;

		images = Array.from(
			{ length: 35 },
			(_, i): Image => ({
				name: i % 2 == 0 ? "reddit" : "twitter",
				x: Math.round(Math.random() * ViewportWidth),
				y: Math.round(Math.random() * ViewportHeight),
				z: parseFloat((Math.random() + 0.5).toFixed(2)),
				duration: -1
			})
		);

		ObserveAnimate();
		AnimationReq = requestAnimationFrame(AnimateBg);

		// Unblock routing by stoping animation
		routingListener(() => {
			AnimationMode = "stop";
			cancelAnimationFrame(AnimationReq);
		});
	});

	const Move = () => {
		images = images.map((image, i) => {
			const { x, y, z } = image;

			let NewYCoord: number = y + 0.5 * z;
			let NewXCoord = i % 2 == 0 ? x - 0.1 : x + 0.1;

			if (NewYCoord > ViewportHeight || NewYCoord <= -1) NewYCoord = 0;
			if (NewXCoord > ViewportWidth || NewXCoord <= -1)
				NewXCoord = Math.round(Math.random() * ViewportWidth);

			return {
				...image,
				y: NewYCoord,
				x: x + Math.cos((NewYCoord * 2 * Math.PI) / 360) * parseFloat(Math.random().toFixed(1))
			};
		});
	};

	const Ping = () => {
		images = images.map((image) => {
			if (image.duration === -1) return { ...image, duration: Math.round(Math.random() * 1500) };
			if (image.duration > 0)
				return {
					...image,
					duration: image.duration - 1,
					z: Math.cos((image.duration * 2 * Math.PI) / 360) + 1
				};
			return {
				...image,
				y: Math.round(Math.random() * ViewportHeight),
				x: Math.round(Math.random() * ViewportWidth),
				duration: Math.round(Math.random() * 1500)
			};
		});
	};

	let lastSrTo = 0;
	const AnimateOnScroll = () => {
		const ScTo = window.pageYOffset;
		images = images.map((image) => ({
			...image,
			y: ScTo > lastSrTo ? image.y - 0.2 : image.y + 0.2
		}));

		lastSrTo = ScTo;
	};

	const AnimateBg = () => {
		switch (AnimationMode) {
			case "stop":
				document.removeEventListener("scroll", AnimateOnScroll);
				return;
			case "Move":
				document.removeEventListener("scroll", AnimateOnScroll);
				Move();
			case "Ping":
				document.addEventListener("scroll", AnimateOnScroll);
				Ping();
		}
		AnimationReq = requestAnimationFrame(AnimateBg);
	};

	const ObserveAnimate = () => {
		const observer = new IntersectionObserver(([{ isIntersecting: isVisible }]) => {
			if (!isVisible) AnimationMode = "Ping";
			else AnimationMode = "Move";
		});
		observer.observe(document.getElementById("HeadlineHero")!);
	};
</script>

<section class="w-screen fixed select-none">
	{#each images as { x, y, z, name }}
		<div
			transition:fade
			style={`top: ${y}px; left: ${x}px; transform: scale(${z}); opacity: ${z}; background-color: ${
				name === "reddit" ? "#FF4500" : "#1DA1F2"
			}`}
			class="absolute w-0.5 h-0.5"
		/>
	{/each}
</section>
