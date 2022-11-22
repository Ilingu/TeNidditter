import { WasmInitiate } from "./wasm";

import { PUBLIC_ENCRYPT_KEY } from "$env/static/public";
import { IsEmptyString } from "$lib/shared/utils";
import type { FunctionJob } from "$lib/shared/types/globals";

export const HashDatas = (str: string): string => {
	if (IsEmptyString(str)) return "";

	if (!WasmInitiate || !Hash) return "";
	const hashed = Hash(str);

	if (IsEmptyString(hashed)) return "";
	return hashed;
};
export const EncryptDatas = (str: string): FunctionJob<string> => {
	if (IsEmptyString(str)) return { success: false };

	if (!WasmInitiate || !EncryptAES) return { success: false };
	const encryptedStr = EncryptAES(str, PUBLIC_ENCRYPT_KEY);

	if (IsEmptyString(encryptedStr)) return { success: false };
	return { success: true, data: encryptedStr };
};
export const DecryptDatas = (str: string): FunctionJob<string> => {
	if (IsEmptyString(str)) return { success: false };

	if (!WasmInitiate || !DecryptAES) return { success: false };
	const decryptedStr = DecryptAES(str, PUBLIC_ENCRYPT_KEY);

	if (IsEmptyString(decryptedStr)) return { success: false };
	return { success: true, data: decryptedStr };
};
