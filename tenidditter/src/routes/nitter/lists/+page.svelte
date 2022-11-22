<script lang="ts">
	import api from "$lib/api";
	import Link from "$lib/client/components/design/Link.svelte";
	import AuthStore from "$lib/client/stores/auth";
	import { IsEmptyString, MakeBearerToken, pushAlert, TrimSpecialChars } from "$lib/utils";
	import { fade } from "svelte/transition";

	let CloseModalBtn: HTMLLabelElement;
	let loading = false;

	let listname = "";
	const CreateNewList = async () => {
		loading = true;
		if (IsEmptyString(listname) || listname.length < 3 || listname.length >= 30)
			return pushAlert("listaname should be between 3 and 30 characters", "warning", 8000);

		const { success: listCreated } = await api.post("/tedinitter/nitter/list", {
			headers: MakeBearerToken($AuthStore.JwtToken ?? ""),
			body: { listname }
		});

		if (!listCreated) pushAlert("Couldn't create your list", "error");
		else CloseModalBtn?.click();

		loading = false;
	};

	const DeleteList = async (listId: string | number) => {
		listId = TrimSpecialChars(listId.toString());

		if (IsEmptyString(listId)) return;
		const { success: deleted } = await api.delete("/tedinitter/nitter/list/%s", {
			params: [listId],
			headers: MakeBearerToken($AuthStore.JwtToken ?? "")
		});
		if (!deleted) pushAlert("Couldn't delete this list", "error");
	};
</script>

<main class="grid place-items-center mt-5">
	<h1 class="text-3xl font-bold flex flex-wrap items-center justify-center gap-x-4 gap-y-2">
		<i class="fa-solid fa-list" /> Manage Your
		<span
			class="from-primary to-secondary bg-gradient-to-br bg-clip-text text-clip text-transparent"
			>Lists</span
		>
		<!-- The button to open modal -->
		<label for="modal-new-list" class="btn gap-x-2"
			><i class="fa-solid fa-square-plus icon" /> New List</label
		>
	</h1>
	{#if $AuthStore?.Lists && $AuthStore.Lists.length > 0}
		<div class="sm:w-[400px] w-[375px] flex flex-col gap-y-3 mt-4">
			{#each $AuthStore.Lists as list, i}
				<div
					in:fade={{ delay: 100 * i }}
					class="flex gap-x-2 group items-center h-14 bg-neutral rounded-md w-full p-2"
				>
					<svg
						width="36"
						height="36"
						viewBox="0 0 24 24"
						xmlns="http://www.w3.org/2000/svg"
						fill-rule="evenodd"
						clip-rule="evenodd"
						class="fill-current"
						><path
							d="M22.672 15.226l-2.432.811.841 2.515c.33 1.019-.209 2.127-1.23 2.456-1.15.325-2.148-.321-2.463-1.226l-.84-2.518-5.013 1.677.84 2.517c.391 1.203-.434 2.542-1.831 2.542-.88 0-1.601-.564-1.86-1.314l-.842-2.516-2.431.809c-1.135.328-2.145-.317-2.463-1.229-.329-1.018.211-2.127 1.231-2.456l2.432-.809-1.621-4.823-2.432.808c-1.355.384-2.558-.59-2.558-1.839 0-.817.509-1.582 1.327-1.846l2.433-.809-.842-2.515c-.33-1.02.211-2.129 1.232-2.458 1.02-.329 2.13.209 2.461 1.229l.842 2.515 5.011-1.677-.839-2.517c-.403-1.238.484-2.553 1.843-2.553.819 0 1.585.509 1.85 1.326l.841 2.517 2.431-.81c1.02-.33 2.131.211 2.461 1.229.332 1.018-.21 2.126-1.23 2.456l-2.433.809 1.622 4.823 2.433-.809c1.242-.401 2.557.484 2.557 1.838 0 .819-.51 1.583-1.328 1.847m-8.992-6.428l-5.01 1.675 1.619 4.828 5.011-1.674-1.62-4.829z"
						/>
					</svg>
					<div
						class="flex-1 hover:underline hover:underline-offset-2 hover:decoration-wavy hover:decoration-accent"
					>
						<Link href={`/nitter/list/${list.list_id}`}>
							<span
								class=" from-accent to-secondary bg-gradient-to-br bg-clip-text text-clip text-transparent text-bold text-xl capitalize"
								>{list.title}</span
							>
						</Link>
					</div>
					<button
						on:click={() => DeleteList(list.list_id)}
						class="opacity-0 group-hover:opacity-100 z-50 transition-opacity btn btn-sm btn-error"
						>Delete</button
					>
				</div>
			{/each}
		</div>
	{:else}
		<h2 class="text-xl mt-4">Nothing to show! Create a list first.</h2>
	{/if}
</main>

<!-- Put this part before </body> tag -->
<input type="checkbox" id="modal-new-list" class="modal-toggle" />
<div class="modal">
	<div class="modal-box">
		<h3 class="font-bold text-lg">Create a new list!</h3>
		<p class="py-4">Please fill the following fields:</p>
		<form on:submit|preventDefault={CreateNewList}>
			<div class="form-control">
				<label class="input-group">
					<span>Name</span>
					<input
						type="text"
						placeholder="my-science-list"
						bind:value={listname}
						on:input={(ev) => (listname = TrimSpecialChars(ev.currentTarget?.value))}
						class="input input-bordered input-accent"
					/>
				</label>
			</div>
		</form>
		<div class="modal-action">
			<label for="modal-new-list" class={`btn btn-warning gap-x-2 ${loading ? "disabled" : ""}`}
				><i class="fa-solid fa-ban" /> Cancel</label
			>
			<button class="btn btn-info gap-x-2" on:click={CreateNewList} disabled={loading}
				><i class="fa-regular fa-square-plus" /> Submit!</button
			>
			<label for="modal-new-list" class="hidden" bind:this={CloseModalBtn} />
		</div>
	</div>
</div>

<style scoped>
</style>
