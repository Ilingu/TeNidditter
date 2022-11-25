/// <reference types="vitest" />

import { defineConfig } from "vite";

// Configure Vitest (https://vitest.dev/config/)
export default defineConfig({
	test: {
		environment: "jsdom"
	}
});
