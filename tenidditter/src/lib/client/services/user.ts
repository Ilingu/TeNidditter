import api from "$lib/shared/api";
import type { FunctionJob } from "$lib/shared/types/globals";
import { MakeBearerToken } from "$lib/shared/utils";
import { pushAlert } from "../ClientUtils";
import { LogOut } from "./auth";

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
