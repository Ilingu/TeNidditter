import type { NeetComment } from "$lib/shared/types/nitter";

export interface NeetInfo {
	/**
	 * The neet datas itself (and its context)
	 */
	main: NeetComment[];
	/**
	 * neet's Comments
	 */
	reply: NeetComment[][];
}
