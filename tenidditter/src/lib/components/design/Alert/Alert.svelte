<script lang="ts">
	import type { AlertShape } from "$lib/types/interfaces";
	import AlertChild from "./AlertChild.svelte";

	export let alert: Omit<AlertShape, "duration">;
</script>

<AlertChild type={alert.type}>{alert.message}</AlertChild>
