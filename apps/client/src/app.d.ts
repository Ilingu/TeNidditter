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

/**
 * Simple AES Encryption algorythm
 * @param {string} key - 32 bytes key
 * @param {string} textToEnc - text to encrypt
 * @returns {string} the encrypted text
 */
declare function EncryptAES(key: string, textToEnc: string): string;
/**
 * Simple AES Decryption algorythm
 * @param {string} key - 32 bytes key used for encryption
 * @param {string} textToDec - previously encrypted string
 * @returns {string} the original/decrypted string
 */
declare function DecryptAES(key: string, textToDec: string): string;
/**
 * Simple sha256 hash function
 * @param {string} toHash - (if you've  multiple value, concatenate them into one string)
 * @returns {string} the hash string of the input string
 */
declare function Hash(toHash: string): string;
/**
 * It generate a string of "length" random characters, ideal to generate ids,uuids,codes...
 * @param {number} length
 * @returns {string} the posts (if afterId set it's only the post after the poist with the id=afterId)
 */
declare function RandomChars(length: number): string;

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
