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
		const resp = await fetch(`https://teddit.net/r/${encodeURI(subredditName)}?api&raw_json=1`);
		if (!resp.ok) throw error(404, "Subreddit Not found");
		const datas: TedditSubredditPosts = await resp.json();

		const { data: SubInfo } = await api.get<SubRedditDatas["Info"]>({
			uri: "/teddit/r",
			param: subredditName
		});

		return { Info: SubInfo, Feed: datas?.links };
	} catch (err) {
		throw error(500, JSON.stringify(err));
	}
};
