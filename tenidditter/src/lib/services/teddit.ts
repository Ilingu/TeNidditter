import type { FeedTypeEnum } from "$lib/types/enums";
import type { FeedResult, TedditHomePageRes } from "$lib/types/interfaces";
import { isValidUrl } from "$lib/utils";

// export const QueryUserPost = () => void;

export const QueryHomePost = async (type: FeedTypeEnum, afterId?: string): Promise<FeedResult> => {
	if (type < 0 || type > 4) return { success: false };
	const TypeToWord = {
		0: "hot",
		1: "new",
		2: "top",
		3: "rising",
		4: "controversial"
	};

	try {
		const url = `https://teddit.net/r/all/${TypeToWord[type]}?api&raw_json=1${
			afterId ? `&t=&after=t3_${afterId}` : ""
		}`;
		if (!isValidUrl(url)) return { success: false };

		const resp = await fetch(url);
		if (!resp.ok) return { success: false, error: "No Post Retuned..." };

		const datas: TedditHomePageRes = await resp.json();
		if (typeof datas !== "object" || !Object.hasOwn(datas, "links"))
			return { success: false, error: "No Post Retuned..." };

		return { success: true, data: datas.links, type: "home_feed" };
	} catch (error) {
		return { success: false, error: error as string };
	}
};
