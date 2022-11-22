<script>
	import Link from "$lib/client/components/design/Link.svelte";
	import { page } from "$app/stores";
	import "highlight.js/styles/atom-one-dark.css"; // theme for the codeblocks
</script>

<div class="md:hidden dropdown w-screen mt-2">
	<div class="grid place-items-center">
		<label tabindex="0" for="" class="text-xl gap-x-2">
			<div class="form-control">
				<div class="input-group">
					<input type="text" placeholder="Search teddit" readonly class="input input-bordered" />
					<button class="btn btn-square">
						<i class="fa-solid fa-magnifying-glass" />
					</button>
				</div>
			</div>
		</label>
	</div>

	<ul
		tabindex="0"
		class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-72"
	>
		<li>
			<Link href={"/teddit/r"}
				><button class="btn btn-primary gap-x-3 w-full"
					><i class="fa-solid fa-magnifying-glass" /> Search Subteddit</button
				></Link
			>
		</li>
		<li>
			<Link href={"/teddit/u"}
				><button class="btn btn-warning gap-x-3 w-full">
					<i class="fas fa-user" />
					Search Teddit's User</button
				></Link
			>
		</li>
	</ul>
</div>
{#if $page.params.subteddit}
	<div class="text-sm breadcrumbs ml-2">
		<ul>
			<li><Link href="/teddit">Feed</Link></li>
			<li><Link href={`/teddit/r/${$page.params.subteddit}`}>{$page.params.subteddit}</Link></li>
			{#if $page.params.postid}
				<li>
					<Link href={`/teddit/r/${$page.params.subteddit}/${$page.params.postid}`}
						>{$page.params.postid}</Link
					>
				</li>
			{/if}
		</ul>
	</div>
{/if}
{#if $page.params.user}
	<div class="text-sm breadcrumbs ml-2">
		<ul>
			<li><Link href="/teddit">Feed</Link></li>
			<li><Link href={`/teddit/u/${$page.params.user}`}>{$page.params.user}</Link></li>
		</ul>
	</div>
{/if}

<slot />
