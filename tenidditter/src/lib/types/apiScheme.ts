import type { User } from "$lib/stores/auth";
import type { TedditHomePageRes, TedditPostInfo, TedditUserShape } from "./interfaces";
import type { FeedHomeType } from "./types";

/* GET */
export type GetRoutes =
	| "/tedinitter/userInfo"
	| "/auth/available"
	| "/teddit/r/%s/about"
	| "/teddit/r/%s/posts"
	| "/teddit/u/%s"
	| "/teddit/home"
	| "/teddit/r/%s/post/%s";
export type GetReturns<T> = T extends "/tedinitter/userInfo"
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
	: never;
export interface GetParams<T> {
	query?: T extends "/auth/available"
		? { username: string }
		: T extends "/teddit/home"
		? { type?: FeedHomeType; afterId?: string }
		: never;
	headers?: T extends "/tedinitter/userInfo" ? { Authorization: string } : never;
	params?: T extends "/teddit/r/%s/about"
		? [subteddit: string]
		: T extends "/teddit/r/%s/posts"
		? [subteddit: string]
		: T extends "/teddit/u/%s"
		? [username: string]
		: T extends "/teddit/r/%s/post/%s"
		? [subteddit: string, postId: string]
		: never;
}

/* POST */
export type PostRoutes = "/auth/" | "/tedinitter/teddit/sub/%s";
export type PostReturns<T> = T extends "/auth/"
	? string
	: T extends "/tedinitter/teddit/sub/%s"
	? null
	: never;
export interface PostParams<T> {
	query?: never;
	headers?: T extends "/tedinitter/teddit/sub/%s" ? { Authorization: string } : never;
	params?: T extends "/tedinitter/teddit/sub/%s" ? [subteddit: string] : never;
	body?: T extends "/auth/" ? { username: string; password: string } : never;
}

/* PUT */
export type PutRoutes = never;
export type PutReturns = never;
export interface PutParams {
	query?: never;
	headers?: never;
	params?: never;
	body?: never;
}

/* DELETE */
export type DeleteRoutes = "/tedinitter/teddit/unsub/%s";
export type DeleteReturns<T> = T extends "/tedinitter/teddit/unsub/%s" ? null : never;
export interface DeleteParams<T> {
	query?: never;
	headers?: T extends "/tedinitter/teddit/unsub/%s" ? { Authorization: string } : never;
	params?: T extends "/tedinitter/teddit/unsub/%s" ? [subteddit: string] : never;
	body?: never;
}
