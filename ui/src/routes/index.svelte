<script context="module">

</script>

<script>
	let servers = [
		{ name: "future tft challengers", id : "1234" },
		{ name: "future league challengers", id : "12342314" },
	]
</script>

<div class="section is-small">
	<div class="container max-width-desktop">
		<div class="discord-box box is-vertical">
			<h2 class="subtitle">Connect to discord</h2>
			<p class="paragraph">Login to see which servers are running on this instance of Orpheus</p>
			<button class="button is-white">Connect</button>
		</div>
		<div class="github-box paragraph box is-vertical">
			<h2 class="subtitle">Support our development</h2>
			<p class="paragraph">Contribute your time or money at our github to keep Orpheus developed</p>
			<button class="button is-white">Connect</button>
		</div>
		<div class="card is-selectable">
			<div class="card-header">
				Available servers
			</div>
			{#each servers as server}
				<a class="no-link card-item" href="/servers/{server.id}">
					<div class="server-name">
						{server.name}
					</div>
				</a>
			{/each}
		</div>
	</div>
</div>

<style>
	.discord-box {
		background-color: #5865F2;
		color: white;
	}

	.github-box {
		background-color: var(--dark);
		color: white;
	}
</style>
