import api from "$lib/api";
import type { NeetInfo } from "$lib/types/interfaces";
import { IsEmptyString, TrimSpecialChars } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

type ReturnType = { NeetComments: NeetInfo };
export const load: import("./$types").PageLoad = async ({ params }): Promise<ReturnType> => {
	const nittosName = params?.nittos,
		neetId = TrimSpecialChars(params?.neet); //trim special char
	if (IsEmptyString(nittosName) || IsEmptyString(neetId) || neetId.length < 19)
		throw error(400, "Bad args");

	const { success, data: NeetComments } = await api.get("/nitter/nittos/%s/neets/%s", {
		params: [nittosName, neetId],
		query: { limit: 3 }
	});
	if (!success || typeof NeetComments !== "object" || !Object.hasOwn(NeetComments, "reply"))
		throw error(404, "this neet does not exist");

	return { NeetComments };
};
