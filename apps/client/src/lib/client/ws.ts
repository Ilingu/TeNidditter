import type { FunctionJob } from "$lib/shared/types/globals";
import { isValidUrl } from "$lib/shared/utils";
import { LogOut } from "./services/auth";

type Args = {
	onMessage: (data: unknown) => void;
	onError?: () => void;
	onClose?: () => void;
};
/**
 * It created and connects to an websocket endpoint (if you want my thought: I should rewrite this to make a class client with typed routes like sse, but I'm too lazy)
 * @param {string} wsEndpoint
 * @returns {Promise<FunctionJob<WebSocket>>} the websocket connection (`window.WebSocket`)
 */
export const NewWsConn = async (wsEndpoint: string): Promise<FunctionJob<WebSocket>> => {
	if (!window.WebSocket) return { success: false };
	if (!isValidUrl(wsEndpoint)) return { success: false };

	const { protocol } = new URL(wsEndpoint);
	if (protocol !== "ws:" && protocol !== "wss:") return { success: false };

	const socket = new WebSocket(wsEndpoint);
	if (socket.readyState !== socket.CONNECTING) return { success: false };

	return new Promise((res) => {
		const reply = () => res({ success: true, data: socket });

		// Connection Timeout
		setTimeout(() => {
			socket.removeEventListener("open", reply);
			res({ success: false });
		}, 10_000);

		// On Connection Succeed
		socket.addEventListener("open", reply);
	});
};

/**
 * It handle the several message sent by the server and act according to it
 * @param {WebSocket} socket
 * @param {Args} events_listeners
 */
export const HandleWsConn = (socket: WebSocket, { onMessage, onClose, onError }: Args) => {
	if (socket.readyState !== socket.OPEN) return;

	onError && socket.addEventListener("error", onError);
	onClose && socket.addEventListener("close", onClose);

	// Listen for messages
	socket.addEventListener("message", (event) => {
		if (event.data === "PING") return socket.send("OK_PING");
		if (event.data === "CLOSING") return socket.close();
		if (event.data === "LOGOUT") return LogOut(true); // received a logout instruction from server
		onMessage(event.data);
	});
};
