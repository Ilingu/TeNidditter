import type { FeedResult, TedditHomePageRes as TedditSubredditPosts } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

export const prerender = false;

export const load: import("./$types").PageServerLoad = async ({ params }): Promise<FeedResult> => {
	const subredditName = params?.subteddit;
	if (IsEmptyString(subredditName)) throw error(404, "Not Found -- Invalid subreddit");

	try {
		const resp = await fetch(`https://teddit.net/r/${encodeURI(subredditName)}?api&raw_json=1`);
		if (!resp.ok) throw error(404, "Subreddit Not found");
		const datas: TedditSubredditPosts = await resp.json();

		return { success: true, data: datas.links };
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
