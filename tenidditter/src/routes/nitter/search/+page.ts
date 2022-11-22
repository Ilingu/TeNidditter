import api from "$lib/shared/api";
import type { NeetComment, NittosPreview } from "$lib/shared/types/nitter";
import { IsEmptyString } from "$lib/shared/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageLoad = async ({
	url
}): Promise<{ searchResult: NeetComment[][] | NittosPreview[] }> => {
	const type = (url.searchParams.get("type") as "tweets" | "users") ?? "tweets";
	const query = decodeURIComponent(url.searchParams.get("q") ?? "");
	if (IsEmptyString(type) || IsEmptyString(query)) return { searchResult: [] };

	const { success, data } = await api.get("/nitter/search", {
		query: { q: query, type, limit: 3 }
	});
	if (!success || typeof data !== "object") throw error(404, "nothing returned from nitter");

	return { searchResult: data };
};
