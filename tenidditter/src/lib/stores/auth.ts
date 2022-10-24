import { ListenToUserChange, LogOut } from "$lib/services/auth";
import { SetJWT } from "$lib/services/localstorage";
import type { UserSubs } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";
import { writable } from "svelte/store";

export interface User {
	username: string;
	exp: number;
	id: number;
}
interface AuthStoreShape {
	loggedIn: boolean;
	user?: User;
	JwtToken?: string;
	Subs?: UserSubs;
}

const JwtTokenRegExp = /^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$/g;

/* STORE */
const AuthStore = writable<AuthStoreShape>({ loggedIn: false });
export const GetUserSession = (): Promise<AuthStoreShape> =>
	new Promise((res) => {
		const UnSub = AuthStore.subscribe((value) => {
			UnSub();
			res(value);
		});
	});

export const SetUserSession = (user: User, JwtToken: string, Subs: UserSubs) => {
	if (!JwtTokenRegExp.test(JwtToken)) return;
	if (typeof user !== "object" || IsEmptyString(user?.username) || typeof user?.exp !== "number")
		return;

	const exp = user.exp * 1000; // convert s to ms
	if (exp <= Date.now()) return LogOut();

	SetJWT(JwtToken);
	window.localStorage.setItem("user", JSON.stringify(user));
	localStorage.setItem("subs", JSON.stringify(Subs));

	AuthStore.set({ loggedIn: true, user, JwtToken, Subs });
	ListenToUserChange(JwtToken);
};

export const UpdateUserSubs = (Subs: UserSubs) => {
	localStorage.setItem("subs", JSON.stringify(Subs));
	AuthStore.update((s) => ({ ...s, Subs }));
};

export default AuthStore;
