<script lang="ts">
	export let data: import("./$types").PageData;
	console.log(data);
</script>

<main class="max-w-[1000px] m-auto flex justify-center gap-x-8 px-2 py-5">
	<header class="bg-[#2e2d2f] w-full p-5 rounded-lg h-[25vh]">
		<h1 class="text-2xl">@ {data.username}</h1>
	</header>
	<div class="max-w-[750px]" />
</main>
