import { goto as naviguate } from "$app/navigation";
import { isValidUrl } from "$lib/shared/utils";

export const pushRoute = (url: string) => {
	if (!isValidUrl(url)) return;

	dispatchRouting();
	naviguate(url);
};

// Faster than beforeNavigation
export const dispatchRouting = () => {
	const changeRoute = new CustomEvent("routing");
	document.dispatchEvent(changeRoute);
};

export const routeWillChange = (cb: () => void) => document.addEventListener("routing", cb);
