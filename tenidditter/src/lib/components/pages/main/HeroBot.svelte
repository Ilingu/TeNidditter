<script>
	import HomeCard from "$lib/components/design/HomeCard.svelte";
</script>

<div class="min-h-[30vh] text-center pt-7">
	<h1 class="text-3xl tracking-wider">
		Privacy <span class="text-accent-content">Matters</span>
	</h1>

	<div class="flex justify-evenly mt-6">
		<HomeCard title="E2E Encryption" icon="fas fa-lock" />
		<HomeCard
			title="<a href='https://github.com/Ilingu/TeNidditter' target='_blank' rel='noopener noreferrer'>Open Source</a>"
			icon="fa-brands fa-github"
		/>
	</div>
</div>
