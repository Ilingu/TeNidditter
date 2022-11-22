// export const QueryUserPost = () => void;

import type { FeedResult } from "$lib/server/types/teddit";
import api from "../api";
import type { FeedHomeType, FeedTypeEnum } from "../types/teddit";

export const QueryHomePost = async (
	type: FeedTypeEnum,
	afterId?: string,
	customFetch?: typeof fetch
): Promise<FeedResult> => {
	if (type < 0 || type > 4) return { success: false };
	const TypeToWord: Record<number, FeedHomeType> = {
		0: "hot",
		1: "new",
		2: "top",
		3: "rising",
		4: "controversial"
	};

	try {
		const { success, data: posts } = await api.get(
			"/teddit/home",
			{
				query: { type: TypeToWord[type], afterId }
			},
			customFetch
		);

		if (!success) return { success: false, error: "No Post Retuned..." };
		if (typeof posts !== "object" || !Object.hasOwn(posts, "links"))
			return { success: false, error: "No Post Retuned..." };

		return { success: true, data: posts.links, type: "home_feed" };
	} catch (error) {
		return { success: false, error: error as string };
	}
};
