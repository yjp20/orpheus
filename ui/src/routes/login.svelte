<script>
	import { onMount } from "svelte";
	import { session } from "$app/stores";
	import Cookie from "js-cookie"

	onMount(() => {
		const params = new URLSearchParams(window.location.hash);
		const token = params.get("access_token")
		const expires_in = params.get("expires_in")
		const expires_at = Date.now() + expires_in * 1000

		Cookie.set("access_token", token)
		Cookie.set("expires_at", expires_at)

		$session.access_token = token
		$session.expires_at = expires_at
	})
</script>
