import { DecryptDatas, EncryptDatas } from "$lib/client/services/wasm/encryption";
import type { FunctionJob } from "$lib/shared/types/globals";
import { IsEmptyString, IsValidJSON } from "$lib/shared/utils";
import type { User, UserSubs } from "../types/auth";
import type { NitterLists } from "../types/nitter";

/* LocalStorage */

/**
 * it fetch, decrypt, check from localstorage the user's jwtToken
 * @returns {Promise<FunctionJob<string>>} the user jwtToken
 */
export const GetJWT = async (): Promise<FunctionJob<string>> => {
	const rawToken = window.localStorage.getItem("JWT_TOKEN");
	if (!rawToken || IsEmptyString(rawToken)) return { success: false };

	const { success, data: Token } = DecryptDatas(rawToken);
	if (!success || !Token || Token?.length <= 0) return { success: false };

	return { success: true, data: Token };
};
/**
 * it encryptand set the user's jwtToken to localstorage
 * @param {string} JwtToken
 */
export const SetJWT = async (JwtToken: string) => {
	const { success, data: eToken } = EncryptDatas(JwtToken);
	if (success && eToken && eToken?.length > 0) window.localStorage.setItem("JWT_TOKEN", eToken);
};

/**
 * it fetch, parse, check from localstorage the user's datas
 * @returns {FunctionJob<User>} the user's datas
 */
export const GetLSUser = (): FunctionJob<User> => {
	const rawToken = window.localStorage.getItem("user");
	if (!rawToken || IsEmptyString(rawToken)) return { success: false };

	try {
		return { success: true, data: JSON.parse(rawToken) };
	} catch (err) {
		return { success: false };
	}
};

/**
 * it fetch, parse, check from localstorage the user's subs
 * @returns {FunctionJob<UserSubs>} the user's subs
 */
export const GetUserSubs = (): FunctionJob<UserSubs> => {
	const rawSubs = window.localStorage.getItem("subs");
	if (!rawSubs || IsEmptyString(rawSubs) || !IsValidJSON(rawSubs)) return { success: false };

	const subs: UserSubs = JSON.parse(rawSubs) ?? undefined;
	if (typeof subs !== "object" || !Object.hasOwn(subs, "teddit") || !Object.hasOwn(subs, "nitter"))
		return { success: false };

	return { success: true, data: subs };
};

/**
 * it fetch, parse, check from localstorage the user's nitter lists
 * @returns {FunctionJob<NitterLists[]>} the user's nitter lists
 */
export const GetUserLists = (): FunctionJob<NitterLists[]> => {
	const rawList = window.localStorage.getItem("lists");
	if (!rawList || IsEmptyString(rawList) || !IsValidJSON(rawList)) return { success: false };

	const lists: NitterLists[] = JSON.parse(rawList) ?? undefined;
	if (typeof lists !== "object" || lists.length <= 0) return { success: false };

	return { success: true, data: lists };
};
