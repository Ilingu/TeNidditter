<script lang="ts">
	import { afterNavigate } from "$app/navigation";
	import { onMount } from "svelte";

	const handleAnimatePlacehover = () => {
		document.querySelectorAll("input").forEach((inp) => {
			inp.addEventListener("focusin", () => PlacehoverAnimate(inp, "add"));
			inp.addEventListener("focusout", () => {
				if (inp.value.trim().length >= 1) return;
				PlacehoverAnimate(inp, "rem");
			});
		});
	};
	onMount(handleAnimatePlacehover);
	afterNavigate(handleAnimatePlacehover);

	/* ANIMATION */
	const PlacehoverAnimate = (inp: HTMLInputElement, mode: "add" | "rem") => {
		const label = inp.parentElement?.parentElement?.firstChild as HTMLElement;
		const hasAnimate = label?.classList.contains("animate");

		if (mode === "add" && hasAnimate) return;
		if (mode === "rem" && !hasAnimate) return;

		label?.classList.toggle("animate");
	};
</script>

<slot />

<style scoped>
	:global(input[disabled]) {
		opacity: 0.5;
	}

	:global(.placehover) {
		transition: all 0.5s;
	}

	:global(.placehover.animate) {
		translate: 0 -40px;
	}
</style>
