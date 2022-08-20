<script lang="ts">
	import type { AlertShape } from "$lib/types/interfaces";
	import AlertChild from "./AlertChild.svelte";

	export let alert: AlertShape;
</script>

<AlertChild type={alert.type}>{alert.message}</AlertChild>
