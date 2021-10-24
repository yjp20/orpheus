<script context="module">
	export async function load({page}) {
		return {
			props: {
				server_id: page.params.id
			}
		}
	}
</script>

<script>
	export let server_id;

	let channel = "main";

	let server = {
		name: "future tft challenger",
		channels: [ "" ]
	}

	let queue_items = [
		{ name: "Hello", artist: "Adele", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
		{ name: "Fly me to the moon", artist: "Frank Sinatra", queued_by: "yourdoge#4501", url: "https://google.com", length: 301 },
	];
	let playing_index = 0;
	let allow_edit = true;
	function formatDuration(seconds) {
		const s = ~~(seconds % 60)
		const m = ~~(seconds / 60)
		return `${m}:${s.toString().padStart(2, '0')}`
	}
	$: playing = queue_items[playing_index];
</script>

<div class="section is-small">
	<div class="container max-width-desktop">
		<div class="columns">
			<div class="column is-4 player-side-container">
				<div class="player-side">
					<form >
						<input placeholder="Search..." />
					</form>
					<div class="card">
						<div class="card-header">
							Currently Playing
						</div>
						<div class="card-item is-vertical">
							<img class="playing-image" src="https://bulma.io/images/placeholders/256x256.png" />
							<div class="playing-name subtitle">{playing.name}</div>
							<div class="playing-artist paragraph">Artist: {playing.artist} </div>
							<div class="playing-queued paragraph">Queued by: {playing.queued_by}</div>
							<div class="player">
								<div class="player-progress">
									<div class="player-progress-bar" style="width: 50%"></div>
								</div>
								<div class="player-buttons">
									<button class="player-left button">prev</button>
									<button class="player-toggle button">toggle</button>
									<button class="player-right button">next</button>
								</div>
							</div>
						</div>
					</div>
					<p class="paragraph">
						<a href="/servers/{server_id}/settings">Playing on channel '{channel}'</a>
					</p>
				</div>
			</div>
			<div class="column is-8">
				<div class="card" class:is-selectable={allow_edit}>
					<div class="card-header">
						Queue
					</div>
					{#each queue_items as item}
						<div class="song card-item is-vertical" class:is-playing={item == playing}>
							<div class="song-header">
								<div class="song-title"> {item.name} - {item.artist} </div>
								<div class="song-length"> {formatDuration(item.length)} </div>
							</div>
							<div class="song-footer">
								<a class="song-url label is-flat" href={item.url}>{item.url}</a>
								<div class="song-queued label is-flat">{item.queued_by}</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.player-side-container {
		position: relative;
	}

	.player-side {
		position: sticky;
		top: 1rem;
	}

	.playing-image {
		width: 100%;
		display: block;
	}

	.playing-name {
		margin-top: 1rem!important;
		margin-bottom: 0.5rem!important;
	}

	.playing-artist,
	.playing-queued {
		margin-top: 0.5rem!important;
		margin-bottom: 0.5rem!important;
	}

	.song.is-playing {
		background-color: var(--blue-light);
	}

	.song-header {
		display: flex;
	}

	.song-title {
		font-weight: bold;
	}

	.song-length {
		margin-left: auto;
	}

	.song-footer {
		display: flex;
	}

	.song-queued {
		margin-left: auto;
	}

	.player {
		margin-top: 1rem;
	}

	.player-progress {
		height: 5px;
		background-color: var(--greylight);
	}

	.player-progress-bar {
		height: 100%;
		background-color: var(--blue);
	}

	.player-buttons {
		display: flex;
		margin-top: 0.5rem;
		justify-content: center;
	}

	.player-buttons button + button {
		margin-left: 0.25rem;
	}
</style>
