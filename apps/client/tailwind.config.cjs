/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./src/**/*.{html,js,svelte,ts}"],
	theme: {
		extend: {
			colors: {
				teddit: "#FF4500",
				nitter: "#1DA1F2",
				"light-dark": "rgba(45,49,49,0.7)"
			},
			fontFamily: {
				normal: ["Fira"],
				nerd: ["Silkscreen"],
				fancy: ["Pacifico"]
			}
		}
	},

	plugins: [require("daisyui")],
	daisyui: {
		darkMode: "dark",
		themes: [
			{
				tenidditter: {
					...require("daisyui/src/colors/themes")["[data-theme=luxury]"]
				},
				nitter: require("daisyui/src/colors/themes")["[data-theme=synthwave]"],
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
