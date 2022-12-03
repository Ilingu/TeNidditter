import type { FunctionJob } from "$lib/shared/types/globals";
import type { NitterLists } from "./nitter";

export interface UserSubs {
	teddit?: string[];
	nitter?: string[];
}

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
export interface AuthStoreShape {
	loggedIn: boolean;
	user: User;
	JwtToken?: string;
	Subs?: UserSubs;
	Lists?: NitterLists[];
}
