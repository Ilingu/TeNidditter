import type { AlertTypes } from "./types";

export interface FunctionJob<T = never> {
	success: boolean;
	data?: T;
}

export interface AlertShape {
	message: string;
	duration: number;
	type: AlertTypes;
}