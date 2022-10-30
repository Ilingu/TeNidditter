<script>
	import Link from "$lib/components/design/Link.svelte";
	import { fly } from "svelte/transition";

	let opened = false;
</script>

<slot />
<div class="flex flex-col fixed bottom-0 left-0 z-10 gap-y-3 h-[23vh]">
	<label class="btn btn-circle btn-primary m-5 text-white swap swap-rotate absolute bottom-0 z-10">
		<!-- this hidden checkbox controls the state -->
		<input type="checkbox" on:change={() => (opened = !opened)} />

		<!-- hamburger icon -->
		<svg
			class="swap-off fill-current"
			xmlns="http://www.w3.org/2000/svg"
			width="32"
			height="32"
			viewBox="0 0 512 512"
			><path d="M64,384H448V341.33H64Zm0-106.67H448V234.67H64ZM64,128v42.67H448V128Z" /></svg
		>

		<!-- close icon -->
		<svg
			class="swap-on fill-current"
			xmlns="http://www.w3.org/2000/svg"
			width="32"
			height="32"
			viewBox="0 0 512 512"
			><polygon
				points="400 145.49 366.51 112 256 222.51 145.49 112 112 145.49 222.51 256 112 366.51 145.49 400 256 289.49 366.51 400 400 366.51 289.49 256 400 145.49"
			/></svg
		>
	</label>

	{#if opened}
		<div
			transition:fly={{ x: -200, duration: 500 }}
			class="w-96 h-full flex flex-col gap-y-5 rounded-lg p-5 bg-[#2e2d2f]"
		>
			<Link href={"/teddit/r"}
				><button class="btn btn-active btn-primary gap-x-3 w-full"
					><i class="fa-solid fa-magnifying-glass" /> Search Subteddit</button
				></Link
			>
			<Link href={"/teddit/u"}
				><button class="btn btn-active btn-warning gap-x-3 w-full">
					<i class="fas fa-user" />
					Search Teddit's User</button
				></Link
			>
		</div>
	{/if}
</div>
