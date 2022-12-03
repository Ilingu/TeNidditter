import { redirect } from "@sveltejs/kit";

export const prerender = false;

// Redirects all routes beginning by "/u/..." by "/teddit/u/..."
export const load: import("./$types").PageServerLoad = ({ url }) => {
	throw redirect(301, `/teddit${url.pathname}`);
};
