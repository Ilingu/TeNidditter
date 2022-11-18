import { redirect } from "@sveltejs/kit";

// Redirects all routes beginning by "/u/..." by "/teddit/u/..."
export const load: import("./$types").PageServerLoad = ({ url }) => {
	throw redirect(301, `/teddit${url.pathname}`);
};
