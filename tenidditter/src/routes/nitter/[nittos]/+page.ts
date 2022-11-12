import api from "$lib/api";
import type { NeetComment, Nittos } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

interface ReturnType {
	userInfo: Nittos;
	userNeets: NeetComment[][] | undefined;
}
export const load: import("./$types").PageLoad = async ({ params }): Promise<ReturnType> => {
	const nittosName = params?.nittos;
	if (IsEmptyString(nittosName)) throw error(400, "Bad nittos name");

	const { success: nittosQuerySucceed, data: Nittos } = await api.get("/nitter/nittos/%s/about", {
		params: [nittosName]
	});
	if (!nittosQuerySucceed || typeof Nittos !== "object")
		throw error(404, "this user does not exist");
	if (!Object.hasOwn(Nittos, "username") || IsEmptyString(Nittos.username))
		throw error(404, "this user does not exist");

	// eslint-disable-next-line prefer-const
	let { success: NittosNeetsQuerySucceed, data: Neets } = await api.get("/nitter/nittos/%s/neets", {
		params: [nittosName],
		query: { limit: 3 }
	});
	if (!NittosNeetsQuerySucceed || typeof Neets !== "object") Neets = undefined;

	return { userInfo: Nittos, userNeets: Neets };
};
