import type { FunctionJob } from "$lib/shared/types/globals";

let wasmBinary: WebAssembly.Instance;
export let WasmInitiate = false;

const WASM_URL = "/wasm/wasmbin.wasm";
export const InitWasm = async (): Promise<FunctionJob> => {
	try {
		const go = new Go();

		if ("instantiateStreaming" in WebAssembly) {
			const obj = await WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject);
			wasmBinary = obj.instance;
			go.run(wasmBinary); // --> run main()

			WasmInitiate = true;
			return { success: true };
		}

		// Classic way if instantiateStreaming not in WebAssembly
		const resp = await fetch(WASM_URL);
		const bytes = await resp.arrayBuffer();
		await WebAssembly.instantiate(bytes, go.importObject).then((obj) => {
			wasmBinary = obj.instance;
			go.run(wasmBinary); // --> run main()
		});

		WasmInitiate = true;
		return { success: true };
	} catch (err) {
		console.error(err);
		return { success: false };
	}
};
