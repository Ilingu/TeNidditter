import type { FeedResult } from "$lib/server/types/teddit";
import api from "$lib/shared/api";
import { QueryHomePost } from "$lib/shared/services/teddit";
import { FeedTypeEnum } from "$lib/shared/types/teddit";
import { IsEmptyString, MakeBearerToken } from "$lib/shared/utils";
import { error } from "@sveltejs/kit";

export const load: import("./$types").PageServerLoad = async ({
	cookies,
	fetch,
	url
}): Promise<FeedResult> => {
	const eToken = cookies.get("JwtToken");
	// if user authenticated then fetch his feed
	if (!eToken || IsEmptyString(eToken) || url.searchParams.get("type") === "home_feed")
		return fetchHomePage(fetch);

	// otherwise fetch the global teddit home posts
	const userFeed = await fetchUserFeed(eToken);
	if (!userFeed.success || !userFeed.data) return fetchHomePage(fetch);

	return userFeed;
};

const fetchUserFeed = async (JwtToken: string): Promise<FeedResult> => {
	const { success, data: Feed } = await api.get("/tedinitter/teddit/feed", {
		headers: MakeBearerToken(JwtToken)
	});
	if (!success || typeof Feed !== "object" || Feed.length <= 0) return { success: false };

	return { success: true, data: Feed, type: "user_feed" };
};

const fetchHomePage = async (customFetch: typeof fetch): Promise<FeedResult> => {
	try {
		const { success, data: posts } = await QueryHomePost(FeedTypeEnum.Hot, undefined, customFetch);
		if (!success) throw error(500, "No Posts returned...");

		if (typeof posts !== "object" || posts.length <= 0) throw error(500, "No Posts returned...");
		return { success: true, data: posts, type: "home_feed" };
	} catch (err) {
		throw error(500, "No Posts returned...");
	}
};

export const prerender = false;
