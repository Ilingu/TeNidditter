<script lang="ts">
	import { page } from "$app/stores";
	import { setContext } from "svelte";
	import Comments from "$lib/client/components/pages/teddit/comments/Comments.svelte";
	import Link from "$lib/client/components/design/Link.svelte";
	import Feeds from "$lib/client/components/pages/teddit/Feeds.svelte";
	import { FormatElapsedTime } from "$lib/client/ClientUtils";

	export let data: import("./$types").PageData;

	setContext("COMMEMTS_CTX", data.comments); // Context because comments will be used to render recursive childrens (with undefined depth, so props thrilling can hit hard)

	// convert the post number of ups to readable and exploitable number
	const PostUps = (): number => {
		const ups = data.post_datas.post_ups;
		if (ups.endsWith("k")) {
			const [num] = ups.split("k");
			const safeNum = Math.round(parseFloat(num));
			return parseInt(`${safeNum}000`);
		}
		return parseInt(ups);
	};

	// generates an custom og image for this post
	const OGImageUrl = `${$page.url.origin}/api/og-image?title=${encodeURIComponent(
		data.post_datas.post_title
	)}&author=${encodeURIComponent(data.post_datas.post_author)}&subreddit=${encodeURIComponent(
		$page.params.subteddit
	)}&ups=${PostUps()}&created=${data.post_datas.post_created}`;
	const OGDesc = `Submitted ${FormatElapsedTime(
		data.post_datas.post_created * 1000,
		Date.now()
	)} by ${data.post_datas.post_author} on r/${$page.params.subteddit}`;
</script>

<svelte:head>
	<!-- Twitter Card -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:url" content={$page.url.origin} />
	<meta name="twitter:title" content={data.post_datas.post_title} />
	<meta name="twitter:description" content={OGDesc} />
	<meta name="twitter:image" content={OGImageUrl} />
	<meta name="twitter:creator" content={"u/" + data.post_datas.post_author} />
	<!-- OG Card -->
	<meta property="og:type" content="website" />
	<meta property="og:title" content={data.post_datas.post_title} />
	<meta property="og:description" content={OGDesc} />
	<meta property="og:site_name" content="TeNidditter" />
	<meta property="og:url" content={$page.url.origin} />
	<meta property="og:image" content={OGImageUrl} />
</svelte:head>

<main class="max-w-[1500px] m-auto flex flex-col items-center gap-x-8 px-2 py-5">
	<div class="max-w-[750px]">
		<Feeds
			rawPosts={[
				{
					id: $page.params.postid,
					title: data.post_datas.post_title,
					author: data.post_datas.post_author,
					created: data.post_datas.post_created,
					body_html: data.post_datas.body_html,

					ups: PostUps(),
					num_comments: data.post_datas.post_nb_comments,

					subreddit: $page.params.subteddit,

					is_self_link: true,
					is_video: false,
					stickied: false,

					url: "",
					domain: "",
					permalink: "",

					blur: false
				}
			]}
		/>

		<details class="w-full max-w-xs mt-5">
			<summary>
				Sorted by <span class="text-white font-bold"
					>{$page.url.searchParams.get("sort") ?? "best"}</span
				>
			</summary>
			<ul style="list-style: inside;">
				<li><Link href={$page.url.href.split("?sort=")[0] + `?sort=best`}>Best</Link></li>
				<li><Link href={$page.url.href.split("?sort=")[0] + `?sort=top`}>Top</Link></li>
				<li><Link href={$page.url.href.split("?sort=")[0] + `?sort=new`}>New</Link></li>
				<li>
					<Link href={$page.url.href.split("?sort=")[0] + `?sort=controversial`}>Controversial</Link
					>
				</li>
				<li><Link href={$page.url.href.split("?sort=")[0] + `?sort=old`}>Old</Link></li>
				<li><Link href={$page.url.href.split("?sort=")[0] + `?sort=qa`}>Q&A</Link></li>
			</ul>
		</details>

		<section class="mt-2">
			{#each data.comments as co, i}
				<Comments idxCtx={i} comment={co[0]} open recursive />
				<div class="my-10" />
			{/each}
		</section>
	</div>
</main>
