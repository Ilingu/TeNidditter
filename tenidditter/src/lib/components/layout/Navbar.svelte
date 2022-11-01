<script>
	import Link from "$lib/components/design/Link.svelte";
	import AuthStore from "$lib/stores/auth";
</script>

<nav class="navbar bg-base-100 sticky top-0 z-10">
	<div class="mx-[2.5%] flex gap-x-4">
		<Link href="/">
			<button class="btn btn-ghost text-3xl font-nerd font-bold">
				<img src="/favicon.ico" alt="app logo" class="w-12" />
				TeNidditter
			</button>
		</Link>
	</div>
	<div class="gap-x-2 mr-[2.5%]">
		<Link href="/nitter">
			<button class="font-semibold font-nerd flex items-center">
				<img src="/Assets/Img/twitter.webp" alt="app logo" class="w-6 hue-rotate-180" />
				Nitter</button
			>
		</Link>
	</div>
	<div class="flex-1 gap-x-2">
		<Link href="/teddit">
			<button class="font-semibold font-nerd flex items-center "
				><img src="/Assets/Img/reddit.svg" alt="app logo" class="w-6 hue-rotate-180" />
				Teddit</button
			>
		</Link>
	</div>

	<Link href="/auth">
		{#if $AuthStore.loggedIn}
			<p class="gap-x-2 font-nerd text-xl">
				<i class="fas fa-user" />
				{$AuthStore.user?.username}
			</p>
		{:else}
			<button class="btn btn-primary btn-sm transition-all gap-x-3 text-lg"
				><i class="fa-solid fa-fire text-white" /> Get Started</button
			>
		{/if}
	</Link>
</nav>
