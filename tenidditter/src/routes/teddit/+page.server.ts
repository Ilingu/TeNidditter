import type { FunctionJob } from "$lib/types/interfaces";
import { IsEmptyString } from "$lib/utils";

// TODO Better encrypt/decrypt system --> WASM
// + api client
export const load: import("./$types").PageServerLoad = ({ cookies }): FunctionJob => {
	const eToken = cookies.get("JWT_TOKEN");
	if (!eToken || IsEmptyString(eToken)) return { success: false };

	return { success: true };
};
