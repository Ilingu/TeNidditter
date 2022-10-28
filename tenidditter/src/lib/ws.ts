import type { FunctionJob } from "./types/interfaces";
import { isValidUrl } from "./utils";

type Args = {
	onMessage: (data: unknown) => void;
	onError?: () => void;
	onClose?: () => void;
};
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

export const HandleWsConn = (socket: WebSocket, { onMessage, onClose, onError }: Args) => {
	if (socket.readyState !== socket.OPEN) return;

	onError && socket.addEventListener("error", onError);
	onClose && socket.addEventListener("close", onClose);

	// Listen for messages
	socket.addEventListener("message", (event) => {
		if (event.data === "PING") return socket.send("OK_PING");
		onMessage(event.data);
	});
};
