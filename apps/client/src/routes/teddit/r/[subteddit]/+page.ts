import api from "$lib/shared/api";
import type { TedditRawPost } from "$lib/shared/types/teddit";
import { IsEmptyString } from "$lib/shared/utils";
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

// fetch subteddit datas (posts and about page)
export const load: import("./$types").PageLoad = async ({
	params,
	fetch
}): Promise<SubRedditDatas> => {
	const subredditName = params?.subteddit;
	if (IsEmptyString(subredditName)) throw error(400, "Invalid subreddit");

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
};
