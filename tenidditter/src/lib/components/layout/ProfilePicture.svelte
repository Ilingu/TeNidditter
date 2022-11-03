<script lang="ts">
	import { GenerateRandomPPUrl } from "$lib/services/ProfilesPictures";

	export let size: "small" | "medium" | "big" = "medium";
</script>

<div class={`${size === "big" ? "w-24" : size === "medium" ? "w-12" : "w-6"} mask mask-hexagon`}>
	<img
		src={GenerateRandomPPUrl(size === "big" ? 96 : size === "medium" ? 48 : 24)}
		alt="fake profile pictures"
	/>
</div>
