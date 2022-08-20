<script lang="ts">
	import AlertProvider from "$lib/services/AlertProvider.svelte";
	import Navbar from "$lib/components/layout/Navbar.svelte";

	import "../style/app.css";
</script>

<Navbar />
<slot />
<AlertProvider />
