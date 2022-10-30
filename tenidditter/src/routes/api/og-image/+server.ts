import { error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import satori from "satori";
import { FormatNumbers, humanElapsedTime, IsEmptyString } from "$lib/utils";

interface ReactNode {
	type: keyof HTMLElementTagNameMap;
	props?: {
		children?: ReactNode | string | (ReactNode | string)[];
		style?: object;
		src?: string;
		width?: number;
		height?: number;
	};
}

export const GET: RequestHandler = async ({ url, fetch }) => {
	const title = decodeURIComponent(url.searchParams.get("title") ?? "");
	const author = decodeURIComponent(url.searchParams.get("author") ?? "");
	const subreddit = decodeURIComponent(url.searchParams.get("subreddit") ?? "");
	const ups = Number(url.searchParams.get("ups") ?? "NaN");
	const created = Number(url.searchParams.get("created") ?? "NaN");

	if (
		IsEmptyString(title) ||
		IsEmptyString(author) ||
		IsEmptyString(subreddit) ||
		isNaN(ups) ||
		isNaN(created)
	)
		throw error(400, "invalid args");

	const HTMLObject: ReactNode = {
		type: "div",
		props: {
			children: [
				{
					type: "h1",
					props: {
						children: title,
						style: {
							color: "white",
							fontSize: "32px"
						}
					}
				},
				{
					type: "p",
					props: {
						children: [
							`Submitted `,
							{
								type: "span",
								props: {
									children: humanElapsedTime(created * 1000, Date.now()),
									style: { color: "#1DA1F2", marginLeft: "10px", marginRight: "10px" }
								}
							},
							" by ",
							{
								type: "span",
								props: {
									children: author,
									style: { color: "#1DA1F2", marginLeft: "10px", marginRight: "10px" }
								}
							},
							" on ",
							{
								type: "span",
								props: {
									children: `r/${subreddit}`,
									style: { color: "#FF4500", marginLeft: "10px" }
								}
							}
						],
						style: {
							color: "white",
							fontSize: "22px"
						}
					}
				},
				{
					type: "p",
					props: {
						children: [
							{
								type: "span",
								props: {
									children: "⬆",
									style: { marginRight: "10px", transform: "translateY(3px)" }
								}
							},
							{
								type: "span",
								props: {
									children: FormatNumbers(ups),
									style: { textTransform: "lowercase" }
								}
							}
						],
						style: {
							position: "absolute",
							top: "5px",
							right: "15px",
							backgroundColor: "#FF4500",
							display: "flex",
							justifyContent: "center",
							borderRadius: "5px",
							minWidth: "64px",
							height: "50px",
							fontSize: "25px",
							padding: "0 10px",
							color: "white"
						}
					}
				},
				{
					type: "img",
					props: {
						style: {
							position: "absolute",
							top: "5px",
							left: "5px"
						},
						src: `${url.origin}/icons/IconTNDT192.png`,
						width: 96,
						height: 96
					}
				},
				{
					type: "p",
					props: {
						children: "TeNidditter",
						style: {
							color: "lightgray",
							position: "absolute",
							bottom: "5px",
							right: "15px"
						}
					}
				}
			],
			style: {
				height: "100%",
				width: "100%",
				position: "relative",
				display: "flex",
				textAlign: "center",
				alignItems: "center",
				justifyContent: "center",
				flexDirection: "column",
				flexWrap: "nowrap",
				backgroundColor: "#09090b",
				backgroundImage:
					"radial-gradient(circle at 25px 25px, #FF4500 2%, transparent 0%), radial-gradient(circle at 75px 75px, #1DA1F2 2%, transparent 0%)",
				backgroundSize: "100px 100px"
			}
		}
	};

	try {
		const svg = await satori(HTMLObject, {
			width: 800,
			height: 400,
			fonts: [
				{
					name: "Pacifico-Regular",
					data: await (await fetch("/Assets/fonts/Pacifico/Pacifico-Regular.ttf")).arrayBuffer(),
					weight: 400,
					style: "normal"
				}
			],
			graphemeImages: {
				"⬆": "https://twemoji.maxcdn.com/v/latest/svg/2b06.svg"
			}
		});

		return new Response(svg, { headers: { "Content-Type": "image/svg+xml" } });
	} catch (err) {
		console.error(err);
		throw error(500, "cannot generate image");
	}
};
