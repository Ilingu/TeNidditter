<script lang="ts">
	import Link from "$lib/components/design/Link.svelte";
	import ProfilePicture from "$lib/components/layout/ProfilePicture.svelte";
	import Feeds from "$lib/components/pages/nitter/Feeds.svelte";
	import AuthStore from "$lib/stores/auth";

	export let data: import("./$types").PageServerData;
</script>

<main class="grid place-items-center mt-5">
	<div class="max-w-[1250px] flex md:flex-row flex-col gap-y-4 justify-center gap-x-10 z-0">
		<div class="drawer rounded-lg h-fit md:sticky top-20">
			<div class="drawer-side flex justify-center">
				<ul class="menu p-4 md:w-80 w-3/4 rounded bg-base-300 text-base-content">
					<!-- Sidebar content here -->
					<li class="flex-row items-center text-accent font-bold text-2xl">
						<i class="fa-solid fa-rss" />
						{$AuthStore.user?.username}'s Feed
					</li>
					<li>
						<Link href="/nitter/search">
							<i class="fa-solid fa-magnifying-glass" />
							Search Nitter</Link
						>
					</li>
					<li>
						<Link href="/nitter/lists"><i class="fa-solid fa-list" /> Lists</Link>
					</li>
					<li class="flex-row items-center text-accent font-bold text-2xl">
						<i class="fa-solid fa-user" /> Subs
					</li>
					<li>
						{#each ($AuthStore?.Subs?.nitter ?? []).slice(0, 20) as nittos}
							<Link href={`/nitter/${nittos}`}>
								<ProfilePicture size="small" />
								<span>{nittos}</span>
							</Link>
						{/each}
					</li>
				</ul>
			</div>
		</div>
		<div class="max-w-[750px]">
			<Feeds neets={data.comments} />
		</div>
	</div>
</main>
