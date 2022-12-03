import type { FunctionJob } from "$lib/shared/types/globals";
import type { FeedType, TedditRawPost } from "$lib/shared/types/teddit";

export interface FeedResult extends FunctionJob<TedditRawPost[]> {
	type?: FeedType;
}
