import { goto as naviguate } from "$app/navigation";
import { isValidUrl } from "$lib/utils";

export const pushRoute = (url: string) => {
	if (!isValidUrl(url)) return;

	dispatchRouting();
	naviguate(url, { noscroll: true });
};

// Faster than beforeNavigation
export const dispatchRouting = () => {
	const changeRoute = new CustomEvent("routing");
	document.dispatchEvent(changeRoute);
};

export const routeWillChange = (cb: () => void) => document.addEventListener("routing", cb);
