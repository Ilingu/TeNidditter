import api from "$lib/shared/api";
import type { NeetComment } from "$lib/shared/types/nitter";
import { IsEmptyString, MakeBearerToken } from "$lib/shared/utils";
import { error, redirect } from "@sveltejs/kit";

type ReturnType = { savedNeets: NeetComment[] | null };
export const load: import("./$types").PageServerLoad = async ({
	cookies,
	params
}): Promise<ReturnType> => {
	const eToken = cookies.get("JwtToken");
	if (!eToken || IsEmptyString(eToken)) throw redirect(307, `/nitter/search`);

	const listId = params?.listId;
	if (IsEmptyString(listId)) throw error(400, "Bad listId");

	const {
		success,
		data: savedNeets,
		status
	} = await api.get("/tedinitter/nitter/list/%s", {
		headers: MakeBearerToken(eToken),
		params: [listId]
	});

	if (status === 204) return { savedNeets: null };
	else if (!success || typeof savedNeets !== "object" || savedNeets.length <= 0)
		throw error(404, "list not found");

	return { savedNeets };
};
