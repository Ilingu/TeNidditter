// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
// and what to do when importing types
declare namespace App {
	// interface Locals {}
	// interface Platform {}
	// interface PrivateEnv {}
	interface PublicEnv {
		PUBLIC_API_URL: string;
		PUBLIC_ENCRYPT_KEY: string;
	}
	// interface Session {}
	// interface Stuff {}
}

declare class Go {
	_callbackTimeouts: Record<string, unknown>;
	_nextCallbackTimeoutID: number;
	importObject: WebAssembly.Imports | undefined;
}

declare function EncryptAES(key: string, textToEnc: string): string;
declare function DecryptAES(key: string, textToDec: string): string;
declare function Hash(toHash: string): string;

interface CustomEventMap {
	alertEvent: CustomEvent<AlertShape>;
}
declare global {
	interface Document {
		addEventListener<K extends keyof CustomEventMap>(
			type: K,
			listener: (this: Document, ev: CustomEventMap[K]) => void
		): void;
	}
	// interface Window {
	// 	opera: any;
	// 	appVersion: () => "Web" | "PWA";
	// }
	// interface Navigator {
	// 	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	// 	userAgentData: any;
	// }
}
