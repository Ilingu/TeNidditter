export enum FeedTypeEnum {
	"Hot",
	"New",
	"Top",
	"Rising",
	"Controversial"
}
export type FeedType = "user_feed" | "home_feed";
export type FeedHomeType = "hot" | "new" | "top" | "rising" | "controversial";

/**
 * API object returned for a teddit post
 */
export interface TedditPostInfo {
	post_datas: TedditPostMetadata;
	comments: (TedditCommmentShape & { id: number; parentId: number })[][];
}

export interface TedditPostMetadata {
	post_author: string;
	post_title: string;
	/**
	 * in **second timestamp** (1s = 1000ms)
	 */
	post_created: number;
	post_ups: string;
	post_nb_comments: number;
	/**
	 * raw unparsed html without script tag
	 */
	body_html: string;
}

export interface TedditUserShape {
	username: string;
	/**
	 * @deprecated - teddit block icon's request (**CORS**)
	 */
	icon_img: string;
	created: number;
	verified: boolean;
	link_karma: number;
	comment_karma: number;
	view_more_posts: boolean;
	user_front: boolean;
	before: string;
	after: string;
	posts: (TedditPost | TedditCommmentShape)[];
}

export interface TedditCommmentShape {
	type: "t1";
	subreddit: string;
	created: number;
	subreddit_name_prefixed: string;
	ups: number;
	url: string;
	edited: boolean;
	body_html: string;
	num_comments: number;
	over_18: boolean;
	permalink: string;
	link_author: string;
	link_title: string;
	user_flair: string;
}

export interface TedditHomePageRes {
	info: Info;
	links: TedditRawPost[];
}

interface Info {
	before: unknown;
	after: string;
}

export interface TedditPost {
	id: string;
	title: string;
	author: string;
	created: number;
	/**
	 * @alias `pinned`
	 */
	stickied: boolean; // pinned by reddit

	ups: number;
	num_comments: number;

	/**
	 * if @property `url` is a external (outside of reddit.com) link, so this is set to `false`
	 */
	is_self_link: boolean;
	is_video: boolean;

	/**
	 * part of the body of the post
	 */
	selftext_html?: string; // if no media/img, display
	/**
	 * part of the body of the post
	 */
	media?: Media;
	/**
	 * part of the body of the post
	 */
	images?: TedditImages;
	/**
	 * part of the body of the post
	 */
	duration?: number;

	url: string;
	domain: string;
	permalink: string;

	subreddit: string;
	link_flair?: string;
	type?: "t3";
}

export interface TedditRawPost extends TedditPost {
	id: string;
	permalink: string;
	created: number;
	author: string;
	title: string;
	over_18: boolean;
	score: number;
	ups: number;
	upvote_ratio: number;
	num_comments: number;
	is_self_link: boolean;
	selftext_html?: string;
	url: string;
	domain: string;
	is_video: boolean;
	media?: Media;
	duration?: number;
	images?: TedditImages;
	locked: boolean;
	stickied: boolean;
	subreddit_front: unknown;
	subreddit: string;
	link_flair: string;
	user_flair: string;
	link_flair_text?: string;
}

interface Media {
	reddit_video?: RedditVideo;
	oembed?: Oembed;
	type?: string;
}

interface RedditVideo {
	bitrate_kbps: number;
	/**
	 * actual video url to put into `src`
	 */
	fallback_url: string;
	height: number;
	width: number;
	scrubber_media_url: string;
	dash_url: string;
	duration: number;
	hls_url: string;
	is_gif: boolean;
	transcoding_status: string;
}

interface Oembed {
	provider_url: string;
	url?: string;
	html: string;
	author_name: string;
	height?: number;
	width: number;
	version: string;
	author_url?: string;
	provider_name: string;
	cache_age?: number;
	type: string;
	description?: string;
	title?: string;
	thumbnail_width?: number;
	thumbnail_url?: string;
	thumbnail_height?: number;
}

interface TedditImages {
	thumb: string;
	preview?: string;
}
