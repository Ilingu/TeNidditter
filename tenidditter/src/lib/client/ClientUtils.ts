import type { AlertShape, AlertTypes } from "./types/alerts";
import type { Themes } from "./types/themes";

/**
 * Whether the client is a mobile device or not
 * @returns {boolean} if `true`, it means that the client **is a mobile device**
 */
export const isMobile = (): boolean =>
	typeof navigator !== "undefined" &&
	/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Windows Phone|Mobile|SamsungBrowser|MiuiBrowser/i.test(
		navigator.userAgent
	);

/**
 * It change the daisyUI app theme in realtime
 * @param {Themes} theme
 */
export const changeAppTheme = (theme: Themes) =>
	document.documentElement.setAttribute("data-theme", theme);

/**
 * It display a new alert message to the user UI interface
 * @param {string} message - alert message
 * @param {AlertTypes} type - type of the alert `"success" | "info" | "warning" | "error"`
 * @param duration - amount of time the alert will be displayed in the user interface, **in millisecond** / `DEFAULT=5000`
 */
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
 * Uses the `Intl` API to Format a number to a human readable number string (e.g: `5000->"5k"`)
 * @param {number} num - the non formatted number
 * @returns {string} the formatted number string
 */
export const FormatNumbers = (num: number): string =>
	new Intl.NumberFormat(Intl.DateTimeFormat().resolvedOptions().locale, {
		notation: "compact"
	}).format(num);

/**
 * It takes 2 dates and returns the time elapsed in a human readable string
 *
 * e.g:
 *
 * ```js
 * const dateA = 1669395600093;
 * const dateB = 1669399200093;
 * FormatElapsedTime(dateA, dateB) // return: `1hr ago` (dateB-dateA = 3600000ms = 1hr)
 * ```
 * @param {number} dateA - dateA timestamp `dateA < dateB`
 * @param {number} dateB - dateB timestamp `dateB > dateA`
 * @returns {string} the formatted Elapsed Time
 */
export const FormatElapsedTime = (dateA: number, dateB: number): string => {
	const elapsedTime = dateB - dateA;

	const toMin = elapsedTime / 1000 / 60;
	if (toMin < 60) return `${Math.round(toMin)} min ago`;

	const toHour = toMin / 60;
	if (toHour < 24) return `${Math.round(toHour)}h ago`;

	const toDays = toHour / 24;
	if (toDays < 30.4375) return `${Math.round(toDays)}d ago`;

	const toMonths = toDays / 30.4375; // 1m = 30d or 31d or 28d or 29d --> 30.4375d is the avg
	if (toDays < 365.25) return `${Math.round(toMonths)}m ago`;

	const toYears = toDays / 365.25; // 1/4y has 365, so in avg 1y=365.25d (or 12m but it'll be less precise to use "toMonth" here)
	return `${Math.round(toYears)}y ago`;
};

/**
 * Copy Text To User Clipbord
 * @param {string} text
 */
export const copyToClipboard = (text: string) => navigator.clipboard.writeText(text);
