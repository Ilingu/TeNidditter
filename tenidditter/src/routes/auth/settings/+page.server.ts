import { IsEmptyString } from "$lib/utils";
import { redirect } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageServerLoad = ({ cookies }) => {
	const eToken = cookies.get("JwtToken");
	if (!eToken || IsEmptyString(eToken)) throw redirect(307, `/auth`);
};
