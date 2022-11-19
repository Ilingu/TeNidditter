<script lang="ts">
	import AuthStore from "$lib/stores/auth";
	import { FormatUsername } from "$lib/utils";
	import { GetZxcvbn, ScoreToColor, ScoreToText } from "$lib/zxcvbn";
	import { Authenticate, LogOut } from "$lib/services/auth";
	import type { PswReport } from "$lib/types/zxcvbn";
	import Link from "$lib/components/design/Link.svelte";
	import RecoveryModal from "./RecoveryModal.svelte";

	/* VAR */
	let Username = "",
		Password = "";

	let AuthMethod: "register" | "login" = "login";

	let PswStrenghtReport: PswReport | null,
		UsernameError = false;

	let loading = false;
	let showRecoveryCodesModal = false;
	const setShowRecoveryCodesModal = (newVal: boolean) => {
		showRecoveryCodesModal = newVal;
	};

	AuthStore.subscribe((value) => value.loggedIn && (Username = value.user?.username || ""));

	/* auth */
	const SignAuth = async () => {
		const recoveryCodeCb = (codes: string[]) => {
			if (codes?.length > 0) {
				recoveryCodes = codes;
				showRecoveryCodesModal = true;
			}
		};

		loading = true;
		await Authenticate(Username, Password, AuthMethod, recoveryCodeCb);
		loading = false;
		Reset();
	};

	/* inputs */
	const Reset = () => {
		Username = "";
		Password = "";
		UsernameError = false;
	};

	let debounce: NodeJS.Timeout;
	const PasswordChange = () => {
		clearTimeout(debounce);
		debounce = setTimeout(async () => {
			if (Password.trim().length <= 0) return (PswStrenghtReport = null);
			const zxcvbn = await GetZxcvbn();
			const { score, crack_times_display } = zxcvbn(Password);
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

	/* recoveryCodes */
	let recoveryCodes: string[] = [];
</script>

<div class="lg:col-span-2 col-span-6 bg-base-200 h-full flex flex-col justify-center items-center">
	<h1 class="font-fancy text-4xl mb-10">Welcome 👋</h1>
	<form on:submit|preventDefault={SignAuth} class="w-fit flex flex-col items-center gap-y-10">
		<!-- Username Input -->
		<div class="form-control relative">
			<span
				class={`placehover absolute select-none top-1/4 left-1/4 ${
					$AuthStore.loggedIn ? "hidden" : ""
				}`}>Username</span
			>
			<label class="input-group">
				<span class="fas fa-at" />
				<input
					autocomplete="off"
					disabled={loading || $AuthStore.loggedIn}
					bind:value={Username}
					on:input={(ev) => (Username = FormatUsername(ev.currentTarget.value))}
					type="text"
					class={`input input-bordered bg-transparent ${UsernameError ? "input-error" : ""}`}
				/>
			</label>
		</div>

		<!-- Password Input -->
		<div class="form-control relative">
			<span class="placehover absolute select-none top-1/4 left-1/4">Password</span>
			<label class="input-group">
				<span class="fas fa-lock" />
				<input
					autocomplete="off"
					disabled={loading || $AuthStore.loggedIn}
					bind:value={Password}
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

			<p class="absolute -top-5 text-sm right-0 hover:underline">
				<Link href="/auth/reset-password">Forgot Password?</Link>
			</p>

			<!-- Password Strenght Report -->
			{#if PswStrenghtReport}
				<div class="text-sm text-white text-center mt-2 -mb-7">
					<p>
						Password is <span class="font-bold" style={`color: ${PswStrenghtReport.score.color};`}
							>{PswStrenghtReport.score.text}</span
						>
					</p>
					<p><span class="font-bold">{PswStrenghtReport.crackTime}</span> to crack</p>
				</div>
			{/if}
		</div>

		<!-- Submit (Login/signup) -->
		{#if $AuthStore.loggedIn}
			<div class="w-full">
				<button
					class="btn btn-error btn-wide gap-2 btn-md"
					on:click={() => LogOut(true, $AuthStore.JwtToken)}
					><i class="fa-solid fa-right-from-bracket icon" /> Sign out</button
				>
			</div>
		{:else}
			<div class="grid grid-cols-2 w-full">
				<button
					disabled={loading ||
						(PswStrenghtReport && PswStrenghtReport?.score.number < 3) ||
						$AuthStore.loggedIn}
					class={`btn bg-base-300 gap-2 btn-md ${loading ? "loading" : ""}`}
					type="submit"
					on:click={() => (AuthMethod = "login")}
					><i class="fas fa-right-to-bracket icon" /> Sign in</button
				>
				<button
					disabled={loading ||
						(PswStrenghtReport && PswStrenghtReport?.score.number < 3) ||
						$AuthStore.loggedIn}
					class={`btn bg-base-300 gap-2 btn-md ${loading ? "loading" : ""}`}
					type="submit"
					on:click={() => (AuthMethod = "register")}><i class="fas fa-user icon" /> Sign up</button
				>
			</div>
		{/if}
	</form>
</div>
<RecoveryModal show={showRecoveryCodesModal} {recoveryCodes} setShow={setShowRecoveryCodesModal} />
