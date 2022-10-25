import api from "$lib/api";
import type {
	TedditHomePageRes as TedditSubredditPosts,
	TedditRawPost
} from "$lib/types/interfaces";
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

export const load: import("./$types").PageServerLoad = async ({
	params
}): Promise<SubRedditDatas> => {
	const subredditName = params?.subteddit;
	if (IsEmptyString(subredditName)) throw error(404, "Not Found -- Invalid subreddit");

	try {
		const { success, data: SubPosts } = await api.get("/teddit/r", {
			param: subredditName + "/posts"
		});
		if (!success || typeof SubPosts !== "object" || !Object.hasOwn(SubPosts, "links"))
			throw error(404, "Subreddit Not found");

		const { data: SubInfo } = await api.get("/teddit/r", {
			param: subredditName + "/about"
		});

		return {
			Info: SubInfo as SubRedditDatas["Info"],
			Feed: (SubPosts as TedditSubredditPosts).links
		};
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
