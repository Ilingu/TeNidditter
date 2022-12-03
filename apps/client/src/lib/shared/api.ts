import { PUBLIC_API_URL } from "$env/static/public";

import { IsEmptyString, isValidUrl } from "./utils";
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
import type { FunctionJob } from "./types/globals";

interface QueryParams {
	uri: string;
	method: "GET" | "POST" | "PUT" | "DELETE";
	body?: object;
	headers?: object;
	credentials?: boolean;
}
interface APIResShape<T = never> {
	success: boolean;
	code: number;
	data: T;
	error: string;
}

/* API Client */

/**
 * @classdesc Tweniditter backend API client library (fully typed)
 * @class
 */
export default class api {
	/**
	 * all the `GET` endpoints
	 * @static
	 * @param {T} uri
	 * @param {GetParams<T>} additionals arguments of this route
	 * @returns {Promise<ApiClientResp<GetReturns<T, U>>>} the api response datas
	 */
	static async get<T extends GetRoutes, U>(
		uri: T,
		{ query, headers, params }: GetParams<T>,
		customFetch?: typeof fetch
	): Promise<ApiClientResp<GetReturns<T, U>>> {
		uri = BuildURI<T>(uri, { params, query });
		return await callApi<GetReturns<T, U>>(
			{
				uri,
				method: "GET",
				headers
			},
			customFetch
		);
	}

	/**
	 * all the `POST` endpoints
	 * @static
	 * @param {T} uri
	 * @param {GetParams<T>} additionals arguments of this route
	 * @returns {Promise<ApiClientResp<PostReturns<T>>>} the api response datas
	 */
	static async post<T extends PostRoutes>(
		uri: T,
		{ body, headers, params, query, credentials }: PostParams<T>
	): Promise<ApiClientResp<PostReturns<T>>> {
		uri = BuildURI<T>(uri, { params, query });
		return await callApi<PostReturns<T>>({
			uri,
			method: "POST",
			body,
			headers,
			credentials
		});
	}

	/**
	 * all the `PUT` endpoints
	 * @static
	 * @param {T} uri
	 * @param {GetParams<T>} additionals arguments of this route
	 * @returns {Promise<ApiClientResp<PutParams<T>>>} the api response datas
	 */
	static async put<T extends PutRoutes>(
		uri: T,
		{ body, headers, params, query }: PutParams<T>
	): Promise<ApiClientResp<PutReturns<T>>> {
		uri = BuildURI<T>(uri, { params, query });
		return await callApi<PutReturns<T>>({
			uri,
			method: "PUT",
			body,
			headers
		});
	}

	/**
	 * all the `DELETE` endpoints
	 * @static
	 * @param {T} uri
	 * @param {GetParams<T>} additionals arguments of this route
	 * @returns {Promise<ApiClientResp<DeleteParams<T>>>} the api response datas
	 */
	static async delete<T extends DeleteRoutes>(
		uri: T,
		{ body, headers, params, query, credentials }: DeleteParams<T>
	): Promise<ApiClientResp<DeleteReturns<T>>> {
		uri = BuildURI<T>(uri, { params, query });
		return await callApi<DeleteReturns<T>>({
			uri,
			method: "DELETE",
			body,
			headers,
			credentials
		});
	}
}

/* Helpers */
interface ApiClientResp<T = never> extends FunctionJob<T> {
	headers?: Headers;
	status?: number;
}

/**
 * under the hood function to call and handle the API repsonses (not typed)
 * @param {QueryParams} { uri, method, body, headers, credentials }
 * @returns {Promise<ApiClientResp<T>>} the api response datas (not typed)
 */
export const callApi = async <T = never>(
	{ uri, method, body, headers, credentials }: QueryParams,
	customFetch?: typeof fetch
): Promise<ApiClientResp<T>> => {
	if (IsEmptyString(uri)) return { success: false, error: "Invalid URI" };

	let url: string;
	if (uri.startsWith("/")) url = PUBLIC_API_URL + uri;
	else url = PUBLIC_API_URL + "/" + uri;

	if (!isValidUrl(url)) return { success: false, error: "Invalid URL" };

	try {
		const resp = await (customFetch || fetch)(url, {
			method,
			body: JSON.stringify(body),
			headers: { "Content-Type": "application/json", ...(headers || {}) },
			credentials: credentials ? "include" : undefined
		});
		if (resp.status === 204 || resp.status === 205) return { success: true, status: resp.status };

		const { success, data: apiRes, error }: APIResShape<T> = await resp.json();
		if (!resp.ok || !success) return { success: false, status: resp.status, error };

		return { success: true, data: apiRes, status: resp.status, headers: resp?.headers };
	} catch (err) {
		return { success: false, error: err as string };
	}
};

/**
 * It build the API route's URI from the given args
 * @param {T} uri
 * @returns {T} the constructed URI
 */
const BuildURI = <T extends string>(
	uri: T,
	{ params, query }: { params?: string[]; query?: object }
): T => {
	if (typeof uri !== "string") return `${uri}`;

	// Params subtitution
	params?.forEach((param) => (uri = uri.replace("%s", param) as T));

	// Queries
	if (query && Object.entries(query).length > 0) {
		(uri as string) += "?";
		for (const [key, val] of Object.entries(query)) {
			if (val === null || typeof val === "undefined") continue;
			(uri as string) += `${key}=${encodeURIComponent(val)}&`;
		}
		(uri as string) = uri.replace(/&+$/, "") as "/tedinitter/userInfo" | `/auth/available`; // trim last &
	}
	return uri;
};
