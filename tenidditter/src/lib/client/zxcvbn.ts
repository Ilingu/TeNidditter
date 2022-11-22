export const GetZxcvbn = async (): Promise<typeof import("zxcvbn")> =>
	(await import("zxcvbn")).default;

export const ScoreToText = {
	0: "too weak 💀",
	1: "weak 💀",
	2: "so so 😑",
	3: "good 👍",
	4: "very good 🔒"
};
export const ScoreToColor = {
	0: "#ff6f6f",
	1: "#ff6f6f",
	2: "#e2d562",
	3: "#87d039",
	4: "#87d039"
};
