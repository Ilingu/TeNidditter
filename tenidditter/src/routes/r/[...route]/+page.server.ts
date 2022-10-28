import { redirect } from "@sveltejs/kit";

export const load: import("./$types").PageServerLoad = ({ url }) => {
	throw redirect(301, `/teddit${url.pathname}`);
};
