export const GetZxcvbn = async (): Promise<typeof import("zxcvbn")> =>
	(await import("zxcvbn")).default;

export const ScoreToText = {
	0: "too weak",
	1: "weak",
	2: "so so",
	3: "good",
	4: "very good"
};
