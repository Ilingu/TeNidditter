export interface AlertShape {
	message: string;
	duration: number;
	type: AlertTypes;
}
export type AlertTypes = "success" | "info" | "warning" | "error";
