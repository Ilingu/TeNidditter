export interface NeetComment extends NeetBasicComment {
	quote?: NeetBasicComment;
}

export interface NeetBasicComment {
	id: string;
	content: string;
	creator: NittosPreview;
	createdAt: number;
	stats: NeetCommentStats;
	attachment?: Attachments;
	externalLink?: string;
	retweeted?: string;
	pinned?: boolean;
}
export interface Attachments {
	images?: string[];
	/**
	 * Hls video urls
	 */
	videos?: string[];
}
export interface NeetCommentStats {
	reply_counts: number;
	rt_counts: number;
	quotes_counts?: number;
	likes_counts: number;
	play_counts?: number;
}

export interface NittosPreview {
	username: string;
	description: string;
	avatarUrl: string;
}
export interface Nittos {
	username: string;
	bio: string;
	avatarUrl: string;
	location: string;
	website: string;
	joinDate: string;
	stats: NittosStats;
	bannerUrl: string;
}

export interface NittosStats {
	tweets_counts: number;
	following_counts: number;
	followers_counts: number;
	likes_counts: number;
}
