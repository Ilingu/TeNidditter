import { IsEmptyString } from "$lib/shared/utils";
import { error } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const prerender = false;

export const actions: Actions = {
	default: async ({ request }): Promise<boolean> => {
		const form = await request.formData();
		const username = form.get("subteddit");

		if (IsEmptyString(username)) throw error(400, "bad args");
		try {
			const resp = await fetch(`https://teddit.net/r/${encodeURI(username as string)}?api`);
			if (!resp.ok) return false;

			const sub = await resp.json();
			if (typeof sub !== "object" || !Object.hasOwn(sub, "links") || sub?.links?.length <= 0)
				return false;

			return true;
		} catch (err) {
			throw error(500, "couldn't check if user exist");
		}
	}
};
