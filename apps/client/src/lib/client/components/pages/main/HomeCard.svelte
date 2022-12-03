<script lang="ts">
	export let icon: string;
	export let title: string;
</script>

<div
	class="scrollAnimate sm:scale-100 scale-75 flex gap-y-2 flex-col items-center group hover:scale-105 transition-all"
>
	<div class="h-20 w-20 bg-gray-100 bg-opacity-20 rounded-md flex justify-center items-center ">
		<span class={icon + " text-xl group-hover:scale-110 transition-all text-accent-content"} />
	</div>
	<p
		class="sm:text-xl text-lg font-nerd font-semibold tracking-wider text-primary group-hover:text-accent-content"
	>
		{@html title}
	</p>
</div>
