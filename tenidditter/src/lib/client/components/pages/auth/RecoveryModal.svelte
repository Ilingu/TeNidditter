<script lang="ts">
	import { copyToClipboard, pushAlert } from "$lib/client/ClientUtils";

	export let show = false;
	export let setShow: (newVal: boolean) => void;

	let recoverySaved = false;

	/* recoveryCodes */
	export let recoveryCodes: string[] = [];
	const DownloadRecoveryCodes = () => {
		if (recoveryCodes.length <= 0) return;
		const codes = recoveryCodes.join("\n");
		const textBlob = new Blob([codes], { type: "text/plain" });
		const fileUrl = URL.createObjectURL(textBlob);

		var a = document.createElement("a");
		a.href = fileUrl;
		a.download = "TeNiditter_Recovery_Codes.txt";
		a.click();
		recoverySaved = true;
	};
</script>

<!-- recovery codes Modal -->
<input type="checkbox" id="recoveryCodeModal" class="modal-toggle" bind:checked={show} />
<div class="modal">
	<div
		class="modal-box bg-neutral-focus flex flex-col items-center text-center w-[37.5%] max-w-[100vw]"
	>
		<h3 class="font-bold text-2xl">Don't forget to save your recovery codes!</h3>
		<p class="py-4 text-lg">
			If you loose your password you can use these codes to recover your account by changing of
			password. <br />
			<span class="text-error font-bold"
				>If you loose them too, consider your account lost forever.</span
			>
			<br />
			<span class="text-info">
				<i class="fa-solid fa-info" /> Please Note that this is the only chance to download them, they
				won't be display another time!
			</span>
		</p>

		<div class="bg-base-100 w-3/4 grid grid-cols-2 gap-y-3 rounded-md p-5">
			{#each recoveryCodes as code}
				<p class="text-white font-bold text-xl select-all">{code}</p>
			{/each}
		</div>

		<div class="flex gap-x-3 mt-2">
			<button class="btn btn-accent gap-x-2" on:click={DownloadRecoveryCodes}
				><i class="fa-solid fa-download icon" /> Download</button
			>
			<a
				class="btn btn-accent gap-x-2"
				href={`/api/recovery-code/print?codes=${encodeURIComponent(JSON.stringify(recoveryCodes))}`}
				target="_blank"
				on:click={() => (recoverySaved = true)}
				rel="noopener noreferrer"><i class="fa-solid fa-print icon" /> Print</a
			>
			<button
				class="btn btn-accent gap-x-2"
				on:click={() => {
					copyToClipboard(recoveryCodes.join("\n"));
					recoverySaved = true;
					pushAlert("Recovery Codes Copied", "info", 1600);
				}}><i class="fa-solid fa-copy icon" /> Copy</button
			>
		</div>
		<div class="modal-action">
			<button
				disabled={!recoverySaved}
				class="btn btn-success"
				on:click={() => {
					recoverySaved = false;
					recoveryCodes = [];
					setShow(false);
				}}>💯 All good!</button
			>
		</div>
	</div>
</div>
