<script lang="ts">
	import { pushAlert } from "$lib/utils/utils";

	import { GetZxcvbn } from "$lib/utils/zxcvbn";
	import { onMount } from "svelte";

	let Username = "",
		Password = "";

	let AuthBtn: HTMLButtonElement;
	let loading = false;

	onMount(() => {
		document.querySelectorAll("input").forEach((inp) => {
			inp.addEventListener("focusin", () => PlacehoverAnmiate(inp));
			inp.addEventListener("focusout", () => PlacehoverAnmiate(inp));
		});
	});

	const PlacehoverAnmiate = (inp: HTMLInputElement) => {
		const label = inp.parentElement?.parentElement?.firstChild as HTMLElement;
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
		if (passwordStrenght.score < 3) return pushAlert("Password is too weak!", "warning");

		const nameLen = Username.trim().length;
		if (nameLen < 3 || nameLen > 15) return pushAlert("Username invalid!", "warning");

		// NEXT: db endpoint to see if username overlap
		// +

		{
			loading = false;
			AuthBtn?.classList.remove("loading");
		}
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
						type="password"
						class="input input-bordered bg-transparent"
					/>
				</label>
			</div>

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
