<script lang="ts">
	import Link from "$lib/client/components/design/Link.svelte";
	import ProfilePicture from "$lib/client/components/layout/ProfilePicture.svelte";
	import RecoveryModal from "$lib/client/components/pages/auth/RecoveryModal.svelte";
	import AuthStore from "$lib/client/stores/auth";
	import { pushAlert } from "$lib/utils";

	let showRecoveryCodesModal = false;
	const setShowRecoveryCodesModal = (newVal: boolean) => {
		showRecoveryCodesModal = newVal;
	};

	let newRecoveryCodes: string[] = [];

	const RegenerateCodes = async () => {
		const { success, data: newCodes } =
			await $AuthStore?.user.action?.regenerateUserRecoveryCodes();

		if (!success || !newCodes) return pushAlert("Failed to regenerate recovery codes", "error");

		newRecoveryCodes = newCodes;
		showRecoveryCodesModal = true;
	};

	/* DELETE Account stuff */
	enum DeleteStates {
		NoAction,
		Confirm,
		DeleteAccount
	}
	let deleteState = DeleteStates.NoAction;
	const HandleDeleteState = () => {
		if (deleteState >= DeleteStates.DeleteAccount) return CancelDeleteTimeout();

		canceled = false;
		deleteState++;

		if (deleteState === DeleteStates.DeleteAccount) HandleDeleteAccount();
	};

	let deleteTimeout: NodeJS.Timeout,
		canceled = false;
	const HandleDeleteAccount = () => {
		pushAlert("10s before your account deletion", "warning", 8_000);
		clearTimeout(deleteTimeout);
		deleteTimeout = setTimeout(DeleteAccount, 10_000);
	};
	const CancelDeleteTimeout = () => {
		canceled = true;
		clearTimeout(deleteTimeout);
		deleteState = DeleteStates.NoAction;
		pushAlert("You've successfully canceled your account deletion", "success", 5000);
	};

	const DeleteAccount = async () => {
		if (canceled) {
			ResetDeleteState();
			return pushAlert("You've successfully canceled your account deletion", "success", 5000);
		}
		$AuthStore?.user.action?.deleteAccount();
	};

	const ResetDeleteState = () => {
		deleteState = DeleteStates.NoAction;
		canceled = false;
	};
</script>

<main class="page-content w-screen flex justify-center">
	<aside
		class="bg-base-300 xl:max-w-[350px] max-w-[750px] min-w-[280px] w-[95%] p-5 rounded-lg h-fit mt-5"
	>
		<div class="flex justify-center avatar mb-1">
			<ProfilePicture size="big" />
		</div>
		<h1 class="text-3xl mb-4 text-center text-teddit font-semibold tracking-wide">
			{$AuthStore.user?.username}
		</h1>
		<ul>
			<li>
				<i class="fa-solid fa-id-card" /> ID:
				<span class="font-bold text-white">#{$AuthStore.user?.id}</span>
			</li>
			<li>
				<Link href="/nitter/lists">
					<i class="fa-solid fa-eye" /> View Lists
				</Link>
			</li>
			<li class="flex  gap-x-2">
				<span>
					<i class="fa-solid fa-users" />
				</span>
				<details>
					<summary> View your subs:</summary>

					<details>
						<summary>Teddit</summary>
						<ul>
							{#each $AuthStore.Subs?.teddit ?? ["nothing!"] as sub}
								<li>
									<Link href={`/teddit/r/${sub}`}>
										{sub}
									</Link>
								</li>
							{/each}
						</ul>
					</details>
					<details>
						<summary>Nitter</summary>
						<ul>
							{#each $AuthStore.Subs?.nitter ?? ["nothing!"] as nittos}
								<li>
									<Link href={`/nitter/${nittos}`}>
										{nittos}
									</Link>
								</li>
							{/each}
						</ul>
					</details>
				</details>
			</li>
		</ul>
		<div class="flex flex-col items-center mt-2 gap-y-2">
			<Link href="/auth">
				<button
					class="btn btn-warning btn-wide gap-x-2 mb-4"
					on:click={$AuthStore?.user.action?.logout}><i class="fa-solid fa-key" /> Logout</button
				>
			</Link>

			<div
				class="tooltip tooltip-bottom z-10 mb-4"
				data-tip="Regenerate all your recovery codes by overriding the current ones"
			>
				<button class="btn btn-accent btn-wide gap-x-2" on:click={RegenerateCodes}
					><i class="fa-solid fa-key" /> Generate recovery codes</button
				>
			</div>

			<Link href="/auth/reset-password">
				<button class="btn btn-info btn-wide gap-x-2"
					><i class="fa-solid fa-lock" /> Change Password</button
				>
			</Link>

			<div
				class="tooltip tooltip-bottom tooltip-warning z-10 mb-4"
				data-tip="This action is irreversible and will delete all the datas associated with this account"
			>
				<button class="btn btn-error btn-wide gap-x-2" on:click={HandleDeleteState}>
					{#if deleteState === DeleteStates.NoAction}
						<i class="fa-solid fa-user-minus" /> Delete Account
					{:else if deleteState === DeleteStates.Confirm}
						<i class="fa-solid fa-check-double" /> Are you sure?
					{:else}
						<i class="fa-solid fa-ban" /> Cancel
					{/if}
				</button>
			</div>
		</div>
	</aside>
</main>
<RecoveryModal
	show={showRecoveryCodesModal}
	recoveryCodes={newRecoveryCodes}
	setShow={setShowRecoveryCodesModal}
/>
