import type { FunctionJob, NitterLists as UserLists, UserSubs } from "$lib/types/interfaces";
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
	if (success && eToken && eToken?.length > 0) window.localStorage.setItem("JWT_TOKEN", eToken);
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

	const subs: UserSubs = JSON.parse(rawSubs) ?? undefined;
	if (typeof subs !== "object" || !Object.hasOwn(subs, "teddit") || !Object.hasOwn(subs, "nitter"))
		return { success: false };

	return { success: true, data: subs };
};

export const GetUserLists = (): FunctionJob<UserLists[]> => {
	const rawList = window.localStorage.getItem("lists");
	if (!rawList || IsEmptyString(rawList) || !IsValidJSON(rawList)) return { success: false };

	const lists: UserLists[] = JSON.parse(rawList) ?? undefined;
	if (typeof lists !== "object" || lists.length <= 0) return { success: false };

	return { success: true, data: lists };
};
