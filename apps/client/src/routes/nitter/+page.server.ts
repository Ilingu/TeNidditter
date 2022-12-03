import api from "$lib/shared/api";
import type { NeetComment } from "$lib/shared/types/nitter";
import { IsEmptyString, MakeBearerToken } from "$lib/shared/utils";
import { redirect } from "@sveltejs/kit";

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
		throw redirect(307, `/nitter/search`);

	return { comments };
};
