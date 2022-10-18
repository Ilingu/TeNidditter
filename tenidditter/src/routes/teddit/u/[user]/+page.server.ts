import api from "$lib/api";
import type { UserShape } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageServerLoad = async ({ params }): Promise<UserShape> => {
	const username = params?.user;
	if (IsEmptyString(username)) throw error(404, "Not Found -- Invalid username");

	try {
		const { success, data: UserInfos } = await api.get<UserShape>({
			uri: "/teddit/u",
			param: username
		});
		if (!success || typeof UserInfos !== "object" || !Object.hasOwn(UserInfos, "username"))
			throw error(404, "User Not found");

		return UserInfos;
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
