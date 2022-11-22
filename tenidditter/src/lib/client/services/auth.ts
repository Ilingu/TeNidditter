/* eslint-disable @typescript-eslint/no-non-null-assertion */
import api from "$lib/shared/api";
import AuthStore, {
	defaultAuthStore,
	SetUserSession,
	UpdateUserList,
	UpdateUserSubs
} from "$lib/client/stores/auth";
import { GetJWT, GetLSUser, GetUserLists, GetUserSubs } from "./localstorage";
import { PUBLIC_API_URL } from "$env/static/public";
import { FormatUsername, IsEmptyString, IsValidJSON, MakeBearerToken } from "$lib/shared/utils";
import { pushAlert } from "../ClientUtils";
import { GetZxcvbn } from "../zxcvbn";
import type { NitterLists } from "../types/nitter";
import type { User, UserSubs } from "../types/auth";
import { HandleWsConn, NewWsConn } from "../ws";
import type { FunctionJob } from "$lib/shared/types/globals";

/* AUTH FUNC */
export const Authenticate = async (
	Username: string,
	Password: string,
	AuthMethod: "login" | "register",
	recoveryCodeCb?: (codes: string[]) => void
) => {
	// Username check
	const username = FormatUsername(Username.trim());
	const nameLen = username.length;
	if (nameLen < 3 || nameLen > 15) return pushAlert("Username invalid!", "warning");

	// Check Password
	const zxcvbn = await GetZxcvbn();
	const passwordStrenght = zxcvbn(Password);
	if (passwordStrenght.score < 3) return pushAlert("Password is too weak! Like you", "warning");

	if (AuthMethod === "register") {
		const { success, data: IsAvailable } = await api.get("/auth/available", {
			query: { username }
		});

		if (!success || !IsAvailable)
			return pushAlert("This Username is already taken.", "warning", 6000);
	}

	const {
		success: AuthSuccess,
		data: JwtToken,
		headers
	} = await api.post("/auth/", {
		credentials: true,
		body: {
			username,
			password: Password,
			method: AuthMethod
		}
	});

	if (!AuthSuccess) return pushAlert("Failed to auth", "error");

	if (AuthMethod === "register") return AfterRegister(headers, recoveryCodeCb);
	if (AuthMethod === "login" && !IsEmptyString(JwtToken))
		return AfterLogin(JwtToken as string, headers);

	pushAlert("Invalid login", "error");
};

const AfterRegister = async (headers?: Headers, recoveryCodeCb?: (codes: string[]) => void) => {
	const codesHeader = headers?.get("RecoveryCodes") ?? "";

	if (codesHeader && !IsEmptyString(codesHeader) && IsValidJSON(codesHeader))
		recoveryCodeCb && recoveryCodeCb(JSON.parse(codesHeader));

	pushAlert("Successfully registered, you can now login", "success", 6000);
};

const AfterLogin = async (JwtToken: string, headers?: Headers) => {
	const tedditSubsHeader = headers?.get("TedditSubs") ?? "",
		nitterSubsHeader = headers?.get("NitterSubs") ?? "",
		nitterListsHeader = headers?.get("NitterLists") ?? ""; // retrieve user subs

	let tedditSubs: string[] = [],
		nitterSubs: string[] = [],
		nitterLists: NitterLists[] = [];

	if (tedditSubsHeader && !IsEmptyString(tedditSubsHeader) && IsValidJSON(tedditSubsHeader))
		tedditSubs = JSON.parse(tedditSubsHeader);
	if (nitterSubsHeader && !IsEmptyString(nitterSubsHeader) && IsValidJSON(nitterSubsHeader))
		nitterSubs = JSON.parse(nitterSubsHeader);
	if (nitterListsHeader && !IsEmptyString(nitterListsHeader) && IsValidJSON(nitterListsHeader))
		nitterLists = JSON.parse(nitterListsHeader);

	const subs = { teddit: tedditSubs, nitter: nitterSubs };
	await FirstLogin(JwtToken, subs, nitterLists); // validate user jwt

	pushAlert("Successfully logged in!", "success");
};

export const FirstLogin = async (JwtToken: string, Subs: UserSubs, Lists: NitterLists[]) => {
	// check if jwt still valid
	const { success, data: user } = await GetUserInfo(JwtToken);
	if (!success || !user) return LogOut();

	SetUserSession(user, JwtToken, Subs, Lists);
};

export const AutoLogin = async () => {
	const { success: fetchJwt, data: JwtToken } = await GetJWT();
	if (!fetchJwt || !JwtToken) return LogOut();

	const { data: Subs } = GetUserSubs();
	const { data: Lists } = GetUserLists();

	const { success, data: user } = GetLSUser();
	if (success && typeof user === "object")
		SetUserSession(user, JwtToken, Subs ?? { nitter: [], teddit: [] }, Lists ?? []);

	// check if jwt still valid
	const { success: validJwt } = await GetUserInfo(JwtToken);
	if (!validJwt) return LogOut();
};

export const LogOut = async (serverLogout = false, JwtToken?: string) => {
	window.localStorage.clear();
	AuthStore.set(defaultAuthStore);

	if (serverLogout)
		return await api.delete("/auth/logout", {
			credentials: true,
			query: JwtToken ? { token: JwtToken } : undefined
		});
	fetch("/api/logout", { method: "delete" });
	// "executionContexts" should handle the reload part
};

const DispatchSessionChange = async (data: unknown) => {
	console.log("Mew WS Msg: ", { data });
	const BinaryData = data as Blob;

	try {
		const BinaryUtf8 = await BinaryData.text();
		if (!IsValidJSON(BinaryUtf8)) return console.warn("WS Msg Failed");

		const NewDatas = JSON.parse(BinaryUtf8);
		if (typeof NewDatas !== "object") return console.warn("WS Msg Failed");

		const IsUserChange = Object.hasOwn(NewDatas, "teddit") && Object.hasOwn(NewDatas, "nitter");
		if (IsUserChange) return UpdateUserSubs(NewDatas);

		const IsListsChange = NewDatas.length > 0 && Object.hasOwn(NewDatas[0], "title");
		if (IsListsChange) return UpdateUserList(NewDatas);
	} catch (err) {
		console.warn("WS Msg Failed");
	}
};
export const ListenToUserChange = async (JwtToken: string) => {
	const URL = `${PUBLIC_API_URL.replace("http", "ws")}/auth/userChanged?token=${JwtToken}`;

	const { success: connected, data: socket } = await NewWsConn(URL);
	if (!connected || !socket) return;
	console.log("Listening to userChange...");

	HandleWsConn(socket, { onMessage: DispatchSessionChange });
};

export const GetUserInfo = async (JwtToken: string): Promise<FunctionJob<User>> => {
	const { success: LoginSuccess, data: user } = await api.get("/tedinitter/userInfo", {
		headers: MakeBearerToken(JwtToken)
	});

	if (!LoginSuccess || !user) return { success: false };
	return { success: true, data: user };
};
