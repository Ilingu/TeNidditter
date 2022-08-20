<script>
	import Links from "$lib/components/design/Links.svelte";
</script>

<nav class="navbar bg-base-100 flex justify-around">
	<div class="flex-1 ml-[5%] flex gap-x-20 ">
		<Links href="/" classStyle="btn btn-ghost text-3xl font-nerd font-bold">
			<img src="/favicon.ico" alt="app logo" class="" />
			TeNidditter
		</Links>
		<Links href="/nitter" classStyle="font-semibold font-nerd">Nitter</Links>
		<Links href="/teddit" classStyle="font-semibold font-nerd">Teddit</Links>
	</div>
</nav>
