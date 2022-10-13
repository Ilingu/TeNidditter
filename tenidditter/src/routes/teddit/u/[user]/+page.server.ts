import type { UserShape } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageServerLoad = async ({ params }): Promise<UserShape> => {
	const username = params?.user;
	if (IsEmptyString(username)) throw error(404, "Not Found -- Invalid username");

	try {
		const resp = await fetch(`https://teddit.net/u/${username}?api&raw_json=1`);
		if (!resp.ok) throw error(404, "User Not found");

		const user: UserShape = await resp.json();
		if (typeof user !== "object") throw error(404, "User Not found");

		return user;
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
