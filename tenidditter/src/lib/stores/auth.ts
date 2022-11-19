/* eslint-disable @typescript-eslint/no-empty-function */
import { ListenToUserChange, LogOut } from "$lib/services/auth";
import { SetJWT } from "$lib/services/localstorage";
import {
	DeleteUserAccount,
	RegenerateUserRecoveryCodes,
	ToggleNitterSubs,
	ToggleTedditSubs
} from "$lib/services/user";
import type { FunctionJob, NitterLists, UserSubs } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { writable } from "svelte/store";

export interface User {
	username: string;
	exp: number;
	id: number;
	action: {
		deleteAccount: () => Promise<void>;
		regenerateUserRecoveryCodes: () => Promise<FunctionJob<string[]>>;
		toggleTedditSubs: (subteddit: string, isSub: boolean, JwtToken: string) => Promise<FunctionJob>;
		toggleNitterSubs: (nittos: string, isSub: boolean, JwtToken: string) => Promise<FunctionJob>;

		logout: () => void;
	};
}
interface AuthStoreShape {
	loggedIn: boolean;
	user: User;
	JwtToken?: string;
	Subs?: UserSubs;
	Lists?: NitterLists[];
}

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
