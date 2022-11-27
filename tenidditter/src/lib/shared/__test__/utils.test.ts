import { expect, test, describe, beforeAll } from "vitest";
import {
	EscapeHTML,
	FormatUsername,
	IsEmptyString,
	IsValidJSON,
	isValidUrl,
	MakeBearerToken,
	removeDuplicates,
	Sleep,
	TrimNonDigitsChars,
	TrimSpecialChars
} from "../utils";

interface TestCase<I, E> {
	input: I;
	excepted: E;
}

describe.concurrent("Testing shared utils", () => {
	beforeAll(() => {
		expect(process.env.TEST).toBe("true");
	});

	test.concurrent("IsEmptyString", () => {
		const tests: TestCase<Parameters<typeof IsEmptyString>[0], ReturnType<typeof IsEmptyString>>[] =
			[
				{ input: null, excepted: true },
				{ input: undefined, excepted: true },
				{ input: [28, 11], excepted: true },
				{ input: { a: true, b: "yolo", c: ["a"] }, excepted: true },
				{ input: "", excepted: true },
				{ input: " ", excepted: true },
				{ input: "a", excepted: false },
				{ input: " AWERTY_ A ", excepted: false }
			];

		for (const { input, excepted } of tests) {
			expect(IsEmptyString(input)).toBe(excepted);
		}
	});

	test.concurrent("isValidUrl", () => {
		const tests: TestCase<Parameters<typeof isValidUrl>[0], ReturnType<typeof isValidUrl>>[] = [
			{ input: "", excepted: false },
			{ input: "not an url", excepted: false },
			{ input: "exemple.com", excepted: false },
			{ input: "http exemple.com", excepted: false },
			{ input: "://ack.vercel.app", excepted: false },
			{ input: "/some/route", excepted: false },

			{ input: "http://ack.vercel.app", excepted: true },
			{ input: "https://ack.vercel.app/some/route", excepted: true },
			{ input: "https://ack.vercel.app/some/route?test=true&vitest=great", excepted: true }
		];

		for (const { input, excepted } of tests) {
			expect(isValidUrl(input)).toBe(excepted);
		}
	});

	test.concurrent(
		"Sleep",
		async () => {
			const tests: TestCase<Parameters<typeof Sleep>[0], number>[] = [
				{ input: 10, excepted: 2 },
				{ input: 100, excepted: 10 },
				{ input: 1000, excepted: 10 }
				// { input: 10_000, excepted: 10 }
			];

			for (const { input, excepted } of tests) {
				const start = Date.now();
				await Sleep(input);

				const elapsedTime = Date.now() - start;
				expect(input - excepted).toBeLessThanOrEqual(elapsedTime);
				expect(input + excepted).toBeGreaterThanOrEqual(elapsedTime);
			}
		},
		{ timeout: 11_100 + 1000 }
	);

	test.concurrent("EscapeHTML", () => {
		const tests: TestCase<Parameters<typeof EscapeHTML>[0], ReturnType<typeof EscapeHTML>>[] = [
			{ input: "some string", excepted: "some string" },
			{ input: "<a></a>", excepted: "&lt;a&gt;&lt;/a&gt;" },
			{
				input: `<div><div class="js-notice js-hide-dashboard-item rounded-2" data-color-mode="dark" data-light-theme="light" data-dark-theme="dark"><form class="Box position-relative rounded-2 mb-4 p-3 js-notice-dismiss overflow-hidden" style="z-index: 1; color: #cdd9e4 !important;" data-turbo="false" action="/settings/dismiss-notice/dashboard_promo_universe_22_on_demand" accept-charset="UTF-8" method="post"><input type="hidden" name="authenticity_token" value="8Yuoj3XiVdnQSgPdImmhXBr1t-i7eiE8TH_PSQ4y7S75PpTJN_YE5prW1pErlOnPs7bJU6Rkmpb9Qn0wGjS-BQ"> <picture> <source srcset="https://github.githubassets.com/images/modules/dashboard/universe22/bg.webp" type="image/webp"> <img src="https://github.githubassets.com/images/modules/dashboard/universe22/bg.jpg" alt="" width="604" height="450" class="position-absolute top-0 left-0 width-full" style="pointer-events: none; z-index: -1; height: 100%; height: 100%; object-fit: cover"> </picture> <div class="position-absolute p-2" style="top: 4px; right: 6px;"> <button style="color: currentColor" aria-label="Close" type="submit" data-view-component="true" class="close-button"><svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-x"> <path fill-rule="evenodd" d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"></path> </svg></button> </div> <img src="https://github.githubassets.com/images/modules/dashboard/universe22/universe22-logo.svg" alt="GitHub Universe 2022" width="173" height="24" style="max-width: 85%; height: auto" class="d-block"> <h3 class="h5 mb-1 mt-3">Let's build from here</h3> <p class="mb-3 f5"> Watch all the latest product announcements and expert-driven sessions from this year's event, available now on-demand. </p> <style> .btn-universe22-dashboard { background-color: rgba(255, 255, 255, 0.16) !important; border-color: #adbac6 !important; color: #cdd9e4 !important; } .btn-universe22-dashboard:hover, .btn-universe22-dashboard:focus { background-color: rgba(255, 255, 255, 0.08) !important; border-color: #8b949e !important; } .btn-universe22-dashboard:active { background-color: rgba(255, 255, 255, 0.1) !important; border-color: #6e7681 !important; } </style> <a href="https://watch.githubuniverse.com/on-demand/?utm_source=github&amp;utm_medium=product&amp;utm_campaign=2022-product-WW-Universe&amp;utm_content=watch-on-demand-promo" target="_blank" data-analytics-event="{&quot;category&quot;:&quot;Dashboard notices&quot;,&quot;action&quot;:&quot;click to watch Universe 2022 on demand&quot;,&quot;label&quot;:&quot;ref_page:/&quot;}" data-view-component="true" class="btn-universe22-dashboard btn btn-block">    Watch now </a></form>  </div> </div>`,
				excepted: `&lt;div&gt;&lt;div class="js-notice js-hide-dashboard-item rounded-2" data-color-mode="dark" data-light-theme="light" data-dark-theme="dark"&gt;&lt;form class="Box position-relative rounded-2 mb-4 p-3 js-notice-dismiss overflow-hidden" style="z-index: 1; color: #cdd9e4 !important;" data-turbo="false" action="/settings/dismiss-notice/dashboard_promo_universe_22_on_demand" accept-charset="UTF-8" method="post"&gt;&lt;input type="hidden" name="authenticity_token" value="8Yuoj3XiVdnQSgPdImmhXBr1t-i7eiE8TH_PSQ4y7S75PpTJN_YE5prW1pErlOnPs7bJU6Rkmpb9Qn0wGjS-BQ"&gt; &lt;picture&gt; &lt;source srcset="https://github.githubassets.com/images/modules/dashboard/universe22/bg.webp" type="image/webp"&gt; &lt;img src="https://github.githubassets.com/images/modules/dashboard/universe22/bg.jpg" alt="" width="604" height="450" class="position-absolute top-0 left-0 width-full" style="pointer-events: none; z-index: -1; height: 100%; height: 100%; object-fit: cover"&gt; &lt;/picture&gt; &lt;div class="position-absolute p-2" style="top: 4px; right: 6px;"&gt; &lt;button style="color: currentColor" aria-label="Close" type="submit" data-view-component="true" class="close-button"&gt;&lt;svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-x"&gt; &lt;path fill-rule="evenodd" d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"&gt;&lt;/path&gt; &lt;/svg&gt;&lt;/button&gt; &lt;/div&gt; &lt;img src="https://github.githubassets.com/images/modules/dashboard/universe22/universe22-logo.svg" alt="GitHub Universe 2022" width="173" height="24" style="max-width: 85%; height: auto" class="d-block"&gt; &lt;h3 class="h5 mb-1 mt-3"&gt;Let's build from here&lt;/h3&gt; &lt;p class="mb-3 f5"&gt; Watch all the latest product announcements and expert-driven sessions from this year's event, available now on-demand. &lt;/p&gt; &lt;style&gt; .btn-universe22-dashboard { background-color: rgba(255, 255, 255, 0.16) !important; border-color: #adbac6 !important; color: #cdd9e4 !important; } .btn-universe22-dashboard:hover, .btn-universe22-dashboard:focus { background-color: rgba(255, 255, 255, 0.08) !important; border-color: #8b949e !important; } .btn-universe22-dashboard:active { background-color: rgba(255, 255, 255, 0.1) !important; border-color: #6e7681 !important; } &lt;/style&gt; &lt;a href="https://watch.githubuniverse.com/on-demand/?utm_source=github&amp;amp;utm_medium=product&amp;amp;utm_campaign=2022-product-WW-Universe&amp;amp;utm_content=watch-on-demand-promo" target="_blank" data-analytics-event="{&amp;quot;category&amp;quot;:&amp;quot;Dashboard notices&amp;quot;,&amp;quot;action&amp;quot;:&amp;quot;click to watch Universe 2022 on demand&amp;quot;,&amp;quot;label&amp;quot;:&amp;quot;ref_page:/&amp;quot;}" data-view-component="true" class="btn-universe22-dashboard btn btn-block"&gt;    Watch now &lt;/a&gt;&lt;/form&gt;  &lt;/div&gt; &lt;/div&gt;`
			}
		];

		for (const { input, excepted } of tests) {
			expect(EscapeHTML(input)).toBe(excepted);
		}
	});

	test.concurrent("TrimNonDigitsChars", () => {
		const tests: TestCase<
			Parameters<typeof TrimNonDigitsChars>[0],
			ReturnType<typeof TrimNonDigitsChars>
		>[] = [
			{
				input: "uk64gh[]{};':sad8734?/kjsaf09[wr-88+RO<>,.OTOF2	089;8GEWR8^&*%@#>,82679OY23QR",
				excepted: "64873409882089888267923"
			}
		];

		for (const { input, excepted } of tests) {
			expect(TrimNonDigitsChars(input)).toBe(excepted);
		}
	});

	test.concurrent("FormatUsername", () => {
		const tests: TestCase<
			Parameters<typeof FormatUsername>[0],
			ReturnType<typeof FormatUsername>
		>[] = [
			{
				input: "Ilingu",
				excepted: "ilingu"
			},
			{
				input: "saSa-miYa3000",
				excepted: "sasamiya"
			},
			{
				input: "XXX_(I'm ridiculous 1)-*_*-#swag+=+3454_xXx",
				excepted: "xxx_imridiculous_swag_xxx"
			}
		];

		for (const { input, excepted } of tests) {
			expect(FormatUsername(input)).toBe(excepted);
		}
	});

	test.concurrent("TrimSpecialChars", () => {
		const tests: TestCase<
			Parameters<typeof TrimSpecialChars>[0],
			ReturnType<typeof TrimSpecialChars>
		>[] = [
			{
				input: "Ilingu",
				excepted: "Ilingu"
			},
			{
				input: "saSa-miYa3000",
				excepted: "saSamiYa3000"
			},
			{
				input: "I love yaoi",
				excepted: "I love yaoi"
			},
			{
				input: "I have 32+BL's in my librar|",
				excepted: "I have 32BLs in my librar"
			},
			{
				input: `A!@#$%^&*()_+{}|\\":';[]<>?/.,=-~A`,
				excepted: "A_A"
			}
		];

		for (const { input, excepted } of tests) {
			expect(TrimSpecialChars(input)).toBe(excepted);
		}
	});

	test.concurrent("MakeBearerToken", () => {
		const tests: TestCase<
			Parameters<typeof MakeBearerToken>[0],
			ReturnType<typeof MakeBearerToken>
		>[] = [
			{
				input: "Ilingu",
				excepted: { Authorization: "Bearer Ilingu" }
			},
			{
				input:
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwibmFtZSI6ImlsaW5ndSIsImV4cCI6MTY2ODYxNTQxN30.7HUMg5vYTwmmBRghUSw5oGlGAXBrTLT3gjBGksu0WuU",
				excepted: {
					Authorization:
						"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwibmFtZSI6ImlsaW5ndSIsImV4cCI6MTY2ODYxNTQxN30.7HUMg5vYTwmmBRghUSw5oGlGAXBrTLT3gjBGksu0WuU"
				}
			}
		];

		for (const { input, excepted } of tests) {
			expect(MakeBearerToken(input)).toStrictEqual(excepted);
		}
	});

	test.concurrent("removeDuplicates", () => {
		const SameMemoryRef = [2, 5, 8, "2"];
		const SameFn = (someArgs: string): number[] => [parseInt(someArgs)];

		const DuplicatesTest: TestCase<
			Parameters<typeof removeDuplicates>[0],
			ReturnType<typeof removeDuplicates>
		>[] = [
			{
				input: [1, 2, 5, 3, 2, 4, 1001, 8, 5, 3, 1001],
				excepted: [1, 2, 5, 3, 4, 1001, 8]
			},
			{
				input: [
					"abc",
					5,
					"d",
					"http://localhost:3000",
					"abc",
					1,
					"d",
					1,
					"http://localhost:3000",
					5
				],
				excepted: ["abc", 5, "d", "http://localhost:3000", 1]
			},
			{
				input: [SameMemoryRef, SameMemoryRef, 28],
				excepted: [SameMemoryRef, 28]
			},
			{
				input: [SameFn, SameFn, 28],
				excepted: [SameFn, 28]
			}
		];

		for (const { input, excepted } of DuplicatesTest) {
			expect(removeDuplicates(input)).toEqual(excepted);
		}
	});
	test.concurrent("IsValidJSON", async () => {
		const tests: TestCase<Parameters<typeof IsValidJSON>[0], ReturnType<typeof IsValidJSON>>[] = [
			{
				input: "Ilingu",
				excepted: false
			},
			{
				input: "{Ilingu}",
				excepted: false
			},
			{
				input: `{"message": "invalid or expired jwt"}`,
				excepted: true
			},
			{
				input: `{ "success": true,
				 "code": 202, "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwibmFtZSI6ImlsaW5ndSIsImV4cCI6MTY2OTc1NzE4OX0.hC2IUYFiftwx51WH_-Psxq4uf6l__3QLWy8M0UWYQwE", "error": "" }`,
				excepted: true
			},
			{
				input: `{"message: "invalid or expired jwt"}`,
				excepted: false
			},
			{
				input: `{"message" "invalid or expired jwt"}`,
				excepted: false
			},
			{
				input: `{"message": "invalid or expired jwt"`,
				excepted: false
			},
			{
				input: `{"message": "invalid or expired jwt",}`,
				excepted: false
			},
			{
				input: JSON.stringify((await import("./jsonTest.json")).default),
				excepted: true
			}
		];

		for (const { input, excepted } of tests) {
			expect(IsValidJSON(input)).toBe(excepted);
		}
	});
});
