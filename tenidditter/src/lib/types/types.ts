export type Themes = "tenidditter" | "nitter" | "teddit";
export type AlertTypes = "success" | "info" | "warning" | "error";
export type FeedType = "user_feed" | "home_feed";
export type FeedHomeType = "hot" | "new" | "top" | "rising" | "controversial";

export type Tuple<TItem, TLength extends number> = [TItem, ...TItem[]] & { length: TLength };
