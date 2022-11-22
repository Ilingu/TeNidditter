import type { NeetComment } from "$lib/shared/types/nitter";

export interface NeetInfo {
	main: NeetComment[];
	reply: NeetComment[][];
}
