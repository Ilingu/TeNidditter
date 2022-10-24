import type { FunctionJob } from "../../types/interfaces";
import { IsEmptyString } from "../../utils";
import { WasmInitiate } from "./wasm";

import { PUBLIC_ENCRYPT_KEY } from "$env/static/public";

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
