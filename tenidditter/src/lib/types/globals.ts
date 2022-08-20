import type { AlertShape } from "./interfaces";

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
}
