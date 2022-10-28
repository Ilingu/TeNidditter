import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageLoad = async ({ params }) => {
	const subredditName = params?.subteddit;
	const postId = params?.postid;
	if (IsEmptyString(subredditName) || IsEmptyString(postId) || postId.length < 6)
		throw error(400, "Invalid args");

	try {
		// api.get("")
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
