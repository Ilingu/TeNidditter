import { PUBLIC_API_URL } from "$env/static/public";

import { IsEmptyString, isValidUrl } from "./utils";
import type { FunctionJob } from "$lib/types/interfaces";

interface QueryParams {
	uri: string;
	method: "GET" | "POST" | "PUT" | "DELETE";
	body?: object;
	headers?: object;
}
interface APIResShape<T = never> {
	success: boolean;
	code: number;
	data: T;
}

interface APIClientParams<T extends { route: string; body?: object; headers?: object }> {
	uri: T["route"];
	body?: T["body"];
	headers?: T["headers"];
	query?: {
		[key: string]: string;
	};
}

type GetType<T = never> = Omit<T, "body">;

export default class api {
	static async get<T = never>({
		uri,
		query,
		headers
	}: GetType<
		APIClientParams<
			| { route: "/tedinitter/userInfo"; headers: { Authorization: string } }
			| { route: `/auth/available`; headers: undefined }
		>
	>): Promise<FunctionJob<T>> {
		if (query && Object.entries(query).length > 0) {
			uri += "?";
			for (const [key, val] of Object.entries(query)) {
				uri += `${key}=${encodeURI(val)}&`;
			}
			uri = uri.replace(/&+$/, "") as "/tedinitter/userInfo" | `/auth/available`; // trim last &
		}

		return await callApi<T>({
			uri,
			method: "GET",
			headers: headers as object
		});
	}
	static async post<T = never>({
		uri,
		body,
		headers
	}: APIClientParams<{ route: "/auth/"; body: { username: string; password: string } }>): Promise<
		FunctionJob<T>
	> {
		return await callApi<T>({
			uri,
			method: "POST",
			body,
			headers: headers as object
		});
	}
	static update({ uri, body, headers }: APIClientParams<{ route: "" }>) {
		console.log({ uri, body, headers });
	}
	static delete({ uri, body, headers }: APIClientParams<{ route: "" }>) {
		console.log({ uri, body, headers });
	}
}

export const callApi = async <T = never>({
	uri,
	method,
	body,
	headers
}: QueryParams): Promise<FunctionJob<T>> => {
	if (IsEmptyString(uri)) return { success: false, error: "Invalid URI" };

	let url: string;
	if (uri.startsWith("/")) url = PUBLIC_API_URL + uri;
	else url = PUBLIC_API_URL + "/" + uri;

	if (!isValidUrl(url)) return { success: false, error: "Invalid URL" };

	try {
		const resp = await fetch(url, {
			method,
			body: JSON.stringify(body),
			headers: { "Content-Type": "application/json", ...(headers || {}) }
		});
		if (!resp.ok) return { success: false, error: "Request Failed" };

		const { success, data: apiRes }: APIResShape<T> = await resp.json();
		if (!success) return { success: false, error: "apiRes" };

		return { success: true, data: apiRes };
	} catch (error) {
		return { success: false, error: error as string };
	}
};
