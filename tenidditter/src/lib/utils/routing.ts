import { isValidUrl } from "./utils";
import { goto as naviguate } from "$app/navigation";

export const pushRoute = (url: string) => {
	if (!isValidUrl(url)) return;
	dispatchRouting();
	naviguate(url, { noscroll: true });
};

export const dispatchRouting = () => {
	const changeRoute = new CustomEvent("routing");
	document.dispatchEvent(changeRoute);
};

export const routingListener = (cb: () => void) => document.addEventListener("routing", cb);
