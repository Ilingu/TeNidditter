<script lang="ts">
	import { onMount } from "svelte";

	const HeadlineText = [
		"Privacy comes first.",
		"Encrypted, private, secure.",
		"❤️ = 🔒 + 🕊",
		"Social Medias with freedom."
	];
	const DescText = [
		`Are you looking for an more privacy friendly alternative to <span class="text-teddit">reddit</span> or <span class="text-nitter">twitter</span>?`,
		`Tenidditter app is what social medias should have become <span class="font-bold">long ago</span>.`,
		`Tenidditter app act as a proxy between twitter/reddit's datas and you, you never directly speak to them, therefore you can keep <span class="font-bold">entertaining yourself without seeing your freedom being stolen</span>.`
	];

	/* ANIMATION */
	const DrawLines = () => {
		const canvas = document.getElementById("LinesCanvas") as HTMLCanvasElement;
		if (!canvas?.getContext) return;
		const ctx = canvas.getContext("2d");
		if (!ctx) return;

		canvas.height *= 4;
		canvas.width *= 5;

		const smoothness = Math.round(Math.random() * 24) + 1;
		const linesNumbers = Math.round(Math.random() * 50) + 50;

		for (let lineID = 0; lineID < linesNumbers; lineID++) {
			const h = Math.round(Math.random() * canvas.height);
			// const a = Math.round(Math.random()) + 0.5;
			// const f = Math.round(Math.random() * 2) + 1;

			// Filled sinwave
			ctx.beginPath();
			ctx.moveTo(0, h);
			let lastCoord: [number, number] = [0, h];
			for (let i = 0; i < canvas.width; i += smoothness) {
				let y = Math.sin(((i % 360) * 2 * Math.PI) / 360); // Calculate y value from x
				ctx.moveTo(i, lastCoord[1]); // Where to start drawing
				ctx.lineTo(lastCoord[0], h + 25 * y); // Where to draw to
				lastCoord = [i, h + 25 * y];
			}

			ctx.strokeStyle = Math.round(Math.random()) === 0 ? "#FF4500" : "#1DA1F2";
			ctx.stroke();
		}
	};

	/* App Start */
	onMount(DrawLines);
</script>

<aside class="col-span-4 relative">
	<canvas id="LinesCanvas" class="w-full h-full" />
	<div class="quote mockup-window border bg-base-300 absolute bottom-0 w-3/4">
		<div class="font-mono px-4 pt-2 h-32 bg-base-200">
			<h2 class="font-nerd text-2xl font-semibold text-primary">
				{HeadlineText[Math.round(Math.random() * (HeadlineText.length - 1))]}
			</h2>
			<p class="mt-3">{@html DescText[Math.round(Math.random() * (DescText.length - 1))]}</p>
		</div>
	</div>
</aside>

<style scoped>
	.quote {
		left: calc(50% - 75% / 2);
		animation: PopIn 1s cubic-bezier(0.6, 0, 0.96, 0.68) 0.1s 1 forwards;
	}
	@keyframes PopIn {
		from {
			translate: 0 400px;
		}
		to {
			translate: 0 -25%;
		}
	}
</style>
