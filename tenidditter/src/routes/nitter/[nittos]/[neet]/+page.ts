import type { NeetInfo } from "$lib/server/types/nitter";
import api from "$lib/shared/api";
import { IsEmptyString, TrimNonDigitsChars } from "$lib/shared/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

type ReturnType = { NeetComments: NeetInfo };
export const load: import("./$types").PageLoad = async ({ params }): Promise<ReturnType> => {
	const nittosName = params?.nittos,
		neetId = TrimNonDigitsChars(params?.neet); //trim special char
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
