import api from "$lib/api";
import type { TedditPostInfo } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageLoad = async ({
	params,
	fetch
}): Promise<TedditPostInfo> => {
	const subredditName = params?.subteddit;
	const postId = params?.postid;
	if (IsEmptyString(subredditName) || IsEmptyString(postId) || postId.length < 6)
		throw error(400, "Invalid args");

	try {
		const { success, data: PostInfo } = await api.get(
			"/teddit/r/%s/post/%s",
			{
				params: [subredditName, postId]
			},
			fetch
		);
		if (!success || typeof PostInfo !== "object" || !Object.hasOwn(PostInfo, "comments"))
			throw error(404, "Post Not Found");

		const hasDuplicates = PostInfo.comments
			.map((p) => p.map(({ id }) => id))
			.map((a) => new Set(a).size === a.length);
		if (!hasDuplicates.every((t) => t === true)) throw error(500, "Corrupted Datas");

		return PostInfo;
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
