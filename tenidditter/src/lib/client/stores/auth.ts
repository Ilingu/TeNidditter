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
const AuthStore = writable<AuthStoreShape>(defaultAuthStore);
export const GetUserSession = (): Promise<AuthStoreShape> =>
	new Promise((res) => {
		const UnSub = AuthStore.subscribe((value) => {
			UnSub();
			res(value);
		});
	});

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

export const UpdateUserSubs = (Subs: UserSubs) => {
	localStorage.setItem("subs", JSON.stringify(Subs));
	AuthStore.update((s) => ({ ...s, Subs }));
};
export const UpdateUserList = (Lists: NitterLists[]) => {
	localStorage.setItem("lists", JSON.stringify(Lists));
	AuthStore.update((s) => ({ ...s, Lists }));
};

export default AuthStore;
