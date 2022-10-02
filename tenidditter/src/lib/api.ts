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

export default class ApiClient {
	static get() {
		//
	}
	static post() {
		//
	}
	static update() {
		//
	}
	static delete() {
		//
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
