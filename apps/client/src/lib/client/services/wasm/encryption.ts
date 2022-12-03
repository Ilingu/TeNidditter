import { WasmInitiate } from "./wasm";

import { PUBLIC_ENCRYPT_KEY } from "$env/static/public";
import { IsEmptyString } from "$lib/shared/utils";
import type { FunctionJob } from "$lib/shared/types/globals";

/**
 * Simple sha256 hash function
 * @param {string} str - strings to hash
 * @returns {string} the hash string of the input string, **if `""` is returned, the operation has failed**
 */
export const HashDatas = (...str: string[]): string => {
	const toHash = str.join("");
	if (IsEmptyString(toHash)) return "";

	if (!WasmInitiate || !Hash) return "";
	const hashed = Hash(toHash);

	if (IsEmptyString(hashed)) return "";
	return hashed;
};
/**
 * Simple AES Encryption algorythm
 * @param {string} str - text to encrypt
 * @returns {FunctionJob<string>} success state and the encrypted text *(if success=true)*
 */
export const EncryptDatas = (str: string): FunctionJob<string> => {
	if (IsEmptyString(str)) return { success: false };

	if (!WasmInitiate || !EncryptAES) return { success: false };
	const encryptedStr = EncryptAES(str, PUBLIC_ENCRYPT_KEY);

	if (IsEmptyString(encryptedStr)) return { success: false };
	return { success: true, data: encryptedStr };
};
/**
 * Simple AES Decryption algorythm
 * @param {string} str - previously encrypted string
 * @returns {FunctionJob<string>} success state and the original/decrypted string *(if success=true)*
 */
export const DecryptDatas = (str: string): FunctionJob<string> => {
	if (IsEmptyString(str)) return { success: false };

	if (!WasmInitiate || !DecryptAES) return { success: false };
	const decryptedStr = DecryptAES(str, PUBLIC_ENCRYPT_KEY);

	if (IsEmptyString(decryptedStr)) return { success: false };
	return { success: true, data: decryptedStr };
};
