<script>
	import Link from "$lib/components/design/Link.svelte";
</script>

<header class="hero min-h-[70vh]  ">
	<div class="hero-overlay">
		<div class="flex justify-center w-full h-full">
			<div class="flex flex-col justify-center items-center w-3/4">
				<div class="w-1/2 h-fit relative flex flex-col items-center">
					<div class="badge badge-primary badge-outline absolute h-8 -top-10 -left-20">
						<div class="badge badge-accent mr-2">@alpha</div>
						<a
							href="https://github.com/Ilingu/TeNidditter"
							target="_blank"
							rel="noopener noreferrer"
						>
							Github <i class="fa-brands fa-github icon" /></a
						>
					</div>

					<h1
						id="HeadlineHero"
						class="scrollAnimate light mb-5 text-5xl font-fancy text-primary tracking-wider leading-snug"
					>
						Use <span class="text-nitter font-bold">Twitter</span> and
						<span class="text-teddit font-bold">Reddit</span>
						without using them.
					</h1>
					<p class="mb-5">
						TeNidditter allow you to browse your twitter and reddit feed without giving them <strong
							>all your data.</strong
						>
					</p>

					<Link
						href="/auth"
						classStyle="btn btn-primary transition-all w-1/2 flex gap-x-3 text-lg hover:bg-gradient-to-r hover:to-teddit hover:from-nitter hover:text-white"
					>
						<i class="fa-solid fa-fire text-orange-600" /> Get Started</Link
					>
				</div>
			</div>
		</div>
	</div>
</header>
