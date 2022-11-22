<script lang="ts">
	import Feeds from "$lib/client/components/pages/nitter/Feeds.svelte";
	import { page } from "$app/stores";
	import api from "$lib/api";

	export let data: import("./$types").PageData;

	let lastLimit = 3;
	const queryMore = async () => {
		const { success, data: NeetComments } = await api.get("/nitter/nittos/%s/neets/%s", {
			params: [$page.params.nittos, $page.params.neet],
			query: { limit: lastLimit + 1 }
		});
		if (!success || typeof NeetComments !== "object" || !Object.hasOwn(NeetComments, "reply"))
			return;
		lastLimit++;
		data = { NeetComments };
	};
</script>

<main class="grid place-items-center mt-5">
	<div class="max-w-[1250px] flex justify-center gap-x-10 z-0">
		<Feeds neets={[data.NeetComments.main]} />
	</div>

	<div class="max-w-[1250px] mt-10 flex justify-center gap-x-10 z-0">
		<Feeds neets={data.NeetComments.reply} queryMoreCb={queryMore} />
	</div>
</main>
