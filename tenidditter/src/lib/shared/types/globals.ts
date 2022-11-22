export interface FunctionJob<T = never> {
	success: boolean;
	data?: T;
	error?: string;
}

export type Tuple<TItem, TLength extends number> = [TItem, ...TItem[]] & { length: TLength };
