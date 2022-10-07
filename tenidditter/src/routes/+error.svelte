<script lang="ts">
	import { page } from "$app/stores";
	import Link from "$lib/components/design/Link.svelte";
</script>

<div class="page-content flex flex-col justify-center items-center text-center">
	<div class="glitch text-error text-5xl">
		<span aria-hidden="true">ERROR -- {$page.status}</span>
		ERROR -- {$page.status}
		<span aria-hidden="true">ERROR -- {$page.status}</span>
	</div>

	<p class="text-center text-xl">Oops! Something went wrong:</p>
	<p class="text-center text-xl font-bold mb-5">{$page.error?.message}</p>
	<Link href="/" classStyle="btn">Return Home</Link>
</div>

<style scoped>
	.glitch {
		user-select: none;
		font-weight: bold;
		text-transform: uppercase;
		position: relative;
		text-shadow: 0.05em 0 0 #1da1f2, -0.03em -0.04em 0 #ff4500, 0.025em 0.04em 0 #f5f8fa;
		animation: glitch 725ms infinite;
	}

	.glitch span {
		position: absolute;
		top: 0;
		left: 0;
	}

	.glitch span:first-child {
		animation: glitch 500ms infinite;
		clip-path: polygon(0 0, 100% 0, 100% 35%, 0 35%);
		transform: translate(-0.04em, -0.03em);
		opacity: 0.75;
	}

	.glitch span:last-child {
		animation: glitch 375ms infinite;
		clip-path: polygon(0 65%, 100% 65%, 100% 100%, 0 100%);
		transform: translate(0.04em, 0.03em);
		opacity: 0.75;
	}

	/* Turn off the animation for people who disable animation in theur OS's settings */
	@media (prefers-reduced-motion) {
		.glitch,
		.glitch span:first-child,
		.glitch span:last-child {
			animation: none;
			text-shadow: none;
			clip-path: none;
			transform: none;
		}
	}

	@keyframes glitch {
		0% {
			text-shadow: 0.05em 0 0 #1da1f2, -0.03em -0.04em 0 #ff4500, 0.025em 0.04em 0 #f5f8fa;
		}
		15% {
			text-shadow: 0.05em 0 0 #1da1f2, -0.03em -0.04em 0 #ff4500, 0.025em 0.04em 0 #f5f8fa;
		}
		16% {
			text-shadow: -0.05em -0.025em 0 #1da1f2, 0.025em 0.035em 0 #ff4500, -0.05em -0.05em 0 #f5f8fa;
		}
		49% {
			text-shadow: -0.05em -0.025em 0 #1da1f2, 0.025em 0.035em 0 #ff4500, -0.05em -0.05em 0 #f5f8fa;
		}
		50% {
			text-shadow: 0.05em 0.035em 0 #1da1f2, 0.03em 0 0 #ff4500, 0 -0.04em 0 #f5f8fa;
		}
		99% {
			text-shadow: 0.05em 0.035em 0 #1da1f2, 0.03em 0 0 #ff4500, 0 -0.04em 0 #f5f8fa;
		}
		100% {
			text-shadow: -0.05em 0 0 #1da1f2, -0.025em -0.04em 0 #ff4500, -0.04em -0.025em 0 #f5f8fa;
		}
	}
</style>
