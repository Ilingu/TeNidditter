<script lang="ts">
	import Art from "$lib/components/pages/auth/Art.svelte";
	import AuthForm from "$lib/components/pages/auth/AuthForm.svelte";
</script>

<section class="w-screen grid grid-cols-6">
	<AuthForm />
	<Art />
</section>

<style scoped>
	section {
		overflow: hidden;
		height: calc(100vh - 64px);
	}
</style>
