/* ONLY PURE FUNCTIONS */

/**
 * Return whether the input is of type string and if it's empty string or not.
 * @param {unknown} val - value of unknow type
 * @returns {boolean} if `true` this is either not a string or an empty string
 */
export const IsEmptyString = (val: unknown): boolean =>
	typeof val !== "string" || val.trim().length <= 0;

/**
 * Return whether the input url is a valid HTTP url or not.
 * @param {string} url - url to check
 * @returns {boolean} if `true` this is a valid http url
 */
export const isValidUrl = (url: string): boolean => {
	try {
		new URL(url);
		return true;
	} catch {
		return false;
	}
};

/**
 * Pause your code execution for the amount of millisecond specified
 * @param {number} dur - duration (`in ms`) of the pause/sleep
 * @returns {Promise<unknown>} a promise to `await`
 */
export const Sleep = (dur: number): Promise<unknown> => new Promise((res) => setTimeout(res, dur));

/**
 * Convert your unparsed html into HTML Entities
 * @see {@link https://www.w3schools.com/html/html_entities.asp| W3S HTML Entities}
 * @param {string} str - unparsed html
 * @returns {string} HTML Entities parsed
 */
export const EscapeHTML = (str: string): string => new Option(str).innerHTML;
/* OLD VERSION OF "EscapeHTML"
const ConvertHTMLEntities = (str: string): string => {
	const htmlEntities = {
		"&": "&amp;",
		"<": "&lt;",
		">": "&gt;",
		'"': "&quot;",
		"'": "&apos;"
	};
	return str.replace(/([&<>"'])/g, (match) => htmlEntities[match as keyof typeof htmlEntities]);
};
*/

/**
 * Removes all non digits characters in a string
 * @param {string} str - string containing alphanumerics/ponctuation/symbol... characters
 * @returns {string} string containing only numerics characters
 */
export const TrimNonDigitsChars = (str: string): string => str.replace(/[\D]+/gi, "");

/**
 * Removes all non-letter in username except "_"
 * @param {string} username
 * @returns {string} the formatted username
 */
export const FormatUsername = (username: string): string =>
	username.replace(/[\W0-9]+/g, "").toLowerCase();

/**
 * Removes all Specials characters in a string (ponctuals, symbol...)
 * @param {string} str - string containing specials characters
 * @returns {string} string without specials characters
 */
export const TrimSpecialChars = (str: string): string => str.replace(/[^\w\s]+/gi, "");

/**
 * Creates a Bearer Token Authorization header object
 * @param {string} JwtToken - user's jwt_token
 * @returns {{Authorization: string;}} Authorization header object
 */
export const MakeBearerToken = (JwtToken: string): { Authorization: string } => ({
	Authorization: "Bearer " + JwtToken
});

/**
 * removes all the duplicates values in a given array
 * @param {T} ary - input array of any type
 * @returns {T[]} same array without all the duplicated values
 */
export const removeDuplicates = <T = never>(ary: T[]): T[] => [...new Set<T>(ary)];

/**
 * Whether the input string is a valid JSON parsable object or not
 * @param {string} jsonBlob
 * @returns {boolean} if `true`, this is a valid, `JSON.parse` safe JSON object string
 */
export const IsValidJSON = (jsonBlob: string): boolean => {
	try {
		JSON.parse(jsonBlob);
		return true;
	} catch (err) {
		return false;
	}
};
