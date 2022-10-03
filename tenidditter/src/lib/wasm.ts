import type { FunctionJob } from "./types/interfaces";
import { IsEmptyString } from "./utils";

let wasmBinary: WebAssembly.Instance;

const WASM_URL = "/wasm/encryption.wasm";
export const InitWasm = async (): Promise<FunctionJob> => {
	try {
		const go = new Go();

		if ("instantiateStreaming" in WebAssembly) {
			const obj = await WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject);
			wasmBinary = obj.instance;
			go.run(wasmBinary); // --> run main()

			return { success: true };
		}

		// Classic way if instantiateStreaming not in WebAssembly
		const resp = await fetch(WASM_URL);
		const bytes = await resp.arrayBuffer();
		await WebAssembly.instantiate(bytes, go.importObject).then((obj) => {
			wasmBinary = obj.instance;
			go.run(wasmBinary); // --> run main()
		});
		return { success: true };
	} catch (err) {
		console.error(err);
		return { success: false };
	}
};

export const EncryptDatas = (str: string): FunctionJob<string> => {
	if (IsEmptyString(str)) return { success: false };

	if (!wasmBinary || !EncryptAES) return { success: false };
	const encryptedStr = EncryptAES("thisis32bitlongpassphraseimusing", str);

	if (IsEmptyString(encryptedStr)) return { success: false };
	return { success: true, data: encryptedStr };
};
export const DecryptDatas = (str: string): FunctionJob<string> => {
	if (IsEmptyString(str)) return { success: false };

	if (!wasmBinary || !DecryptAES) return { success: false };
	const decryptedStr = DecryptAES("thisis32bitlongpassphraseimusing", str);

	if (IsEmptyString(decryptedStr)) return { success: false };
	return { success: true, data: decryptedStr };
};
