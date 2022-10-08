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

export const humanElapsedTime = (dateA: number, dateB: number): string => {
	const elapsedTime = dateB - dateA;

	const toMin = elapsedTime / 1000 / 60;
	if (toMin < 60) return `${Math.round(toMin)} min ago`;

	const toHour = toMin / 60;
	if (toHour < 24) return `${Math.round(toHour)}h ago`;

	const toDays = toHour / 24;
	if (toDays < 30) return `${Math.round(toDays)}d ago`;

	const toMonths = toDays / 30.4375; // 1m = 30d or 31d or 28d or 29d --> 30.4375d is the avg
	if (toMonths < 12) return `${Math.round(toMonths)}m ago`;

	const toYears = toDays / 365.25; // 1/4y has 365, so in avg 1y=365.25d (or 12m but it'll be less precise to use "toMonth" here)
	return `${Math.round(toYears)}y ago`;
};
