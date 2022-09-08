<script>
	import Link from "$lib/components/design/Link.svelte";
</script>

<nav class="navbar bg-base-100 sticky top-0 z-10">
	<div class="flex-1 ml-[5%] flex gap-x-4">
		<Link href="/" classStyle="btn btn-ghost text-3xl font-nerd font-bold">
			<img src="/favicon.ico" alt="app logo" class="w-12" />
			TeNidditter
		</Link>
	</div>
	<Link href="/nitter" classStyle="font-semibold font-nerd flex gap-x-2 mr-[5%] items-center">
		<img src="/Assets/Img/twitter.webp" alt="app logo" class="w-6 hue-rotate-180" />
		Nitter
	</Link>
	<Link href="/teddit" classStyle="font-semibold font-nerd flex gap-x-2 mr-[5%] items-center">
		<img src="/Assets/Img/reddit.svg" alt="app logo" class="w-6 hue-rotate-180" />
		Teddit
	</Link>
</nav>
