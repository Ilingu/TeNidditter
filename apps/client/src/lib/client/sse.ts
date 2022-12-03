import { PUBLIC_API_URL } from "$env/static/public";

import { IsEmptyString, IsValidJSON } from "$lib/shared/utils";
import type { SSEReturns, SSERoutes } from "./types/sseScheme";

enum SEEState {
	NOT_INITIALIZED,
	CONNECTING,
	CONNECTED,
	CLOSED
}

/**
 * @classdesc Tweniditter backend Server-Sent-Events client library (fully typed)
 * @class
 */
export default class SSEHandler<T extends SSERoutes> {
	/**
	 * @private
	 * @readonly
	 * @property
	 */
	private readonly route: T;
	/**
	 * @private
	 * @property
	 */
	private state = SEEState.NOT_INITIALIZED;
	/**
	 * @private
	 * @property
	 */
	private sseConn?: EventSource;

	/**
	 * open connection timeout, if exceeded it means that the connection cannot be opened
	 *
	 * default: `20_000ms` <=> `20s`
	 * @public
	 * @property
	 */
	public connectionTimeout = 20_000;

	/**
	 * @private
	 * @property
	 */
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	private msgListeners: Record<string, (msg: any) => void> = {};
	/**
	 * @private
	 * @property
	 */
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	private closeListeners: Record<string, () => void> = {};
	/**
	 * @private
	 * @property
	 */
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	private errListeners: Record<string, (msg: any) => void> = {};

	constructor(route: T) {
		this.route = route;
	}

	/**
	 * @public
	 * @property
	 */
	public get readyState(): SEEState {
		return this.state;
	}

	/**
	 * it opens a SSE connection to the route
	 * @public
	 * @method
	 * @returns {Promise<boolean>} if `true`, it means that the connection opened successfully
	 */
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

	/**
	 * add a event listener to this connection
	 * @public
	 * @method
	 * @param {U} ev - event to listen to `"message" | "error" | "close"`
	 * @param {(msg?: any) => void} cb - callback function to call when the vent is triggered
	 */
	public on<U extends "message" | "error" | "close">(
		ev: U,
		cb: (msg?: U extends "message" ? SSEReturns<T> : U extends "error" ? Event : never) => void
	) {
		// let Id = RandomChars(8);
		const Id = `${Date.now() - Math.random() * 100}-${ev}`;
		// if (IsEmptyString(Id)) Id = `${Date.now() - Math.random() * 100}`;

		if (ev === "message") this.msgListeners[Id] = cb;
		if (ev === "error") this.errListeners[Id] = cb;
		if (ev === "close") this.closeListeners[Id] = cb;
	}
	/**
	 * remove a already added event listener from this connection
	 * @public
	 * @method
	 * @param {U} ev - event to unlisten `"message" | "error" | "close"`
	 * @param {string} Id - Id indentifying the callback function to remove
	 */
	public off(ev: "message" | "error" | "close", Id: string) {
		if (ev === "message") delete this.msgListeners[Id];
		if (ev === "error") delete this.errListeners[Id];
		if (ev === "close") delete this.closeListeners[Id];
	}

	/**
	 * it closes the SSE connection, trigger the "close" event and reset everything back to the class initialization (all the listeners are removed)
	 * @public
	 * @method
	 */
	public close() {
		this.state === SEEState.CLOSED;
		this.sseConn?.close();
		Object.values(this.closeListeners).forEach((cb) => cb());

		this.sseConn = undefined;
		this.msgListeners = {};
		this.errListeners = {};
		this.closeListeners = {};
	}

	/**
	 * helper method that wait until the conn is opened or that an error occurred or that the timeout timed out
	 * @private
	 * @method
	 * @param {EventSource} source - the SSE connection
	 * @returns {Promise<boolean>} if `true`, it means that the connection is opened
	 */
	private waitConnectionOpen(source: EventSource): Promise<boolean> {
		return new Promise<boolean>((resolve) => {
			setTimeout(() => resolve(false), this.connectionTimeout);
			source.onopen = () => resolve(true);
			source.onerror = () => resolve(false);
		});
	}

	/**
	 * it listen, handle and parse in backgroud all the events (message/error) coming from the server and triggers the cb event listeners
	 * @private
	 * @method
	 */
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
