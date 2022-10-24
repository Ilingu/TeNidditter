import type { FunctionJob, UserSubs } from "$lib/types/interfaces";
import { IsEmptyString, IsValidJSON } from "$lib/utils";
import { DecryptDatas, EncryptDatas } from "$lib/services/wasm/encryption";
import type { User } from "$lib/stores/auth";

/* LocalStorage */
export const GetJWT = async (): Promise<FunctionJob<string>> => {
	const rawToken = window.localStorage.getItem("JWT_TOKEN");
	if (!rawToken || IsEmptyString(rawToken)) return { success: false };

	const { success, data: Token } = DecryptDatas(rawToken);
	if (!success || !Token || Token?.length <= 0) return { success: false };

	return { success: true, data: Token };
};
export const SetJWT = async (JwtToken: string) => {
	const { success, data: eToken } = EncryptDatas(JwtToken);
	if (success && eToken && eToken?.length > 0) {
		window.localStorage.setItem("JWT_TOKEN", eToken);
		document.cookie = `JWT_TOKEN=${eToken}; expires=${new Date(
			Date.now() + 1000 * 60 * 60 * 24 * 90
		).toISOString()}; path=/`;
	}
};

export const GetLSUser = (): FunctionJob<User> => {
	const rawToken = window.localStorage.getItem("user");
	if (!rawToken || IsEmptyString(rawToken)) return { success: false };

	try {
		return { success: true, data: JSON.parse(rawToken) };
	} catch (err) {
		return { success: false };
	}
};

export const GetUserSubs = (): FunctionJob<UserSubs> => {
	const rawSubs = window.localStorage.getItem("subs");
	if (!rawSubs || IsEmptyString(rawSubs) || !IsValidJSON(rawSubs)) return { success: false };

	const subs: UserSubs = JSON.parse(rawSubs);
	if (typeof subs !== "object" || !Object.hasOwn(subs, "teddit") || !Object.hasOwn(subs, "nitter"))
		return { success: false };

	return { success: true, data: subs };
};
