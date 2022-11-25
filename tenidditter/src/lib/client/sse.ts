import { PUBLIC_API_URL } from "$env/static/public";

import { IsEmptyString, IsValidJSON } from "$lib/shared/utils";
import type { SSEReturns, SSERoutes } from "./types/sseScheme";

enum SEEState {
	NOT_INITIALIZED,
	CONNECTING,
	CONNECTED,
	CLOSED
}

export default class SSEHandler<T extends SSERoutes> {
	private readonly route: T;
	private state = SEEState.NOT_INITIALIZED;
	private sseConn?: EventSource;

	public connectionTimeout = 20_000;

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	private msgListeners: Record<string, (msg: any) => void> = {};
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	private closeListeners: Record<string, () => void> = {};
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	private errListeners: Record<string, (msg: any) => void> = {};

	constructor(route: T) {
		this.route = route;
	}

	get readyState(): SEEState {
		return this.state;
	}

	public async connect(): Promise<boolean> {
		if (this.state === SEEState.CONNECTED || this.state === SEEState.CONNECTING) return false;
		this.state = SEEState.CONNECTING;

		const source = new EventSource(`${PUBLIC_API_URL}${this.route}`);
		const connected = await this.waitConnectionOpen(source);
		if (!connected) return false;

		this.state = SEEState.CONNECTED;
		this.sseConn = source;
		this.listenSSE();

		return true;
	}

	public on<U extends "message" | "error" | "close">(
		ev: U,
		cb: (msg?: U extends "message" ? SSEReturns<T> : U extends "error" ? Event : never) => void
	) {
		// let Id = RandomChars(8);
		const Id = `${Date.now() - Math.random() * 100}-${ev}`;
		// if (IsEmptyString(Id)) Id = `${Date.now() - Math.random() * 100}`;

		if (ev === "message") {
			if (this.state !== SEEState.CONNECTED || !this.sseConn) return;
			this.msgListeners[Id] = cb;
		}
		if (ev === "error") this.errListeners[Id] = cb;
		if (ev === "close") this.closeListeners[Id] = cb;
	}
	public off(ev: "message" | "error" | "close", Id: string) {
		if (ev === "message") delete this.msgListeners[Id];
		if (ev === "error") delete this.errListeners[Id];
		if (ev === "close") delete this.closeListeners[Id];
	}

	public close() {
		this.state === SEEState.CLOSED;
		this.sseConn?.close();
		Object.values(this.closeListeners).forEach((cb) => cb());

		this.sseConn = undefined;
		this.msgListeners = {};
		this.errListeners = {};
		this.closeListeners = {};
	}

	private waitConnectionOpen(source: EventSource): Promise<boolean> {
		return new Promise<boolean>((resolve) => {
			setTimeout(() => resolve(false), this.connectionTimeout);
			source.onopen = () => resolve(true);
			source.onerror = () => resolve(false);
		});
	}

	private listenSSE() {
		if (!this.sseConn) return;
		this.sseConn.onmessage = (sourceEv) => {
			if (!sourceEv?.data || IsEmptyString(sourceEv.data)) return;
			if (sourceEv.data === "CLOSING") return this.close();

			if (!IsValidJSON(sourceEv.data)) return;
			const ParsedData = JSON.parse(sourceEv.data);
			if (!ParsedData) return;

			for (const cb of Object.values(this.msgListeners)) {
				cb(ParsedData);
			}
		};

		this.sseConn.onerror = (ev) => {
			for (const cb of Object.values(this.errListeners)) {
				cb(ev);
			}
		};
	}
}
