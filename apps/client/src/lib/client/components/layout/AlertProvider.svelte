<script lang="ts">
	import Alert from "$lib/client/components/design/Alert/Alert.svelte";
	import type { AlertShape } from "$lib/client/types/alerts";
	import { IsEmptyString } from "$lib/shared/utils";

	import { onMount } from "svelte";

	type EnhancedAlert = AlertShape & { id: string };
	let alerts: EnhancedAlert[] = [];

	onMount(() => {
		document.addEventListener("alert", handleAlert as EventListener);
	});

	const handleAlert = ({ detail: alertEvt }: CustomEvent<AlertShape>) => {
		if (!alertEvt) return;
		if (IsEmptyString(alertEvt.message)) return;

		const alert: EnhancedAlert = {
			id: crypto.randomUUID(),
			...alertEvt
		};
		alerts = [...alerts, alert];

		setTimeout(() => {
			alerts = alerts.filter((currAlert) => currAlert.id !== alert.id);
		}, alert?.duration || 5000);
	};
</script>

<div class="fixed z-[1000] top-20 right-2 xl:w-1/4 sm:w-1/2 w-3/4 flex flex-col gap-y-2">
	{#each alerts as alert}
		<Alert {alert} />
	{/each}
</div>
