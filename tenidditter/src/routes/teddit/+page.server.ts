import { QueryHomePost } from "$lib/services/teddit";
import { FeedTypeEnum } from "$lib/types/enums";
import type { FeedResult } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";

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
		const { success, data: posts } = await QueryHomePost(FeedTypeEnum.Hot);
		if (!success) throw error(500, "No Posts returned...");

		if (typeof posts !== "object" || posts.length <= 0) throw error(500, "No Posts returned...");
		return { success: true, data: posts, type: "home_feed" };
	} catch (err) {
		throw error(500, "No Posts returned...");
	}
};

export const prerender = false;
