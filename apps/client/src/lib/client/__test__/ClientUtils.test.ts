import { expect, test, describe, beforeAll } from "vitest";
import {
	changeAppTheme,
	FormatElapsedTime,
	FormatNumbers,
	isMobile,
	pushAlert
} from "../ClientUtils";
import type { AlertShape } from "../types/alerts";

interface TestCase<I, E> {
	input: I;
	expected: E;
}

describe.concurrent("Testing client utils", () => {
	beforeAll(() => {
		expect(process.env.TEST).toBe("true");
	});

	test.concurrent("isMobile", () => {
		const setUserAgent = (ua: string) => {
			Object.defineProperty(navigator, "userAgent", {
				value: ua,
				configurable: true
				// writable: true
			});
		};

		const tests: TestCase<string, ReturnType<typeof isMobile>>[] = [
			{ input: navigator.userAgent, expected: false },
			{
				input:
					"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1", // safari mobile
				expected: true
			},
			{
				input:
					"Mozilla/5.0 (Linux; U; Android 4.4.2; en-us; SCH-I535 Build/KOT49H) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30", // Android Browser
				expected: true
			},
			{
				input:
					"Mozilla/5.0 (Linux; Android 7.0; SM-G930V Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36", // Chrome mobile
				expected: true
			},
			{
				input:
					"Opera/9.80 (J2ME/MIDP; Opera Mini/5.1.21214/28.2725; U; ru) Presto/2.8.119 Version/11.10",
				expected: true // opera mini
			},
			{
				input: "Mozilla/5.0 (Android 7.0; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
				expected: true // Firefox android
			},
			{
				input:
					"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/7.5b3349 Mobile/14F89 Safari/603.2.4", // firefox ios
				expected: true
			},
			{
				input:
					"Mozilla/5.0 (Linux; Android 7.0; SAMSUNG SM-G955U Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/5.4 Chrome/51.0.2704.106 Mobile Safari/537.36", // SamsungBrowser
				expected: true
			},
			{
				input:
					"Mozilla/5.0 (Linux; Android 6.0; Lenovo K50a40 Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.137 YaBrowser/17.4.1.352.00 Mobile Safari/537.36", // Yandex mobile
				expected: true
			},
			{
				input:
					"Mozilla/5.0 (Linux; U; Android 7.0; en-us; MI 5 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.146 Mobile Safari/537.36 XiaoMi/MiuiBrowser/9.0.3", // MIUI mobile
				expected: true
			},
			{
				input: "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0", // Windows Firefox
				expected: false
			},
			{
				input: "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0", // Mac Firefox
				expected: false
			},
			{
				input:
					"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36", // Desktop chrome
				expected: false
			}
		];

		for (const { input, expected } of tests) {
			setUserAgent(input);
			expect(isMobile()).toBe(expected);
		}
	});

	test.concurrent("changeAppTheme", () => {
		const tests: TestCase<
			Parameters<typeof changeAppTheme>[0],
			Parameters<typeof changeAppTheme>[0]
		>[] = [
			{ input: "nitter", expected: "nitter" },
			{ input: "teddit", expected: "teddit" },
			{ input: "tenidditter", expected: "tenidditter" }
		];

		for (const { input, expected } of tests) {
			changeAppTheme(input);
			expect(document.documentElement.getAttribute("data-theme")).toBe(expected);
		}
	});

	test.concurrent("pushAlert", async () => {
		const tests: TestCase<Parameters<typeof pushAlert>, AlertShape>[] = [
			{
				input: ["test", "success"],
				expected: { message: "test", type: "success", duration: 5000 }
			},
			{
				input: ["Invalid Username", "error", 2000],
				expected: { message: "Invalid Username", type: "error", duration: 2000 }
			},
			{
				input: ["Not logged in!", "warning", 3600],
				expected: { message: "Not logged in!", type: "warning", duration: 3600 }
			},
			{
				input: ["This feature is not implemented yet -_-", "info", 10_000],
				expected: {
					message: "This feature is not implemented yet -_-",
					type: "info",
					duration: 10_000
				}
			}
		];

		let eventsReceived = 0;
		document.addEventListener("alert", () => eventsReceived++);

		for (const { input, expected } of tests) {
			const succeed = await new Promise<boolean>((resolve) => {
				const checkEvContent = ({ detail: alertEvt }: CustomEvent<AlertShape>) => {
					expect(alertEvt).toStrictEqual(expected);
					document.removeEventListener("alert", checkEvContent as EventListener);
					resolve(true);
				};

				document.addEventListener("alert", checkEvContent as EventListener);
				setTimeout(() => resolve(false), 5000);

				pushAlert(...input);
			});
			expect(succeed).toBe(true);
		}

		expect(eventsReceived).toBe(tests.length);
	});

	test.concurrent("FormatNumbers", () => {
		const tests: TestCase<Parameters<typeof FormatNumbers>[0], ReturnType<typeof FormatNumbers>>[] =
			[
				{ input: NaN, expected: "NaN" },
				{ input: 2, expected: "2" },
				{ input: 25, expected: "25" },
				{ input: 256, expected: "256" },
				{ input: 2_566, expected: "2.6K" },
				{ input: 12_200, expected: "12K" },
				{ input: 28_000, expected: "28K" },
				{ input: 28_200, expected: "28K" },
				{ input: 28_500, expected: "29K" },
				{ input: 28_700, expected: "29K" },
				{ input: 256_600, expected: "257K" },
				{ input: 2_566_000, expected: "2.6M" },
				{ input: 25_660_000, expected: "26M" },
				{ input: 256_600_000, expected: "257M" },
				{ input: 2_566_000_000, expected: "2.6B" },
				{ input: 25_660_000_000, expected: "26B" },
				{ input: 256_600_000_000, expected: "257B" },
				{ input: 2_566_000_000_000, expected: "2.6T" }
			];

		for (const { input, expected } of tests) {
			expect(FormatNumbers(input)).toBe(expected);
		}
	});

	test.concurrent("FormatElapsedTime", async () => {
		const now = Date.now();

		const MINUTE = 60 * 1000,
			HOUR = 60 * MINUTE,
			DAY = 24 * HOUR,
			MONTH = 30.4375 * DAY,
			YEAR = 12 * MONTH;
		const tests: TestCase<
			Parameters<typeof FormatElapsedTime>,
			ReturnType<typeof FormatElapsedTime>
		>[] = [
			// MIN
			{ input: [now, now + 1 * MINUTE], expected: "1 min ago" },
			{ input: [now, now + 15 * MINUTE], expected: "15 min ago" },
			{ input: [now, now + 24.6 * MINUTE], expected: "25 min ago" },
			{ input: [now, now + 59 * MINUTE], expected: "59 min ago" },
			{ input: [now, now + 59.75 * MINUTE], expected: "60 min ago" },
			// HOURS
			{ input: [now, now + 120 * MINUTE], expected: "2h ago" },
			{ input: [now, now + 1 * HOUR], expected: "1h ago" },
			{ input: [now, now + 5 * HOUR], expected: "5h ago" },
			{ input: [now, now + 9.2 * HOUR], expected: "9h ago" },
			{ input: [now, now + 23 * HOUR], expected: "23h ago" },
			{ input: [now, now + 23.9 * HOUR], expected: "24h ago" },
			{ input: [now, now + 25 * HOUR], expected: "1d ago" },
			// DAYS
			{ input: [now, now + 1440 * MINUTE], expected: "1d ago" },
			{ input: [now, now + 1 * DAY], expected: "1d ago" },
			{ input: [now, now + 48 * HOUR], expected: "2d ago" },
			{ input: [now, now + 7 * DAY], expected: "7d ago" },
			{ input: [now, now + 30 * DAY], expected: "30d ago" },
			{ input: [now, now + 30.4375 * DAY], expected: "1m ago" },
			{ input: [now, now + 31.6 * DAY], expected: "1m ago" },
			// MONTHS
			{ input: [now, now + 62 * DAY], expected: "2m ago" },
			{ input: [now, now + 1 * MONTH], expected: "1m ago" },
			{ input: [now, now + 6.456 * MONTH], expected: "6m ago" },
			{ input: [now, now + 11.6 * MONTH], expected: "12m ago" },
			{ input: [now, now + 12 * MONTH], expected: "1y ago" },
			// YEARS
			{ input: [now, now + 365 * DAY], expected: "12m ago" },
			{ input: [now, now + 365.25 * DAY], expected: "1y ago" },
			{ input: [now, now + 22.7 * MONTH], expected: "2y ago" },
			{ input: [now, now + 1 * YEAR], expected: "1y ago" },
			{ input: [now, now + 1652 * YEAR], expected: "1652y ago" }
		];

		for (const { input, expected } of tests) {
			expect(FormatElapsedTime(...input)).toBe(expected);
		}
	});
});
