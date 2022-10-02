import type { FunctionJob } from "./types/interfaces";

let wasm: WebAssembly.Instance;

const WASM_URL = "/wasm/encryption.wasm";
export const InitWasm = async (): Promise<FunctionJob> => {
	try {
		const go = new Go();

		if ("instantiateStreaming" in WebAssembly) {
			const obj = await WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject);
			wasm = obj.instance;
			// go.run(wasm); // --> run main()
			console.log(wasm.exports);

			return { success: true };
		}

		// Classic way if instantiateStreaming not in WebAssembly
		const resp = await fetch(WASM_URL);
		const bytes = await resp.arrayBuffer();
		await WebAssembly.instantiate(bytes, go.importObject).then((obj) => {
			wasm = obj.instance;
		});
		return { success: true };
	} catch (err) {
		console.error(err);
		return { success: false };
	}
};

export const EncryptAES = (str: string): FunctionJob<string> => {
	if (!wasm || !wasm?.exports?.EncryptAES) return { success: false };

	const encryptedStr = wasm.exports.EncryptAES(
		"thisis32bitlongpassphraseimusing",
		"This is a secret"
	);

	console.log(encryptedStr);

	return { success: true, data: encryptedStr };
};
