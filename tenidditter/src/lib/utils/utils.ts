import { goto as naviguate } from "$app/navigation";
import type { Themes } from "$lib/types/types";

export const IsEmptyString = (str: string) => typeof str !== "string" || str.trim().length <= 0;
export const isValidUrl = (url: string): boolean => {
	try {
		new URL(url);
		return true;
	} catch {
		return false;
	}
};

export const pushRoute = (url: string) => {
	if (!isValidUrl(url)) return;
	naviguate(url, { noscroll: true });
};

export const changeAppTheme = (theme: Themes) =>
	document.documentElement.setAttribute("data-theme", theme);
