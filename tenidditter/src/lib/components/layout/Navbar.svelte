<script>
	import Link from "$lib/components/design/Link.svelte";
	import AuthStore from "$lib/stores/auth";
	import { page } from "$app/stores";
	import ProfilePicture from "./ProfilePicture.svelte";
	import { LogOut } from "$lib/services/auth";

	let active = "";
	$: active =
		$page.url.pathname === "/"
			? "one"
			: $page.url.pathname.includes("/teddit")
			? "two"
			: $page.url.pathname.includes("/nitter")
			? "three"
			: $page.url.pathname.includes("/auth")
			? "four"
			: "";
</script>

<!-- Desktop Version -->
<nav class="navbar bg-base-200 sticky top-0 z-10 hidden md:flex">
	<div class="navbar-start">
		<Link href="/">
			<button class="btn btn-ghost lg:text-3xl text-2xl font-nerd font-bold">
				<img src="/favicon.ico" alt="app logo" class="w-12" />
				TeNidditter
			</button>
		</Link>
	</div>

	<div class="navbar-center gap-x-5">
		<Link href="/nitter">
			<button class="font-semibold font-nerd flex items-center">
				<img src="/Assets/Img/twitter.webp" alt="app logo" class="w-6 hue-rotate-180" />
				Nitter</button
			>
		</Link>
		<Link href="/teddit">
			<button class="font-semibold font-nerd flex items-center "
				><img src="/Assets/Img/reddit.svg" alt="app logo" class="w-6 hue-rotate-180" />
				Teddit</button
			>
		</Link>
	</div>

	<div class="navbar-end mr-5 gap-x-2">
		{#if $page.url.pathname.includes("/teddit")}
			<div class="dropdown dropdown-end">
				<label tabindex="0" for="" class="btn btn-ghost text-xl">
					<i class="fa-solid fa-magnifying-glass" />
				</label>
				<ul
					tabindex="0"
					class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-72"
				>
					<li>
						<Link href="/teddit/r"
							><button class="btn btn-primary gap-x-3 w-full"
								><i class="fa-solid fa-magnifying-glass" /> Search Subteddit</button
							></Link
						>
					</li>
					<li>
						<Link href="/teddit/u"
							><button class="btn btn-warning gap-x-3 w-full">
								<i class="fas fa-user" />
								Search Teddit's User</button
							></Link
						>
					</li>
				</ul>
			</div>
		{:else if $page.url.pathname.includes("/nitter")}
			<Link href="/nitter/search">
				<label tabindex="0" for="" class="btn btn-ghost text-xl">
					<i class="fa-solid fa-magnifying-glass" />
				</label></Link
			>
		{/if}
		{#if $AuthStore.loggedIn}
			<div class="dropdown dropdown-end">
				<label tabindex="0" for="" class="btn btn-ghost btn-circle avatar">
					<div class="flex justify-center avatar mb-1" title="Your Profile">
						<ProfilePicture size="medium" />
					</div>
				</label>
				<ul
					tabindex="0"
					class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
				>
					<li><Link href="/auth"><i class="fa-solid fa-gear" /> Settings</Link></li>
					<li>
						<button on:click={() => LogOut(true, $AuthStore.JwtToken)}
							><i class="fa-solid fa-right-from-bracket icon" /> Logout</button
						>
					</li>
				</ul>
			</div>
		{:else}
			<Link href="/auth">
				<button class="btn btn-primary btn-sm transition-all gap-x-3 text-lg"
					><i class="fa-solid fa-fire text-white" /> Get Started</button
				></Link
			>
		{/if}
	</div>
</nav>

<!-- Mobile Version -->
<div class={`btm-nav bg-base-200 ${active} flex md:hidden z-10`}>
	<Link href="/" className={$page.url.pathname === "/" ? "one" : ""}>
		<button class="text-primary">
			<i class="fa-solid fa-house" />
		</button>
	</Link>
	<Link href="/teddit" className={$page.url.pathname.includes("/teddit") ? "two" : ""}>
		<button class="text-teddit active">
			<i class="fa-brands fa-reddit-alien" />
		</button>
	</Link>
	<Link href="/nitter" className={$page.url.pathname.includes("/nitter") ? "three" : ""}>
		<button class="text-nitter">
			<i class="fa-brands fa-twitter" />
		</button>
	</Link>
	<Link href="/auth" className={$page.url.pathname.includes("/auth") ? "four" : ""}>
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
