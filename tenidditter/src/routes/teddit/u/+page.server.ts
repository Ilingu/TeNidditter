import { IsEmptyString } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const prerender = false;

export const actions: Actions = {
	default: async ({ request }): Promise<boolean> => {
		const form = await request.formData();
		const username = form.get("username");

		if (IsEmptyString(username)) throw error(400, "bad args");
		try {
			const resp = await fetch(`https://teddit.net/u/${encodeURI(username as string)}?api`);
			if (!resp.ok) return false;

			const htmlPage = await resp.text();
			if (htmlPage.includes("reddit-error")) return false;

			return true;
		} catch (err) {
			throw error(500, "couldn't check if user exist");
		}
	}
};
