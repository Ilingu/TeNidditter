<script lang="ts">
	import type { AlertShape } from "$lib/client/types/alerts";
	import AlertChild from "./AlertChild.svelte";

	export let alert: Omit<AlertShape, "duration">;
</script>

<AlertChild type={alert.type}>{alert.message}</AlertChild>
