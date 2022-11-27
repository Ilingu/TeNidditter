import { expect, test, describe, expectTypeOf, beforeAll } from "vitest";
import { GetZxcvbn, ScoreToColor, ScoreToText } from "../zxcvbn";
import type { ZXCVBNResult } from "zxcvbn";

interface TestCase<I, E> {
	input: I;
	excepted: E;
}

describe.concurrent("Testing client zxcvbn lib", () => {
	beforeAll(() => {
		expect(process.env.TEST).toBe("true");
	});

	test.concurrent("GetZxcvbn", async () => {
		const zxcvbn = await GetZxcvbn();
		expectTypeOf(zxcvbn).toBeFunction();
		expectTypeOf(zxcvbn).toEqualTypeOf(
			(password: string, userInputs?: string[] | undefined): ZXCVBNResult =>
				zxcvbn(password, userInputs)
		);
	});

	test.concurrent("Zxcvbn Password", async () => {
		const zxcvbn = await GetZxcvbn();

		const tests: TestCase<Parameters<typeof zxcvbn>[0], ZXCVBNResultOptional>[] = [
			{
				input: "test",
				excepted: {
					password: "test",
					guesses: 94,
					guesses_log10: 1.9731278535996983,
					sequence: [
						{
							pattern: "dictionary",
							i: 0,
							j: 3,
							token: "test",
							matched_word: "test",
							rank: 93,
							dictionary_name: "passwords",
							reversed: false,
							l33t: false,
							base_guesses: 93,
							uppercase_variations: 1,
							l33t_variations: 1,
							guesses: 93,
							guesses_log10: 1.968482948553935
						}
					],
					crack_times_seconds: {
						online_throttling_100_per_hour: 3384,
						online_no_throttling_10_per_second: 9.4,
						offline_slow_hashing_1e4_per_second: 0.0094,
						offline_fast_hashing_1e10_per_second: 9.4e-9
					},
					crack_times_display: {
						online_throttling_100_per_hour: "56 minutes",
						online_no_throttling_10_per_second: "9 seconds",
						offline_slow_hashing_1e4_per_second: "less than a second",
						offline_fast_hashing_1e10_per_second: "less than a second"
					},
					score: 0,
					feedback: {
						warning: "This is a top-100 common password",
						suggestions: ["Add another word or two. Uncommon words are better."]
					}
				}
			},
			{
				input: "golangTypescript",
				excepted: {
					password: "golangTypescript",
					guesses: 2077542016000,
					guesses_log10: 12.317549815672722,
					sequence: [
						{
							pattern: "dictionary",
							i: 0,
							j: 1,
							token: "go",
							matched_word: "go",
							rank: 43,
							dictionary_name: "us_tv_and_film",
							reversed: false,
							l33t: false,
							base_guesses: 43,
							uppercase_variations: 1,
							l33t_variations: 1,
							guesses: 50,
							guesses_log10: 1.6989700043360185
						},
						{
							pattern: "dictionary",
							i: 2,
							j: 5,
							token: "lang",
							matched_word: "lang",
							rank: 478,
							dictionary_name: "surnames",
							reversed: false,
							l33t: false,
							base_guesses: 478,
							uppercase_variations: 1,
							l33t_variations: 1,
							guesses: 478,
							guesses_log10: 2.6794278966121188
						},
						{
							pattern: "dictionary",
							i: 6,
							j: 9,
							token: "Type",
							matched_word: "type",
							rank: 472,
							dictionary_name: "english_wikipedia",
							reversed: false,
							l33t: false,
							base_guesses: 472,
							uppercase_variations: 2,
							l33t_variations: 1,
							guesses: 944,
							guesses_log10: 2.9749719942980684
						},
						{
							pattern: "dictionary",
							i: 10,
							j: 15,
							token: "script",
							matched_word: "script",
							rank: 1990,
							dictionary_name: "english_wikipedia",
							reversed: false,
							l33t: false,
							base_guesses: 1990,
							uppercase_variations: 1,
							l33t_variations: 1,
							guesses: 1990,
							guesses_log10: 3.2988530764097064
						}
					],
					crack_times_seconds: {
						online_throttling_100_per_hour: 74791512576000,
						online_no_throttling_10_per_second: 207754201600,
						offline_slow_hashing_1e4_per_second: 207754201.6,
						offline_fast_hashing_1e10_per_second: 207.7542016
					},
					crack_times_display: {
						online_throttling_100_per_hour: "centuries",
						online_no_throttling_10_per_second: "centuries",
						offline_slow_hashing_1e4_per_second: "6 years",
						offline_fast_hashing_1e10_per_second: "3 minutes"
					},
					score: 4,
					feedback: {
						warning: "",
						suggestions: []
					}
				}
			},
			{
				input: "H3^v6WJ3Qb*%9^pXw00TmY8##31XE5xT",
				excepted: {
					password: "H3^v6WJ3Qb*%9^pXw00TmY8##31XE5xT",
					guesses: 1e32,
					guesses_log10: 32,
					sequence: [
						{
							pattern: "bruteforce",
							token: "H3^v6WJ3Qb*%9^pXw00TmY8##31XE5xT",
							i: 0,
							j: 31,
							guesses: 1e32,
							guesses_log10: 32
						}
					],
					crack_times_seconds: {
						online_throttling_100_per_hour: 3.6e33,
						online_no_throttling_10_per_second: 1.0000000000000001e31,
						offline_slow_hashing_1e4_per_second: 1e28,
						offline_fast_hashing_1e10_per_second: 1e22
					},
					crack_times_display: {
						online_throttling_100_per_hour: "centuries",
						online_no_throttling_10_per_second: "centuries",
						offline_slow_hashing_1e4_per_second: "centuries",
						offline_fast_hashing_1e10_per_second: "centuries"
					},
					score: 4,
					feedback: {
						warning: "",
						suggestions: []
					}
				}
			}
		];

		for (const { input, excepted } of tests) {
			const res = zxcvbn(input) as ZXCVBNResultOptional;
			delete res.calc_time;
			expect(res).toStrictEqual(excepted);
		}
	});

	test.concurrent("Test ScoreToText", async () => {
		type ValueOf<T> = T[keyof T];
		const tests: TestCase<keyof typeof ScoreToText, ValueOf<typeof ScoreToText>>[] = [
			{
				input: 0,
				excepted: "too weak 💀"
			},
			{
				input: 1,
				excepted: "weak 💀"
			},
			{
				input: 2,
				excepted: "so so 😑"
			},
			{
				input: 3,
				excepted: "good 👍"
			},
			{
				input: 4,
				excepted: "very good 🔒"
			}
		];

		for (const { input, excepted } of tests) {
			expect(ScoreToText[input]).toBe(excepted);
		}
	});

	test.concurrent("Test ScoreToColor", async () => {
		type ValueOf<T> = T[keyof T];
		const tests: TestCase<keyof typeof ScoreToColor, ValueOf<typeof ScoreToColor>>[] = [
			{
				input: 0,
				excepted: "#ff6f6f"
			},
			{
				input: 1,
				excepted: "#ff6f6f"
			},
			{
				input: 2,
				excepted: "#e2d562"
			},
			{
				input: 3,
				excepted: "#87d039"
			},
			{
				input: 4,
				excepted: "#87d039"
			}
		];

		for (const { input, excepted } of tests) {
			expect(ScoreToColor[input]).toBe(excepted);
		}
	});
});

interface ZXCVBNResultOptional {
	password: string;
	guesses: number;
	guesses_log10: number;
	crack_times_seconds: Partial<ZXCVBNAttackTime>;
	crack_times_display: Partial<ZXCVBNAttackTime>;
	score: Partial<ZXCVBNScore>;
	feedback: Partial<ZXCVBNFeedback>;
	sequence?: Partial<ZXCVBNSequence>[];
	calc_time?: number;
}

type ZXCVBNScore = 0 | 1 | 2 | 3 | 4;

interface ZXCVBNAttackTime {
	online_throttling_100_per_hour: string | number;
	online_no_throttling_10_per_second: string | number;
	offline_slow_hashing_1e4_per_second: string | number;
	offline_fast_hashing_1e10_per_second: string | number;
}

interface ZXCVBNFeedback {
	warning: Partial<ZXCVBNFeedbackWarning>;
	suggestions: string[];
}

type ZXCVBNFeedbackWarning =
	| "Straight rows of keys are easy to guess"
	| "Short keyboard patterns are easy to guess"
	| "Use a longer keyboard pattern with more turns"
	| 'Repeats like "aaa" are easy to guess'
	| 'Repeats like "abcabcabc" are only slightly harder to guess than "abc"'
	| "Sequences like abc or 6543 are easy to guess"
	| "Recent years are easy to guess"
	| "Dates are often easy to guess"
	| "This is a top-10 common password"
	| "This is a top-100 common password"
	| "This is a very common password"
	| "This is similar to a commonly used password"
	| "A word by itself is easy to guess"
	| "Names and surnames by themselves are easy to guess"
	| "Common names and surnames are easy to guess"
	| "";

interface ZXCVBNSequence {
	ascending: boolean;
	base_guesses: number;
	base_matches: string;
	base_token: string;
	dictionary_name: string;
	guesses: number;
	guesses_log10: number;
	i: number;
	j: number;
	l33t: boolean;
	l33t_variations: number;
	matched_word: string;
	pattern: string;
	rank: number;
	repeat_count: number;
	reversed: boolean;
	sequence_name: string;
	sequence_space: number;
	token: string;
	uppercase_variations: number;
}
