<script lang="ts">
	import { pushAlert } from "$lib/client/ClientUtils";
	import type { PswReport } from "$lib/client/types/zxcvbn";
	import { GetZxcvbn, ScoreToColor, ScoreToText } from "$lib/client/zxcvbn";
	import api from "$lib/shared/api";
	import { FormatUsername, IsEmptyString } from "$lib/shared/utils";

	/* VAR */
	let Username = "",
		NewPassword = "",
		RecoveryCode = "";

	let PswStrenghtReport: PswReport | null,
		UsernameError = false,
		RecoveryCodeError = false;

	let loading = false;

	let debounce: NodeJS.Timeout;
	const PasswordChange = () => {
		clearTimeout(debounce);
		debounce = setTimeout(async () => {
			if (NewPassword.trim().length <= 0) return (PswStrenghtReport = null);
			const zxcvbn = await GetZxcvbn();

			const { score, crack_times_display } = zxcvbn(NewPassword);
			PswStrenghtReport = {
				crackTime: crack_times_display.online_throttling_100_per_hour.toString(),
				score: {
					number: score,
					text: ScoreToText[score],
					color: ScoreToColor[score]
				}
			};
		}, 500);
	};

	const ChangePassword = async () => {
		// Username check
		const username = FormatUsername(Username.trim());
		const nameLen = username.length;
		if (nameLen < 3 || nameLen > 15)
			return (UsernameError = true) && pushAlert("Username invalid!", "warning");

		// Recovery code check
		if (IsEmptyString(RecoveryCode) || RecoveryCode.length !== 10)
			return (RecoveryCodeError = true) && pushAlert("Recovery Code invalid!", "warning");

		// Check Password
		const zxcvbn = await GetZxcvbn();
		const passwordStrenght = zxcvbn(NewPassword);
		if (passwordStrenght.score < 3) return pushAlert("Password is too weak! Like you", "warning");

		loading = true;
		const { success, error } = await api.put("/auth/reset-password", {
			body: { username, RecoveryCode, NewPassword }
		});
		if (!success) pushAlert(`Failed: ${error}`, "error", 15_000);
		else
			pushAlert(
				"Password successfully updated. The code you've used is now desactivated, if it was the last one of your lists please go regenerated a new list before you find yourself locked out of your account.",
				"info",
				25_000
			);

		Reset();
		loading = false;
	};

	const Reset = () => {
		Username = "";
		NewPassword = "";
		RecoveryCode = "";
		UsernameError = false;
		RecoveryCodeError = false;
	};
</script>

<section class="page-content w-screen">
	<main class="bg-base-200 h-full flex flex-col justify-center items-center">
		<h1 class="font-fancy text-4xl mb-10">Lost something <i class="fa-solid fa-question" /></h1>
		<form
			on:submit|preventDefault={ChangePassword}
			class="w-fit flex flex-col items-center gap-y-5"
		>
			<!-- Username Input -->
			<div class="form-control relative">
				<span class="placehover absolute select-none top-1/4 left-1/4">Username</span>
				<label class="input-group">
					<span class="fas fa-at" />
					<input
						autocomplete="off"
						disabled={loading}
						bind:value={Username}
						on:input={(ev) => (Username = FormatUsername(ev.currentTarget.value))}
						type="text"
						class={`input input-bordered bg-transparent ${UsernameError ? "input-error" : ""}`}
					/>
				</label>
			</div>

			<!-- Recovery Input -->
			<div
				class="tooltip tooltip-bottom tooltip-info z-10"
				data-tip="Single use code, if you've already used it, it won't work anymore."
			>
				<div class="form-control relative">
					<span class="placehover absolute select-none top-1/4 left-1/4">Recovery Code</span>
					<label class="input-group">
						<span class="fa-solid fa-key" />
						<input
							autocomplete="off"
							disabled={loading}
							bind:value={RecoveryCode}
							on:input={(ev) => (RecoveryCode = ev.currentTarget.value)}
							type="text"
							class={`input input-bordered bg-transparent ${
								RecoveryCodeError ? "input-error" : ""
							}`}
						/>
					</label>
				</div>
			</div>

			<!-- Password Input -->
			<div class="form-control relative">
				<span class="placehover absolute select-none top-1/4 left-1/4">NEW Password</span>
				<label class="input-group">
					<span class="fas fa-lock" />
					<input
						autocomplete="off"
						bind:value={NewPassword}
						on:input={PasswordChange}
						type="password"
						class={`input input-bordered bg-transparent ${
							PswStrenghtReport
								? PswStrenghtReport.score.number < 3
									? "input-error"
									: "input-success"
								: ""
						}`}
					/>
				</label>

				<!-- Password Strenght Report -->
				{#if PswStrenghtReport}
					<div class="text-sm text-white text-center mt-2">
						<p>
							Password is <span class="font-bold" style={`color: ${PswStrenghtReport.score.color};`}
								>{PswStrenghtReport.score.text}</span
							>
						</p>
						<p><span class="font-bold">{PswStrenghtReport.crackTime}</span> to crack</p>
					</div>
				{/if}
			</div>

			<div class="w-full flex justify-center">
				<button class="btn btn-wide btn-info gap-x-2 btn-md" type="submit"
					><i class="fa-solid fa-right-from-bracket icon" /> Change Password</button
				>
			</div>
		</form>
	</main>
</section>
