/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./src/**/*.{html,js,svelte,ts}"],
	theme: {
		extend: {}
	},

	plugins: [require("daisyui")],
	daisyui: {
		darkMode: "dark",
		themes: [
			{
				tenidditter: {
					...require("daisyui/src/colors/themes")["[data-theme=luxury]"]
				},
				nitter: {
					...require("daisyui/src/colors/themes")["[data-theme=night]"],
					primary: "#1DA1F2",
					secondary: "#657786",
					accent: "#AAB8C2",
					neutral: "#F5F8FA",
					"base-100": "#14171A"
				},
				teddit: {
					...require("daisyui/src/colors/themes")["[data-theme=halloween]"],
					primary: "#FF4500",
					secondary: "#92bddf",
					accent: "#5296dd",
					neutral: "#FFFFFF",
					"base-100": "#000000"
				}
			}
		]
	}
};
