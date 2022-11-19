<script lang="ts">
	import { page } from "$app/stores";
	import Link from "$lib/components/design/Link.svelte";
</script>

{#if $page.params.nittos}
	<div class="text-sm breadcrumbs ml-2">
		<ul>
			<li><Link href="/nitter">Feed</Link></li>
			<li><Link href={`/nitter/${$page.params.nittos}`}>{$page.params.nittos}</Link></li>
			{#if $page.params.neet}
				<li>
					<Link href={`/nitter/${$page.params.nittos}/${$page.params.neet}`}
						>{$page.params.neet}</Link
					>
				</li>
			{/if}
		</ul>
	</div>
{/if}
<slot />
