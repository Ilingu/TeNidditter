import type { AlertTypes } from "./types";

export interface FunctionJob<T = never> {
	success: boolean;
	data?: T;
	error?: string;
}

export interface AlertShape {
	message: string;
	duration: number;
	type: AlertTypes;
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
	stickied: boolean; // pinned by reddit

	ups: number;
	num_comments: number;

	is_self_link: boolean;
	is_video: boolean;

	selftext_html?: string; // if no media/img, display
	media?: Media;
	images?: Images;
	duration?: number;

	url: string;
	permalink: string;

	subreddit: string;
	link_flair_text?: string;
}

export interface TedditRawPost {
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
	images?: Images;
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

interface Images {
	thumb: string;
	preview?: string;
}
