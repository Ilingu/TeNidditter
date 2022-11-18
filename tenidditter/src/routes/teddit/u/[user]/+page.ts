import api from "$lib/api";
import type { TedditUserShape } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

// fetch the teddit user datas
export const load: import("./$types").PageLoad = async ({
	params,
	fetch
}): Promise<TedditUserShape> => {
	const username = params?.user;
	if (IsEmptyString(username)) throw error(400, "Invalid username");

	const { success, data: UserInfos } = await api.get("/teddit/u/%s", { params: [username] }, fetch);
	if (!success || typeof UserInfos !== "object" || !Object.hasOwn(UserInfos, "username"))
		throw error(404, "User Not found");

	return UserInfos;
};
