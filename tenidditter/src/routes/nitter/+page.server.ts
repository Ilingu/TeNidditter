import api from "$lib/api";
import type { NeetComment } from "$lib/types/interfaces";
import { IsEmptyString, MakeBearerToken } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageServerLoad = async ({
	cookies
}): Promise<{ comments: NeetComment[][] }> => {
	const eToken = cookies.get("JwtToken");
	if (!eToken || IsEmptyString(eToken)) throw redirect(307, `/nitter/search`);

	const { success, data: comments } = await api.get("/tedinitter/nitter/feed", {
		headers: MakeBearerToken(eToken)
	});
	if (!success || typeof comments !== "object" || comments.length <= 0)
		throw error(404, "feed not found");

	return { comments };
};
