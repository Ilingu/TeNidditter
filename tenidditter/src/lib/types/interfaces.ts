export interface FunctionJob<T = never> {
	success: boolean;
	data?: T;
}
