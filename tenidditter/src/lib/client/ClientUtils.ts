import type { AlertShape, AlertTypes } from "./types/alerts";
import type { Themes } from "./types/themes";

export const isMobile = () =>
	typeof navigator !== "undefined" &&
	/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Windows Phone/i.test(
		navigator.userAgent
	);

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

/**
 * Copy Text To User Clipbord
 * @param {string} text
 */
export const copyToClipboard = (text: string) => navigator.clipboard.writeText(text);
