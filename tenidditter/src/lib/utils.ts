import type { AlertShape, FunctionJob } from "$lib/types/interfaces";
import type { AlertTypes, Themes } from "$lib/types/types";

export const IsEmptyString = (str: unknown) => typeof str !== "string" || str.trim().length <= 0;
export const isValidUrl = (url: string): boolean => {
	try {
		new URL(url);
		return true;
	} catch {
		return false;
	}
};
export const isMobile = () =>
	/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Windows Phone/i.test(
		navigator.userAgent
	);

export const Sleep = (dur: number) => new Promise((res) => setTimeout(res, dur));

export const changeAppTheme = (theme: Themes) =>
	document.documentElement.setAttribute("data-theme", theme);

export const pushAlert = (message: string, type: AlertTypes, duration = 5000) => {
	const alert = new CustomEvent("alert", {
		detail: {
			message,
			type,
			duration
		} as AlertShape
	});
	document.dispatchEvent(alert);
};

/**
 * Remove all non-letter in username except "_"
 * @param {string} username
 * @returns {string} the formatted username
 */
export const FormatUsername = (username: string): string =>
	username.replace(/[\W0-9]+/g, "").toLowerCase();

/**
 * Encrypt data
 * @param {string} datas
 */
export const encryptDatas = async (datas: string): Promise<FunctionJob<string>> => {
	try {
		const response = await fetch(`/api/encrypt?enc=${encodeURIComponent(datas)}`, {
			method: "GET"
		});
		if (!response.ok) return { success: false };

		const encryptedDatas = decodeURIComponent(await response.text());
		return { success: true, data: encryptedDatas };
	} catch (err) {
		return { success: false };
	}
};

export const decryptDatas = async (datas: string): Promise<FunctionJob<string>> => {
	try {
		const response = await fetch(`/api/decrypt?dec=${encodeURIComponent(datas)}`, {
			method: "GET"
		});
		if (!response.ok) return { success: false };

		const decryptedDatas = decodeURIComponent(await response.text());
		return { success: true, data: decryptedDatas };
	} catch (err) {
		return { success: false };
	}
};
