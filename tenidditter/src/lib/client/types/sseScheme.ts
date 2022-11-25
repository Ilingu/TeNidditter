import type { ExternalLinksDatas } from "./nitter";

export type SSERoutes = "/nitter/stream-in-external-links";
export type SSEReturns<T> = T extends "/nitter/stream-in-external-links"
	? ExternalLinksDatas
	: never;
