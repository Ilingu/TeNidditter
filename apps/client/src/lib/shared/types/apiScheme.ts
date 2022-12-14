import type { User } from "$lib/client/types/auth";
import type { DBSubtedditsShape } from "$lib/client/types/teddit";
import type { NeetInfo } from "$lib/server/types/nitter";
import type { NeetComment, Nittos, NittosPreview } from "./nitter";
import type {
	FeedHomeType,
	TedditHomePageRes,
	TedditPostInfo,
	TedditRawPost,
	TedditUserShape
} from "./teddit";

/* GET */
export type GetRoutes =
	| "/tedinitter/userInfo"
	| "/auth/available"
	| "/teddit/r/%s/about"
	| "/teddit/r/%s/posts"
	| "/teddit/u/%s"
	| "/teddit/home"
	| "/teddit/r/%s/post/%s"
	| "/teddit/r/search"
	| "/tedinitter/teddit/feed"
	| "/tedinitter/nitter/feed"
	| "/nitter/search"
	| "/nitter/nittos/%s/about"
	| "/nitter/nittos/%s/neets"
	| "/nitter/nittos/%s/neets/%s"
	| "/tedinitter/nitter/list/%s";
export type GetReturns<T, U> = T extends "/tedinitter/userInfo"
	? User
	: T extends "/auth/available"
	? boolean
	: T extends "/teddit/r/%s/about"
	? {
			subs: string;
			description: string;
			rules: string;
	  }
	: T extends "/teddit/r/%s/posts"
	? TedditHomePageRes
	: T extends "/teddit/u/%s"
	? TedditUserShape
	: T extends "/teddit/home"
	? TedditHomePageRes
	: T extends "/teddit/r/%s/post/%s"
	? TedditPostInfo
	: T extends "/teddit/r/search"
	? DBSubtedditsShape[]
	: T extends "/tedinitter/teddit/feed"
	? TedditRawPost[]
	: T extends "/tedinitter/nitter/feed"
	? NeetComment[][]
	: T extends "/nitter/search"
	? U extends "tweets"
		? NeetComment[][]
		: U extends "users"
		? NittosPreview[]
		: NeetComment[][] | NittosPreview[]
	: T extends "/nitter/nittos/%s/about"
	? Nittos
	: T extends "/nitter/nittos/%s/neets"
	? NeetComment[][]
	: T extends "/nitter/nittos/%s/neets/%s"
	? NeetInfo
	: T extends "/tedinitter/nitter/list/%s"
	? NeetComment[]
	: never;
export interface GetParams<T> {
	query?: T extends "/auth/available"
		? { username: string }
		: T extends "/teddit/home"
		? { type?: FeedHomeType; afterId?: string }
		: T extends "/teddit/r/%s/post/%s"
		? { sort: string }
		: T extends "/teddit/r/search"
		? { q: string }
		: T extends "/nitter/search"
		? { type: "users" | "tweets"; q: string; limit?: number }
		: T extends "/nitter/nittos/%s/neets"
		? { limit?: number }
		: T extends "/nitter/nittos/%s/neets/%s"
		? { limit?: number }
		: never;
	headers?: T extends "/tedinitter/userInfo"
		? { Authorization: string }
		: T extends "/tedinitter/teddit/feed"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/feed"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/list/%s"
		? { Authorization: string }
		: never;
	params?: T extends "/teddit/r/%s/about"
		? [subteddit: string]
		: T extends "/teddit/r/%s/posts"
		? [subteddit: string]
		: T extends "/teddit/u/%s"
		? [username: string]
		: T extends "/teddit/r/%s/post/%s"
		? [subteddit: string, postId: string]
		: T extends "/nitter/nittos/%s/about"
		? [nittosname: string]
		: T extends "/nitter/nittos/%s/neets"
		? [nittosname: string]
		: T extends "/nitter/nittos/%s/neets/%s"
		? [nittosname: string, neetId: string]
		: T extends "/tedinitter/nitter/list/%s"
		? [listId: string]
		: never;
}

/* POST */
export type PostRoutes =
	| "/auth/"
	| "/tedinitter/teddit/sub/%s"
	| "/tedinitter/nitter/sub/%s"
	| "/tedinitter/nitter/list"
	| "/tedinitter/nitter/list/%s/saveNeet";
export type PostReturns<T> = T extends "/auth/"
	? string
	: T extends "/tedinitter/teddit/sub/%s"
	? null
	: T extends "/tedinitter/nitter/sub/%s"
	? null
	: T extends "/tedinitter/nitter/list"
	? null
	: T extends "/tedinitter/nitter/list/%s/saveNeet"
	? null
	: never;
export interface PostParams<T> {
	query?: never;
	headers?: T extends "/tedinitter/teddit/sub/%s"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/sub/%s"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/list"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/list/%s/saveNeet"
		? { Authorization: string }
		: never;
	params?: T extends "/tedinitter/teddit/sub/%s"
		? [subteddit: string]
		: T extends "/tedinitter/nitter/sub/%s"
		? [nittos: string]
		: T extends "/tedinitter/nitter/list/%s/saveNeet"
		? [listId: string]
		: never;
	body?: T extends "/auth/"
		? { username: string; password: string; method: "login" | "register" }
		: T extends "/tedinitter/nitter/list"
		? { listname: string }
		: T extends "/tedinitter/nitter/list/%s/saveNeet"
		? NeetComment
		: never;
	credentials?: T extends "/auth/" ? true : never;
}

/* PUT */
export type PutRoutes = "/auth/reset-password" | "/tedinitter/regererate-recovery-codes";
export type PutReturns<T> = T extends "/auth/reset-password"
	? null
	: T extends "/tedinitter/regererate-recovery-codes"
	? string[]
	: never;
export interface PutParams<T> {
	query?: never;
	headers?: T extends "/tedinitter/regererate-recovery-codes" ? { Authorization: string } : never;
	params?: never;
	body?: T extends "/auth/reset-password"
		? { username: string; NewPassword: string; RecoveryCode: string }
		: never;
}

/* DELETE */
export type DeleteRoutes =
	| "/tedinitter/teddit/unsub/%s"
	| "/tedinitter/nitter/unsub/%s"
	| "/tedinitter/nitter/list/%s"
	| "/tedinitter/nitter/list/%s/removeNeet/%s"
	| "/auth/logout"
	| "/auth/";
export type DeleteReturns<T> = T extends "/tedinitter/teddit/unsub/%s"
	? null
	: T extends "/tedinitter/nitter/unsub/%s"
	? null
	: T extends "/tedinitter/nitter/list/%s"
	? null
	: T extends "/tedinitter/nitter/list/%s/removeNeet/%s"
	? null
	: never;
export interface DeleteParams<T> {
	query?: T extends "/auth/logout"
		? { token?: string }
		: T extends "/auth/"
		? { token: string }
		: never;
	headers?: T extends "/tedinitter/teddit/unsub/%s"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/unsub/%s"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/list/%s"
		? { Authorization: string }
		: T extends "/tedinitter/nitter/list/%s/removeNeet/%s"
		? { Authorization: string }
		: never;
	params?: T extends "/tedinitter/teddit/unsub/%s"
		? [subteddit: string]
		: T extends "/tedinitter/nitter/unsub/%s"
		? [nittos: string]
		: T extends "/tedinitter/nitter/list/%s"
		? [listId: string]
		: T extends "/tedinitter/nitter/list/%s/removeNeet/%s"
		? [listId: string, neetId: string]
		: never;
	body?: never;
	credentials?: T extends "/auth/logout" ? true : T extends "/auth/" ? true : never;
}
