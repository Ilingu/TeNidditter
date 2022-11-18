<script lang="ts">
	import AuthStore from "$lib/stores/auth";
	import api from "$lib/api";
	import {
		copyToClipboard,
		FormatUsername,
		IsEmptyString,
		IsValidJSON,
		pushAlert
	} from "$lib/utils";
	import { GetZxcvbn, ScoreToColor, ScoreToText } from "$lib/zxcvbn";
	import { LogOut, SignIn } from "$lib/services/auth";
	import type { NitterLists } from "$lib/types/interfaces";
	import type { PswReport } from "$lib/types/zxcvbn";
	import Link from "$lib/components/design/Link.svelte";

	/* VAR */
	let Username = "",
		Password = "";

	let AuthMethod: "signup" | "login" = "login";

	let PswStrenghtReport: PswReport | null,
		UsernameError = false;

	let loading = false;

	let showRecoveryCodesModal: HTMLLabelElement;
	let hideRecoveryCodesModal: HTMLLabelElement;
	let recoverySaved = false;

	AuthStore.subscribe((value) => value.loggedIn && (Username = value.user?.username || ""));

	/* auth */
	const Authenticate = async () => {
		// Username check
		const username = FormatUsername(Username.trim());
		const nameLen = username.length;
		if (nameLen < 3 || nameLen > 15)
			return (UsernameError = true) && pushAlert("Username invalid!", "warning");

		// Check Password
		const zxcvbn = await GetZxcvbn();
		const passwordStrenght = zxcvbn(Password);
		if (passwordStrenght.score < 3) return pushAlert("Password is too weak! Like you", "warning");

		if (AuthMethod === "signup") {
			const { success, data: IsAvailable } = await api.get("/auth/available", {
				query: { username }
			});

			if (!success || !IsAvailable)
				return pushAlert("This Username is already taken.", "warning", 6000);
		}

		const {
			success: AuthSuccess,
			data: JwtToken,
			headers
		} = await api.post("/auth/", {
			credentials: true,
			body: {
				username,
				password: Password
			}
		});

		if (!AuthSuccess) return pushAlert("Failed to auth", "error");

		if (AuthMethod === "signup") return AfterRegister(headers);
		if (AuthMethod === "login" && !IsEmptyString(JwtToken))
			return AfterLogin(JwtToken as string, headers);

		pushAlert("Invalid login", "error");
	};

	const AfterRegister = async (headers?: Headers) => {
		const codesHeader = headers?.get("RecoveryCodes") ?? "";

		if (codesHeader && !IsEmptyString(codesHeader) && IsValidJSON(codesHeader))
			recoveryCodes = JSON.parse(codesHeader);
		if (recoveryCodes.length > 0) showRecoveryCodesModal?.click();

		pushAlert("Successfully registered, you can now login", "success", 6000);

		loading = true;
		AuthMethod = "login";

		await Authenticate();
		loading = false;

		Reset();
	};

	const AfterLogin = async (JwtToken: string, headers?: Headers) => {
		const tedditSubsHeader = headers?.get("TedditSubs") ?? "",
			nitterSubsHeader = headers?.get("NitterSubs") ?? "",
			nitterListsHeader = headers?.get("NitterLists") ?? ""; // retrieve user subs

		let tedditSubs: string[] = [],
			nitterSubs: string[] = [],
			nitterLists: NitterLists[] = [];

		if (tedditSubsHeader && !IsEmptyString(tedditSubsHeader) && IsValidJSON(tedditSubsHeader))
			tedditSubs = JSON.parse(tedditSubsHeader);
		if (nitterSubsHeader && !IsEmptyString(nitterSubsHeader) && IsValidJSON(nitterSubsHeader))
			nitterSubs = JSON.parse(nitterSubsHeader);
		if (nitterListsHeader && !IsEmptyString(nitterListsHeader) && IsValidJSON(nitterListsHeader))
			nitterLists = JSON.parse(nitterListsHeader);

		const subs = { teddit: tedditSubs, nitter: nitterSubs };
		await SignIn(JwtToken, subs, nitterLists); // validate user jwt

		Reset();
		if (!$AuthStore.loggedIn) return pushAlert("Invalid login", "error");

		pushAlert("Successfully logged in!", "success");
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

<div class="lg:col-span-2 col-span-6 bg-base-200 h-full flex flex-col justify-center items-center">
	<h1 class="font-fancy text-4xl mb-10">Welcome 👋</h1>
	<form
		on:submit|preventDefault={async () => {
			loading = true;
			await Authenticate();
			loading = false;
		}}
		class="w-fit flex flex-col items-center gap-y-10"
	>
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
					on:click={() => (AuthMethod = "signup")}><i class="fas fa-user icon" /> Sign up</button
				>
			</div>
		{/if}
	</form>

	<!-- The button to open modal -->
	<label for="recoveryCodeModal" class="hidden" bind:this={showRecoveryCodesModal} />
</div>

<!-- recovery codes Modal -->
<input type="checkbox" id="recoveryCodeModal" class="modal-toggle" />
<div class="modal ">
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
					hideRecoveryCodesModal?.click();
				}}>💯 All good!</button
			>
			<label for="recoveryCodeModal" class="hidden" bind:this={hideRecoveryCodesModal} />
		</div>
	</div>
</div>
