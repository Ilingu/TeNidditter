export const pushRoute = async (uri: string) => {
	const naviguate = (await import("$app/navigation")).goto;
	dispatchRouting();
	naviguate(uri);
};

/**
 * It dispatch the custom "routing" event to all the app, this is executed when a new route (internal to the app) is pushed by the `Link` components (thus it is faster than the premaid "beforeNavigation" hook from svelte)
 */
export const dispatchRouting = () => {
	const changeRoute = new CustomEvent("routing");
	document.dispatchEvent(changeRoute);
};

/**
 * It listen to the custom "routing" event, when trigger it means that the app will change it's route soonly
 * @param {() => void} cb - callback function that will be called when the event is triggered
 */
export const routeWillChange = (cb: () => void) => document.addEventListener("routing", cb);
