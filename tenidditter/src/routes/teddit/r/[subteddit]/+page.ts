import api from "$lib/api";
import type { TedditRawPost } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

interface SubRedditDatas {
	Info?: {
		subs: string;
		description: string;
		rules: string;
	};
	Feed?: TedditRawPost[];
}

export const load: import("./$types").PageLoad = async ({
	params,
	fetch
}): Promise<SubRedditDatas> => {
	const subredditName = params?.subteddit;
	if (IsEmptyString(subredditName)) throw error(400, "Invalid subreddit");

	try {
		const { success, data: SubPosts } = await api.get(
			"/teddit/r/%s/posts",
			{
				params: [subredditName]
			},
			fetch
		);
		if (!success || typeof SubPosts !== "object" || !Object.hasOwn(SubPosts, "links"))
			throw error(404, "Subreddit Not found");

		const { data: SubInfo } = await api.get("/teddit/r/%s/about", {
			params: [subredditName]
		});

		return {
			Info: SubInfo,
			Feed: SubPosts.links
		};
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
