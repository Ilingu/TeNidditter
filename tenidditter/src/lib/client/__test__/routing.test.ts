import { expect, test, describe, beforeAll } from "vitest";

describe.concurrent("Testing client routing lib", () => {
	beforeAll(() => {
		expect(process.env.TEST).toBe("true");
	});

	test.concurrent("dispatchRouting", async () => {
		const routeWillChange = (await import("../routing")).routeWillChange;
		const dispatchRouting = (await import("../routing")).dispatchRouting;

		const succeed = await new Promise<boolean>((res) => {
			routeWillChange(() => {
				res(true);
			});
			dispatchRouting();
			setTimeout(() => res(false), 1000);
		});

		expect(succeed).toBe(true);
	});
});
