// import adapter from "@sveltejs/adapter-auto";
import adapter from "@sveltejs/adapter-vercel";
import preprocess from "svelte-preprocess";

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: preprocess({ postcss: true }),

	kit: {
		adapter: adapter({ edge: true, split: false })
	}
};

export default config;
