import { IsEmptyString } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";

export const load: import("./$types").PageServerLoad = async ({ cookies }) => {
	const eToken = cookies.get("JwtToken");
	if (IsEmptyString(eToken)) throw redirect(307, `/nitter/search`);

	return "";
};
