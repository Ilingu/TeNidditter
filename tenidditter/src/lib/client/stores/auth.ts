/* eslint-disable @typescript-eslint/no-empty-function */
import { ListenToUserChange, LogOut } from "$lib/client/services/auth";
import { SetJWT } from "$lib/client/services/localstorage";
import {
	DeleteUserAccount,
	RegenerateUserRecoveryCodes,
	ToggleNitterSubs,
	ToggleTedditSubs
} from "$lib/client/services/user";
import { IsEmptyString } from "$lib/shared/utils";
import { writable } from "svelte/store";
import type { AuthStoreShape, User, UserSubs } from "../types/auth";
import type { NitterLists } from "../types/nitter";

const JwtTokenRegExp = /^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$/g;

/* STORE */
export const defaultAuthStore: AuthStoreShape = {
	loggedIn: false,
	user: {
		id: NaN,
		username: "NaU",
		exp: 0,
		action: {
			deleteAccount: async () => {},
			regenerateUserRecoveryCodes: async () => ({ success: false }),
			toggleTedditSubs: async () => ({ success: false }),
			toggleNitterSubs: async () => ({ success: false }),
			logout: () => {}
		}
	}
};

/**
 * Global App Auth store, contains all datas about the currently logged in user
 * @constant
 */
const AuthStore = writable<AuthStoreShape>(defaultAuthStore);

/**
 * Get an AuthStore data snapshot
 * @returns {Promise<AuthStoreShape>} auth datas
 */
export const GetUserSession = (): Promise<AuthStoreShape> => {
	return new Promise((res) => {
		const UnSub = AuthStore.subscribe((value) => {
			UnSub();
			res(value);
		});
	});
};

/**
 * Overwrite the store to instanciate a new user session & sync with localstorage; **only if the session is valid** (no invalid jwtTk (wrong shape or expired) nor invalid inputs)
 * @param {User} user - `User` object
 * @param {string} JwtToken - user's jwtToken
 * @param {UserSubs} Subs - user's subs
 * @param {NitterLists[]} Lists = user's lists
 */
export const SetUserSession = (
	user: User,
	JwtToken: string,
	Subs: UserSubs,
	Lists: NitterLists[]
) => {
	if (IsEmptyString(JwtToken) || !JwtTokenRegExp.test(JwtToken)) return;
	if (typeof user !== "object" || IsEmptyString(user?.username) || typeof user?.exp !== "number")
		return;

	const exp = user.exp * 1000; // convert s to ms
	if (exp <= Date.now()) return LogOut();

	SetJWT(JwtToken);
	window.localStorage.setItem("user", JSON.stringify(user));
	localStorage.setItem("subs", JSON.stringify(Subs));
	localStorage.setItem("lists", JSON.stringify(Lists));

	user.action = {
		deleteAccount: () => DeleteUserAccount(JwtToken),
		regenerateUserRecoveryCodes: () => RegenerateUserRecoveryCodes(JwtToken),
		toggleTedditSubs: ToggleTedditSubs,
		toggleNitterSubs: ToggleNitterSubs,
		logout: () => LogOut(true, JwtToken)
	};

	AuthStore.set({ loggedIn: true, user, JwtToken, Subs, Lists });
	ListenToUserChange(JwtToken);
};

/**
 * update the store's `subs` field with the input data and sync localstorage
 * @param {UserSubs} Subs - new user's subs
 */
export const UpdateUserSubs = (Subs: UserSubs) => {
	localStorage.setItem("subs", JSON.stringify(Subs));
	AuthStore.update((s) => ({ ...s, Subs }));
};
/**
 * update the store's `lists` field with the input data and sync localstorage
 * @param {NitterLists[]} Lists - new user's subs
 */
export const UpdateUserList = (Lists: NitterLists[]) => {
	localStorage.setItem("lists", JSON.stringify(Lists));
	AuthStore.update((s) => ({ ...s, Lists }));
};

export default AuthStore;
