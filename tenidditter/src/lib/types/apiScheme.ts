import type { User } from "$lib/stores/auth";
import type { TedditHomePageRes, TedditUserShape } from "./interfaces";
import type { FeedHomeType } from "./types";

/* GET */
export type GetRoutes =
	| "/tedinitter/userInfo"
	| "/auth/available"
	| "/teddit/r"
	| "/teddit/u"
	| "/teddit/home";
export type GetReturns<T> = T extends "/tedinitter/userInfo"
	? User
	: T extends "/auth/available"
	? boolean
	: T extends "/teddit/r"
	?
			| TedditHomePageRes
			| {
					subs: string;
					description: string;
					rules: string;
			  }
	: T extends "/teddit/u"
	? TedditUserShape
	: T extends "/teddit/home"
	? TedditHomePageRes
	: never;
export interface GetParams<T> {
	query?: T extends "/auth/available"
		? { username: string }
		: T extends "/teddit/home"
		? { type?: FeedHomeType; afterId?: string }
		: never;
	headers?: T extends "/tedinitter/userInfo" ? { Authorization: string } : never;
	param?: T extends "/teddit/r" ? string : T extends "/teddit/u" ? string : never;
}

/* POST */
export type PostRoutes = "/auth/" | "/tedinitter/teddit/sub";
export type PostReturns<T> = T extends "/auth/"
	? string
	: T extends "/tedinitter/teddit/sub"
	? null
	: never;
export interface PostParams<T> {
	query?: never;
	headers?: T extends "/tedinitter/teddit/sub" ? { Authorization: string } : never;
	param?: T extends "/tedinitter/teddit/sub" ? string : never;
	body?: T extends "/auth/" ? { username: string; password: string } : never;
}

/* PUT */
export type PutRoutes = never;
export type PutReturns = never;
export interface PutParams {
	query?: never;
	headers?: never;
	param?: never;
	body?: never;
}

/* DELETE */
export type DeleteRoutes = "/tedinitter/teddit/unsub";
export type DeleteReturns<T> = T extends "/tedinitter/teddit/unsub" ? null : never;
export interface DeleteParams<T> {
	query?: never;
	headers?: T extends "/tedinitter/teddit/unsub" ? { Authorization: string } : never;
	param?: T extends "/tedinitter/teddit/unsub" ? string : never;
	body?: never;
}
