import type { FunctionJob } from "$lib/types/interfaces";
import api from "$lib/api";
import { IsEmptyString, pushAlert } from "$lib/utils";
import { writable } from "svelte/store";
import { DecryptDatas, EncryptDatas } from "$lib/encryption";

export interface User {
	username: string;
	exp: number;
	admin: false;
}
interface AuthStoreShape {
	loggedIn: boolean;
	user?: User;
	JwtToken?: string;
}

const JwtTokenRegExp = /^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$/g;

/* AUTH FUNC */
export const AutoLogin = async (JwtToken?: string) => {
	if (!JwtToken) {
		const { success, data: jwt } = await GetJWT();
		if (!success || !jwt) return LogOut();
		JwtToken = jwt;
	}
	{
		const { success, data: user } = GetLSUser();
		if (success && typeof user === "object") SetUserSession(user, JwtToken);
	}

	const { success, data: user } = await GetUserInfo(JwtToken);
	if (!success || !user) return LogOut();

	SetUserSession(user, JwtToken);
};

export const GetUserInfo = async (JwtToken: string): Promise<FunctionJob<User>> => {
	const { success: LoginSuccess, data: user } = await api.get<User>({
		uri: "/tedinitter/userInfo",
		headers: { Authorization: "Bearer " + JwtToken }
	});

	if (!LoginSuccess || !user) return { success: false };
	return { success: true, data: user };
};

const LogOut = () => {
	pushAlert("Failed to login", "error", 3600);

	window.localStorage.removeItem("JWT_TOKEN");
	window.localStorage.removeItem("user");
	AuthStore.set({ loggedIn: false });
};

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

/* STORE */
const AuthStore = writable<AuthStoreShape>({ loggedIn: false });
export const GetUserSession = (): Promise<AuthStoreShape> =>
	new Promise((res) => {
		const UnSub = AuthStore.subscribe((value) => {
			UnSub();
			res(value);
		});
	});

export const SetUserSession = (user: User, JwtToken: string) => {
	if (!JwtTokenRegExp.test(JwtToken)) return;
	if (typeof user !== "object" || IsEmptyString(user?.username) || typeof user?.exp !== "number")
		return;

	const exp = user.exp * 1000; // convert s to ms
	if (exp <= Date.now()) return LogOut();

	SetJWT(JwtToken);
	window.localStorage.setItem("user", JSON.stringify(user));

	AuthStore.set({ loggedIn: true, user, JwtToken });
};

export default AuthStore;
