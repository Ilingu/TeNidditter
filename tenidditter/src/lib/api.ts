import { PUBLIC_API_URL } from "$env/static/public";

import { IsEmptyString, isValidUrl } from "./utils";
import type { FunctionJob } from "$lib/types/interfaces";
import type {
	// GET type
	GetParams,
	GetReturns,
	GetRoutes,
	// POST type
	PostParams,
	PostReturns,
	PostRoutes,
	// PUT type
	PutParams,
	PutReturns,
	PutRoutes,
	// DELETE type
	DeleteParams,
	DeleteReturns,
	DeleteRoutes
} from "./types/apiScheme";

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

/* API Client */
export default class api {
	static async get<T extends GetRoutes>(
		uri: T,
		{ query, headers, param }: GetParams<T>
	): Promise<ApiClientResp<GetReturns<T>>> {
		uri = BuildURI<T>(uri, { param, query });
		return await callApi<GetReturns<T>>({
			uri,
			method: "GET",
			headers
		});
	}

	static async post<T extends PostRoutes>(
		uri: T,
		{ body, headers, param, query }: PostParams<T>
	): Promise<ApiClientResp<PostReturns<T>>> {
		uri = BuildURI<T>(uri, { param, query });
		return await callApi<PostReturns<T>>({
			uri,
			method: "POST",
			body,
			headers
		});
	}

	static async update<T extends PutRoutes>(
		uri: T,
		{ body, headers, param, query }: PutParams
	): Promise<ApiClientResp<PutReturns>> {
		uri = BuildURI<T>(uri, { param, query });
		return await callApi<PutReturns>({
			uri,
			method: "PUT",
			body,
			headers
		});
	}

	static async delete<T extends DeleteRoutes>(
		uri: T,
		{ body, headers, param, query }: DeleteParams<T>
	): Promise<ApiClientResp<DeleteReturns<T>>> {
		uri = BuildURI<T>(uri, { param, query });
		return await callApi<DeleteReturns<T>>({
			uri,
			method: "DELETE",
			body,
			headers
		});
	}
}

/* Helpers */
interface ApiClientResp<T = never> extends FunctionJob<T> {
	headers?: Headers;
}

export const callApi = async <T = never>({
	uri,
	method,
	body,
	headers
}: QueryParams): Promise<ApiClientResp<T>> => {
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

		return { success: true, data: apiRes, headers: resp?.headers };
	} catch (error) {
		return { success: false, error: error as string };
	}
};

const BuildURI = <T extends string>(
	uri: T,
	{ param, query }: { param?: string; query?: object }
): T => {
	if (typeof uri !== "string") return `${uri}`;

	if (!IsEmptyString(param)) (uri as string) += `/${param}`;
	if (query && Object.entries(query).length > 0) {
		(uri as string) += "?";
		for (const [key, val] of Object.entries(query)) {
			if (val === null || typeof val === "undefined") continue;
			(uri as string) += `${key}=${encodeURI(val)}&`;
		}
		(uri as string) = uri.replace(/&+$/, "") as "/tedinitter/userInfo" | `/auth/available`; // trim last &
	}
	return uri;
};
