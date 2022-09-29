<script lang="ts">
	import { pushAlert } from "$lib/utils/utils";

	import { GetZxcvbn, ScoreToText } from "$lib/utils/zxcvbn";
	import { onMount } from "svelte";

	interface PswReport {
		score: string;
		feedback: string[];
		crackTime: string;
	}

	let Username = "",
		Password = "";

	let PswStrenghtReport: PswReport;

	let AuthBtn: HTMLButtonElement;
	let loading = false;

	onMount(() => {
		document.querySelectorAll("input").forEach((inp) => {
			inp.addEventListener("focusin", () => PlacehoverAnimate(inp, "add"));
			inp.addEventListener("focusout", () => {
				if (inp.value.trim().length >= 1) return;
				PlacehoverAnimate(inp, "rem");
			});
		});
	});
	const PlacehoverAnimate = (inp: HTMLInputElement, mode: "add" | "rem") => {
		const label = inp.parentElement?.parentElement?.firstChild as HTMLElement;
		const hasAnimate = label?.classList.contains("animate");

		if (mode === "add" && hasAnimate) return;
		if (mode === "rem" && !hasAnimate) return;

		label?.classList.toggle("animate");
	};

	const Authenticate = async () => {
		{
			loading = true;
			AuthBtn?.classList.add("loading");
		}

		// Checks Password
		const zxcvbn = await GetZxcvbn();
		const passwordStrenght = zxcvbn(Password);
		if (passwordStrenght.score < 3) return pushAlert("Password is too weak! Like you", "warning");

		const nameLen = Username.trim().length;
		if (nameLen < 3 || nameLen > 15) return pushAlert("Username invalid!", "warning");

		// NEXT: db endpoint to see if username overlap
		// + password bar

		{
			loading = false;
			AuthBtn?.classList.remove("loading");
		}
	};

	let debounce: NodeJS.Timeout;
	const PasswordChange = () => {
		clearTimeout(debounce);
		debounce = setTimeout(async () => {
			const zxcvbn = await GetZxcvbn();
			const { score, feedback, crack_times_display } = zxcvbn(Password);
			PswStrenghtReport = {
				crackTime: crack_times_display.online_throttling_100_per_hour.toString(),
				score: ScoreToText[score],
				feedback: feedback?.suggestions
			};
		}, 500);
	};
</script>

<section class="w-screen grid grid-cols-6">
	<div class="col-span-2 bg-base-200 h-full flex flex-col justify-center items-center">
		<h1 class="font-fancy text-4xl mb-10">Welcome 👋</h1>
		<form on:submit|preventDefault={Authenticate} class="w-fit flex flex-col items-center gap-y-10">
			<div class="form-control relative">
				<span class="placehover absolute select-none top-1/4 left-1/4">Username</span>
				<label class="input-group">
					<span class="fas fa-at" />
					<input
						autocomplete="off"
						disabled={loading}
						bind:value={Username}
						type="text"
						class="input input-bordered bg-transparent"
					/>
				</label>
			</div>
			<div class="form-control relative">
				<span class="placehover absolute select-none top-1/4 left-1/4">Password</span>
				<label class="input-group">
					<span class="fas fa-lock" />
					<input
						autocomplete="off"
						disabled={loading}
						bind:value={Password}
						on:input={PasswordChange}
						type="password"
						class="input input-bordered bg-transparent"
					/>
				</label>
			</div>
			{#if PswStrenghtReport}
				<div class="text-sm text-white leading-4">
					<p>
						Password is <span class="font-bold">{PswStrenghtReport.score}</span>
						-- <span class="font-bold">{PswStrenghtReport.crackTime}</span> to crack
					</p>

					{#if PswStrenghtReport?.feedback}
						<ul>
							{#each PswStrenghtReport?.feedback as feedback}
								<li>{feedback}</li>
							{/each}
						</ul>
					{/if}
				</div>
			{/if}

			<button
				disabled={loading}
				bind:this={AuthBtn}
				class="btn btn-wide flex gap-2  btn-sm md:btn-md"
				type="submit"><span class="fas fa-user" /> Login/Sign Up</button
			>
		</form>
	</div>
	<aside class="col-span-4" />
</section>

<style scoped>
	input[disabled] {
		opacity: 0.5;
	}

	section {
		height: calc(100vh - 64px);
	}

	.placehover {
		transition: all 0.5s;
	}

	:global(.placehover.animate) {
		transform: translateY(-40px);
	}
</style>
