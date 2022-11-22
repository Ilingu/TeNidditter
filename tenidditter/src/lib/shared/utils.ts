/* ONLY PURE FUNCTIONS */

export const IsEmptyString = (str: unknown) => typeof str !== "string" || str.trim().length <= 0;
export const isValidUrl = (url: string): boolean => {
	try {
		new URL(url);
		return true;
	} catch {
		return false;
	}
};
export const Sleep = (dur: number) => new Promise((res) => setTimeout(res, dur));
export const ConvertHTMLEntities = (str: string): string => {
	const htmlEntities = {
		"&": "&amp;",
		"<": "&lt;",
		">": "&gt;",
		'"': "&quot;",
		"'": "&apos;"
	};
	return str.replace(/([&<>"'])/g, (match) => htmlEntities[match as keyof typeof htmlEntities]);
};

export const TrimNonDigitsChars = (str: string): string => str.replace(/[\D]+/gi, "");

/**
 * Remove all non-letter in username except "_"
 * @param {string} username
 * @returns {string} the formatted username
 */
export const FormatUsername = (username: string): string =>
	username.replace(/[\W0-9]+/g, "").toLowerCase();

export const TrimSpecialChars = (str: string): string => str.replace(/[^\w\s]+/gi, "");

export const MakeBearerToken = (JwtToken: string) => ({ Authorization: "Bearer " + JwtToken });

export const EscapeHTML = (str: string): string => new Option(str).innerHTML;

export const removeDuplicates = <T = never>(ary: T[]) => [...new Set<T>(ary)];

export const IsValidJSON = (jsonBlob: string): boolean => {
	try {
		JSON.parse(jsonBlob);
		return true;
	} catch (err) {
		return false;
	}
};
