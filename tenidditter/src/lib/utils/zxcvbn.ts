export const GetZxcvbn = async (): Promise<typeof import("zxcvbn")> =>
	(await import("zxcvbn")).default;
