<script context="module">
	import "briefscss";
	let client_id = "901438105679187998";
</script>

<script>
	import { onMount } from "svelte";
	import { session } from "$app/stores"

	export let hostname;

	onMount(() => {
		hostname = window.location.origin
	});
</script>

<svelte:head>
	<title>Orpheus</title>
</svelte:head>

<div class="navbar">
	<div class="navbar-container container max-width-desktop">
		<a class="navbar-title" href="/">Orpheus</a>
		<div class="navbar-group">
			{#if $session.access_token}
				<a class="navbar-item" href="/">Servers</a>
				<div class="navbar-item">Logged in as <b>{$session.access_token}</b></div>
			{:else}
				<a class="navbar-item" href="https://discord.com/api/oauth2/authorize?response_type=token&scope=identify%20guilds&client_id={client_id}&scope=identify&redirect_uri={encodeURIComponent(hostname + "/login")}">Connect to Discord</a>
			{/if}
		</div>
	</div>
</div>

<slot></slot>

<style>
	:global(body) {
		--border: var(--grey);
		background-color: var(--light);;
	}

	.navbar {
		padding: 1em 0;
		background-color: var(--white);
		border-bottom: 1px solid var(--border);
	}

	.navbar-container {
		display: flex;
	}

	.navbar-title {
		font-weight: bold;
		font-size: 1.25em;
	}

	.navbar-group {
		display: flex;
		margin-left: auto;
		align-items: center;
	}

	.navbar-item {
		display: block;
	}

	.navbar-item + .navbar-item{
		margin-left: 1rem;
	}
</style>
