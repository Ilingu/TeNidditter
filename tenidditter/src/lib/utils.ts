import type { AlertShape } from "$lib/types/interfaces";
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

export const ConvertHTMLEntities = (str: string): string => {
	const htmlEntities = {
		"&": "&amp;",
		"<": "&lt;",
		">": "&gt;",
		'"': "&quot;",
		"'": "&apos;"
	};
	return str.replace(/([&<>"'])/g, (match) => htmlEntities[match as keyof typeof htmlEntities]);
};

export const FormatNumbers = (num: number) =>
	new Intl.NumberFormat(Intl.DateTimeFormat().resolvedOptions().locale, {
		notation: "compact"
	}).format(num);
