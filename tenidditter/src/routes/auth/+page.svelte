<script lang="ts">
	import { callApi } from "$lib/utils/server";
	import { FormatUsername, pushAlert } from "$lib/utils/utils";

	import { GetZxcvbn, ScoreToColor, ScoreToText } from "$lib/utils/zxcvbn";
	import { onMount } from "svelte";

	/* TYPES */
	interface PswReport {
		score: { text: string; number: number; color: string };
		crackTime: string;
	}

	/* VAR */
	let Username = "",
		Password = "";

	let PswStrenghtReport: PswReport | null,
		UsernameError = false;

	let loading = false;

	const HeadlineText = [
		"Privacy comes first.",
		"Encrypted, private, secure.",
		"❤️ = 🔒 + 🕊",
		"Social Medias with freedom."
	];
	const DescText = [
		`Are you looking for an more privacy friendly alternative to <span class="text-teddit">reddit</span> or <span class="text-nitter">twitter</span>?`,
		`Tenidditter app is what social medias should have become <span class="font-bold">long ago</span>.`,
		`Tenidditter app act as a proxy between twitter/reddit's datas and you, you never directly speak to them, therefore you can keep <span class="font-bold">entertaining yourself without seeing your freedom being stolen</span>.`
	];

	/* App Start */
	onMount(() => {
		document.querySelectorAll("input").forEach((inp) => {
			inp.addEventListener("focusin", () => PlacehoverAnimate(inp, "add"));
			inp.addEventListener("focusout", () => {
				if (inp.value.trim().length >= 1) return;
				PlacehoverAnimate(inp, "rem");
			});
		});
		DrawLines();
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

		// TODO: db endpoint to see if username overlap
		const { success, data } = await callApi<boolean>({
			uri: `/auth/available?username=${encodeURI(username)}`,
			method: "GET"
		});
		if (!success || !data) return pushAlert("This Username is already taken.", "warning", 6000);

		console.log(data);
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

	const DrawLines = () => {
		const canvas = document.getElementById("LinesCanvas") as HTMLCanvasElement;
		if (!canvas?.getContext) return;
		const ctx = canvas.getContext("2d");
		if (!ctx) return;

		canvas.height *= 4;
		canvas.width *= 5;

		const smoothness = Math.round(Math.random() * 24) + 1;
		const linesNumbers = Math.round(Math.random() * 50) + 50;

		for (let lineID = 0; lineID < linesNumbers; lineID++) {
			const h = Math.round(Math.random() * canvas.height);
			// const a = Math.round(Math.random()) + 0.5;
			// const f = Math.round(Math.random() * 2) + 1;

			// Filled sinwave
			ctx.beginPath();
			ctx.moveTo(0, h);
			let lastCoord: [number, number] = [0, h];
			for (let i = 0; i < canvas.width; i += smoothness) {
				let y = Math.sin(((i % 360) * 2 * Math.PI) / 360); // Calculate y value from x
				ctx.moveTo(i, lastCoord[1]); // Where to start drawing
				ctx.lineTo(lastCoord[0], h + 25 * y); // Where to draw to
				lastCoord = [i, h + 25 * y];
			}

			ctx.strokeStyle = Math.round(Math.random()) === 0 ? "#FF4500" : "#1DA1F2";
			ctx.stroke();
		}
	};
</script>

<section class="w-screen grid grid-cols-6">
	<div class="col-span-2 bg-base-200 h-full flex flex-col justify-center items-center">
		<h1 class="font-fancy text-4xl mb-10">Welcome 👋</h1>
		<form
			on:submit|preventDefault={() => {
				loading = true;
				Authenticate();
				loading = false;
			}}
			class="w-fit flex flex-col items-center gap-y-10"
		>
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
						class={`input input-bordered bg-transparent ${
							PswStrenghtReport
								? PswStrenghtReport.score.number < 3
									? "input-error"
									: "input-success"
								: ""
						}`}
					/>
				</label>
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

			<button
				disabled={loading || (PswStrenghtReport && PswStrenghtReport?.score.number < 3)}
				class={`btn btn-wide bg-base-300 flex gap-2 btn-sm md:btn-md ${loading ? "loading" : ""}`}
				type="submit"><span class="fas fa-user" /> Login/Sign Up</button
			>
		</form>
	</div>
	<aside class="col-span-4 relative">
		<canvas id="LinesCanvas" class="w-full h-full" />
		<div class="quote mockup-window border bg-base-300 absolute bottom-0 w-3/4">
			<div class="font-mono px-4 pt-2 h-32 bg-base-200">
				<h2 class="font-nerd text-2xl font-semibold text-primary">
					{HeadlineText[Math.round(Math.random() * (HeadlineText.length - 1))]}
				</h2>
				<p class="mt-3">{@html DescText[Math.round(Math.random() * (DescText.length - 1))]}</p>
			</div>
		</div>
	</aside>
</section>

<style scoped>
	section {
		overflow: hidden;
	}

	input[disabled] {
		opacity: 0.5;
	}

	.quote {
		left: calc(50% - 75% / 2);
		animation: PopIn 1s cubic-bezier(0.6, 0, 0.96, 0.68) 0.1s 1 forwards;
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

	@keyframes PopIn {
		from {
			transform: translateY(400px);
		}
		to {
			transform: translateY(-25%);
		}
	}
</style>
