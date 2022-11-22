<script>
	import Link from "$lib/client/components/design/Link.svelte";
</script>

<header class="hero lg:min-h-[70vh] min-h-[40vh]">
	<div class="hero-overlay">
		<div class="flex justify-center w-full h-full">
			<div class="flex flex-col justify-center items-center 2xl:w-3/4 w-full">
				<div class="xl:w-1/2 w-full h-fit pt-10 flex flex-col items-center">
					<h1
						id="HeadlineHero"
						class="scrollAnimate md:p-0 px-2 relative light mb-5 xl:text-5xl md:text-4xl text-3xl md:text-left text-center  font-fancy text-primary tracking-wider leading-snug"
					>
						<div class="badge badge-primary badge-outline absolute h-8 -top-10 lg:-left-20 left-0">
							<div class="badge badge-accent mr-2">@alpha</div>
							<a
								href="https://github.com/Ilingu/TeNidditter"
								target="_blank"
								rel="noopener noreferrer"
							>
								Github <i class="fa-brands fa-github icon" /></a
							>
						</div>
						<p>
							Use <span class="text-nitter font-bold">Twitter</span> and
							<span class="text-teddit font-bold">Reddit</span>
							without using them.
						</p>
					</h1>
					<p class="mb-5 lg:p-0 lg:text-left text-justify px-10">
						TeNidditter allow you to browse your twitter and reddit feed without giving them <strong
							>all your data.</strong
						>
					</p>

					<button
						class="md:w-1/2 w-full lg:scale-100 md:scale-90 scale-75 btn btn-primary transition-all flex gap-x-3 text-lg hover:bg-gradient-to-r hover:to-teddit hover:from-nitter hover:text-white"
					>
						<Link href="/auth">
							<div class="w-[360px]">
								<i class="fa-solid fa-fire text-orange-600" /> Get Started
							</div>
						</Link>
					</button>
				</div>
			</div>
		</div>
	</div>
</header>
