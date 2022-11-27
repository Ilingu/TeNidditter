import { sveltekit } from "@sveltejs/kit/vite";
/// <reference types="vitest" />

/** @type {import('vite').UserConfig} */
const config = {
	plugins: [sveltekit()],
	test: {
		environment: "jsdom",
		env: { SSR: false, TEST: true }
	}
};

export default config;
