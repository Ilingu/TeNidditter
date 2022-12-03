import type { RequestHandler } from "@sveltejs/kit";

export const DELETE: RequestHandler = async () => {
	return new Response(null, {
		status: 205,
		headers: { "Clear-Site-Data": `"cache", "cookies", "storage", "executionContexts"` }
	});
};
