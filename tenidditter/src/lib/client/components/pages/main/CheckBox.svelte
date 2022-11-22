<script lang="ts">
	import { onMount } from "svelte";

	let IconRef: HTMLDivElement;

	onMount(() => {
		if (!IconRef) return;

		let CheckboxObserver = new IntersectionObserver(AnimationIcon, { threshold: 1.0 });
		CheckboxObserver.observe(IconRef);
	});

	const AnimationIcon: IntersectionObserverCallback = ([{ isIntersecting }]) => {
		if (!IconRef) return;

		const IconEl = IconRef.children[0];
		if (isIntersecting) {
			IconRef?.classList.add("animate");

			{
				IconEl?.classList.remove("fa-check");
				IconEl?.classList.add("fa-xmark");
			}

			setTimeout(() => {
				IconRef?.classList.remove("animate");

				{
					IconEl?.classList.remove("fa-xmark");
					IconEl?.classList.add("fa-check");
				}
			}, Math.round(Math.random() * 3000) + 1000);
		}
	};
</script>

<div
	bind:this={IconRef}
	class="checkbox sm:scale-100 scale-75 animate flex justify-center items-center h-10 w-10 rounded-full"
>
	<i class="icon fas text-lg text-success" />
</div>

<style scoped>
	.checkbox::after {
		content: "";
		border: 2px dashed #fff;
		height: 2.5rem;
		width: 2.5rem;
		position: absolute;
		border-radius: 100%;
	}
	.checkbox.animate::after {
		animation: RotateBorder 4s linear 0.1s infinite forwards;
	}

	@keyframes RotateBorder {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}
</style>
