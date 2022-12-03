import api from "$lib/shared/api";
import type { FunctionJob } from "$lib/shared/types/globals";
import { MakeBearerToken } from "$lib/shared/utils";
import { pushAlert } from "../ClientUtils";
import { LogOut } from "./auth";

/**
 * [user action]: Wrapper function that deletes the currently logged in user
 * @param {string} JwtToken - token that identifies the user
 */
export const DeleteUserAccount = async (JwtToken: string) => {
	pushAlert("Deleting your account...", "info", 1500);

	const { success: deleted, error } = await api.delete("/auth/", {
		credentials: true,
		query: { token: JwtToken }
	});
	if (!deleted) return pushAlert(`Failed to delete your account: ${error}`, "error", 6000);

	pushAlert("This account vanished mysteriously 👻, Bye~~", "success", 5000);
	LogOut();
};

/**
 * [user action]: Wrapper function that **regenerates** all the user's recovery codes (it **overwrites** the old one), the new codes should be displayed to the user
 * @param {string} JwtToken - token that identifies the user
 */
export const RegenerateUserRecoveryCodes = async (
	JwtToken: string
): Promise<FunctionJob<string[]>> => {
	const { success, data: newCodes } = await api.put("/tedinitter/regererate-recovery-codes", {
		headers: MakeBearerToken(JwtToken)
	});
	if (!success || !newCodes || typeof newCodes !== "object" || newCodes?.length !== 6)
		return { success: false };

	return { success: true, data: newCodes };
};

/**
 * [user action]: Wrapper function that handle the user's teddit subs, if already sub to `subreddit` it'll unsub him and vice versa
 * @param {string} subteddit - the subreddit name
 * @param {boolean} isSub - is the user already subscribed to this subreddit?
 * @param {string} JwtToken - token that identifies the user
 * @returns {Promise<FunctionJob>} return the success (or not) of the operation
 */
export const ToggleTedditSubs = async (
	subteddit: string,
	isSub: boolean,
	JwtToken: string
): Promise<FunctionJob> => {
	if (isSub) {
		const { success } = await api.delete("/tedinitter/teddit/unsub/%s", {
			params: [subteddit],
			headers: MakeBearerToken(JwtToken || "")
		});
		if (!success) pushAlert("Couldn't unsubscribe you, try again", "error");
		return { success };
	}

	const { success } = await api.post("/tedinitter/teddit/sub/%s", {
		params: [subteddit],
		headers: MakeBearerToken(JwtToken || "")
	});
	if (!success) pushAlert("Couldn't subscribe you, try again", "error");
	return { success };
};

/**
 * [user action]: Wrapper function that handle the user's nittos subs, if already sub to `nittos` it'll unsub him and vice versa
 * @param {string} nittos - the nittos name
 * @param {boolean} isSub - is the user already subscribed to this nittos?
 * @param {string} JwtToken - token that identifies the user
 * @returns {Promise<FunctionJob>} return the success (or not) of the operation
 */
export const ToggleNitterSubs = async (
	nittos: string,
	isSub: boolean,
	JwtToken: string
): Promise<FunctionJob> => {
	if (isSub) {
		isSub = false;
		const { success } = await api.delete("/tedinitter/nitter/unsub/%s", {
			params: [nittos],
			headers: MakeBearerToken(JwtToken)
		});
		if (!success) pushAlert("couldn't unsubscribe you", "error");

		return { success };
	}

	isSub = true;
	const { success } = await api.post("/tedinitter/nitter/sub/%s", {
		params: [nittos],
		headers: MakeBearerToken(JwtToken)
	});
	if (!success) pushAlert("couldn't subscribe you", "error");

	return { success };
};
