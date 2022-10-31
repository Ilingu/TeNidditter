<script lang="ts">
	import AuthStore from "$lib/stores/auth";
	import api from "$lib/api";
	import { FormatUsername, IsEmptyString, IsValidJSON, pushAlert } from "$lib/utils";
	import { GetZxcvbn, ScoreToColor, ScoreToText } from "$lib/zxcvbn";
	import { onMount } from "svelte";
	import { SignIn } from "$lib/services/auth";

	/* TYPES */
	interface PswReport {
		score: { text: string; number: number; color: string };
		crackTime: string;
	}

	/* VAR */
	let Username = "",
		Password = "";

	let AuthMethod: "signup" | "login" = "login";

	let PswStrenghtReport: PswReport | null,
		UsernameError = false;

	let loading = false;

	AuthStore.subscribe((value) => value.loggedIn && (Username = value.user?.username || ""));

	/* App Start */
	onMount(() => {
		document.querySelectorAll("input").forEach((inp) => {
			inp.addEventListener("focusin", () => PlacehoverAnimate(inp, "add"));
			inp.addEventListener("focusout", () => {
				if (inp.value.trim().length >= 1) return;
				PlacehoverAnimate(inp, "rem");
			});
		});
	});

	/* INPUTS */
	const Authenticate = async () => {
		// Username check
		const username = FormatUsername(Username.trim());
		const nameLen = username.length;
		if (nameLen < 3 || nameLen > 15)
			return (UsernameError = true) && pushAlert("Username invalid!", "warning");

		// Checks Password
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

		if (!AuthSuccess) pushAlert("Failed to auth", "error");
		else if (AuthMethod === "signup")
			pushAlert("Successfully registered, you can now login", "success", 6000);
		else if (!IsEmptyString(JwtToken)) AfterLogin(JwtToken as string, headers);
		else pushAlert("Invalid login", "error");

		Reset();
	};

	const AfterLogin = async (JwtToken: string, headers?: Headers) => {
		const tedditSubs = headers?.get("TedditSubs"); // retrieve user subs
		if (!tedditSubs || IsEmptyString(tedditSubs) || !IsValidJSON(tedditSubs))
			return pushAlert("Invalid login", "error");

		await SignIn(JwtToken, { teddit: JSON.parse(tedditSubs), nitter: [] }); // validate user jwt
		if (!$AuthStore.loggedIn) return;

		pushAlert("Successfully logged in!", "success");
	};

	const Reset = () => {
		Username = "";
		Password = "";
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

	/* ANIMATION */
	const PlacehoverAnimate = (inp: HTMLInputElement, mode: "add" | "rem") => {
		const label = inp.parentElement?.parentElement?.firstChild as HTMLElement;
		const hasAnimate = label?.classList.contains("animate");

		if (mode === "add" && hasAnimate) return;
		if (mode === "rem" && !hasAnimate) return;

		label?.classList.toggle("animate");
	};
</script>

<div class="col-span-2 bg-base-200 h-full flex flex-col justify-center items-center">
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
		<div class="grid grid-cols-2 w-full">
			<button
				disabled={loading ||
					(PswStrenghtReport && PswStrenghtReport?.score.number < 3) ||
					$AuthStore.loggedIn}
				class={`btn  bg-base-300 flex gap-2 btn-sm md:btn-md ${loading ? "loading" : ""}`}
				type="submit"
				on:click={() => (AuthMethod = "login")}
				><i class="fas fa-right-to-bracket icon" /> Sign in</button
			>
			<button
				disabled={loading ||
					(PswStrenghtReport && PswStrenghtReport?.score.number < 3) ||
					$AuthStore.loggedIn}
				class={`btn  bg-base-300 flex gap-2 btn-sm md:btn-md ${loading ? "loading" : ""}`}
				type="submit"
				on:click={() => (AuthMethod = "signup")}><i class="fas fa-user icon" /> Sign up</button
			>
		</div>
	</form>
</div>

<style scoped>
	input[disabled] {
		opacity: 0.5;
	}

	.placehover {
		transition: all 0.5s;
	}

	:global(.placehover.animate) {
		translate: 0 -40px;
	}
</style>
