<script lang="ts">
	export let href: string;
	export let classStyle = "";
</script>

<a sveltekit:prefetch {href} class={classStyle}><slot /></a>
