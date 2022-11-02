<script>
	import Link from "$lib/components/design/Link.svelte";
	import AuthStore from "$lib/stores/auth";
	import { page } from "$app/stores";

	let active = "";
	$: active =
		$page.url.pathname === "/"
			? "one"
			: $page.url.pathname === "/teddit"
			? "two"
			: $page.url.pathname === "/nitter"
			? "three"
			: $page.url.pathname === "/auth"
			? "four"
			: "";
</script>

<!-- Desktop Version -->
<nav class="navbar bg-base-100 sticky top-0 z-10 hidden md:flex">
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

<!-- Mobile Version -->
<div class={`btm-nav bg-base-200 ${active} flex md:hidden z-10`}>
	<Link href="/" className={$page.url.pathname === "/" ? "one" : ""}>
		<button class="text-primary">
			<i class="fa-solid fa-house" />
		</button>
	</Link>
	<Link href="/teddit" className={$page.url.pathname.includes("/teddit") ? "two" : ""}>
		<button class="text-primary active">
			<i class="fa-brands fa-reddit-alien" />
		</button>
	</Link>
	<Link href="/nitter" className={$page.url.pathname.includes("/nitter") ? "three" : ""}>
		<button class="text-primary">
			<i class="fa-brands fa-twitter" />
		</button>
	</Link>
	<Link href="/auth" className={$page.url.pathname === "/auth" ? "four" : ""}>
		<button class="text-primary">
			<i class="fa-solid fa-user" />
		</button>
	</Link>
</div>

<style scoped>
	.btm-nav::after {
		content: "";
		position: absolute;
		top: 0;
		left: 0;
		width: 25%;
		height: 2px;
		background-color: rgb(220, 165, 76);
		transition: 0.5s all;
	}

	.btm-nav.one::after {
		left: 0;
	}
	.btm-nav.two::after {
		left: 25%;
	}
	.btm-nav.three::after {
		left: 50%;
	}
	.btm-nav.four::after {
		left: 75%;
	}
</style>
