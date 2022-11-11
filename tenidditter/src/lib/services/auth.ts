import api from "$lib/api";
import AuthStore, { SetUserSession, UpdateUserSubs, type User } from "$lib/stores/auth";
import type { FunctionJob, UserSubs } from "$lib/types/interfaces";
import { HandleWsConn, NewWsConn } from "$lib/ws";
import { GetJWT, GetLSUser, GetUserSubs } from "./localstorage";
import { PUBLIC_API_URL } from "$env/static/public";
import { IsValidJSON, MakeBearerToken } from "$lib/utils";

/* AUTH FUNC */
export const SignIn = async (JwtToken: string, Subs: UserSubs) => {
	// check if jwt still valid
	const { success, data: user } = await GetUserInfo(JwtToken);
	if (!success || !user) return LogOut();

	SetUserSession(user, JwtToken, Subs);
};

export const AutoLogin = async () => {
	const { success: fetchJwt, data: JwtToken } = await GetJWT();
	if (!fetchJwt || !JwtToken) return LogOut();

	const { success: FetchUserSubsSucceed, data: Subs } = GetUserSubs();
	if (!FetchUserSubsSucceed || typeof Subs !== "object" || !Object.hasOwn(Subs, "teddit"))
		return LogOut();

	const { success, data: user } = GetLSUser();
	if (success && typeof user === "object") SetUserSession(user, JwtToken, Subs);

	// check if jwt still valid
	const { success: validJwt } = await GetUserInfo(JwtToken);
	if (!validJwt) return LogOut();
};

export const LogOut = async (serverLogout = false, JwtToken?: string) => {
	window.localStorage.clear();
	AuthStore.set({ loggedIn: false });

	if (serverLogout)
		return await api.delete("/auth/logout", {
			credentials: true,
			query: JwtToken ? { token: JwtToken } : undefined
		});
	fetch("/api/logout", { method: "delete" });
	// "executionContexts" should handle the reload part
};

const DispatchUserChange = async (data: unknown) => {
	console.log("Mew WS Msg: ", { data });
	const BinaryData = data as Blob;

	try {
		const BinaryUtf8 = await BinaryData.text();
		if (!IsValidJSON(BinaryUtf8)) return console.warn("WS Msg Failed");

		const NewSubs = JSON.parse(BinaryUtf8);
		if (
			typeof NewSubs !== "object" ||
			!Object.hasOwn(NewSubs, "teddit") ||
			!Object.hasOwn(NewSubs, "nitter")
		)
			return console.warn("WS Msg Failed");
		UpdateUserSubs(NewSubs);
	} catch (err) {
		console.warn("WS Msg Failed");
	}
};
export const ListenToUserChange = async (JwtToken: string) => {
	const URL = `${PUBLIC_API_URL.replace("http", "ws")}/auth/userChanged?token=${JwtToken}`;

	const { success: connected, data: socket } = await NewWsConn(URL);
	if (!connected || !socket) return;
	console.log("Listening to userChange...");

	HandleWsConn(socket, { onMessage: DispatchUserChange });
};

export const GetUserInfo = async (JwtToken: string): Promise<FunctionJob<User>> => {
	const { success: LoginSuccess, data: user } = await api.get("/tedinitter/userInfo", {
		headers: MakeBearerToken(JwtToken)
	});

	if (!LoginSuccess || !user) return { success: false };
	return { success: true, data: user };
};
