import { redirect } from "@sveltejs/kit";

// Redirects all routes beginning by "/search/..." by "/nitter/search/..."
export const load: import("./$types").PageServerLoad = ({ url }) => {
	throw redirect(301, `/nitter${url.pathname + url.search}`);
};
