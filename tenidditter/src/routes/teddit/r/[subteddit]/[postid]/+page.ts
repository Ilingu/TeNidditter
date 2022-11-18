import api from "$lib/api";
import type { TedditPostInfo } from "$lib/types/interfaces";
import type { Tuple } from "$lib/types/types";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

// fetch post's comment by types
const acceptedSort: Tuple<string, 6> = ["best", "top", "new", "controversial", "old", "qa"];
export const load: import("./$types").PageLoad = async ({
	params,
	fetch,
	url
}): Promise<TedditPostInfo> => {
	const subredditName = params?.subteddit;
	const postId = params?.postid;
	if (IsEmptyString(subredditName) || IsEmptyString(postId) || postId.length < 6)
		throw error(400, "Invalid args");

	let sort = (url.searchParams.get("sort") ?? "").toLowerCase();
	if (!acceptedSort.includes(sort)) sort = acceptedSort[0];

	try {
		const { success, data: PostInfo } = await api.get(
			"/teddit/r/%s/post/%s",
			{
				params: [subredditName, postId],
				query: { sort }
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
