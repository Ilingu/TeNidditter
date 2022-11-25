/**
 * It imports the Zxcvbn library dynamically
 * @returns {Promise<typeof import("zxcvbn")>} zxcvbn library default object
 */
export const GetZxcvbn = async (): Promise<typeof import("zxcvbn")> =>
	(await import("zxcvbn")).default;

/**
 * It converts a zxcvbn score (0-4) to an human readable string (e.g: 0->"too weak")
 */
export const ScoreToText = {
	0: "too weak 💀",
	1: "weak 💀",
	2: "so so 😑",
	3: "good 👍",
	4: "very good 🔒"
};
/**
 * It converts a zxcvbn score (0-4) to an hexadecial color (e.g: 0->"#ff6f6f", which is red)
 */
export const ScoreToColor = {
	0: "#ff6f6f",
	1: "#ff6f6f",
	2: "#e2d562",
	3: "#87d039",
	4: "#87d039"
};
