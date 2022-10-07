import type { TedditRawPost, TedditHomePageRes, FunctionJob } from "$lib/types/interfaces";
import type { FeedType } from "$lib/types/types";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

interface FeedResult extends FunctionJob<TedditRawPost[]> {
	type: FeedType;
}

export const load: import("./$types").PageServerLoad = async ({ cookies }): Promise<FeedResult> => {
	const eToken = cookies.get("JWT_TOKEN");
	if (!eToken || IsEmptyString(eToken)) return fetchHomePage();

	const userFeed = await fetchUserFeed();
	if (!userFeed.success || !userFeed.data) return fetchHomePage();

	return userFeed;
};

const fetchUserFeed = async (): Promise<FeedResult> => {
	return { success: false, type: "user_feed" };
};

const fetchHomePage = async (): Promise<FeedResult> => {
	try {
		const resp = await fetch("https://teddit.net/?api&raw_json=1");
		if (!resp.ok) throw error(500, "No Posts returned...");

		const datas: TedditHomePageRes = await resp.json();
		if (typeof datas !== "object" || !Object.hasOwn(datas, "links"))
			throw error(500, "No Posts returned...");

		return { success: true, data: datas.links, type: "home_feed" };
	} catch (err) {
		throw error(500, "No Posts returned...");
	}
};
