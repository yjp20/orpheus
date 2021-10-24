<script context="module">
  export async function load({ page, fetch }) {
    const v = await fetch(
      `http://localhost:4000/api/queue?guild_id=${page.params.id}`
    );
    const o = await v.json();
    console.log(o);
    return {
      props: {
        server: o,
        server_id: page.params.id,
      },
    };
  }
</script>

<script>
  export let server_id;

  let channel = "General";

  /* let server = { */
  /* 	name: "calhacks orpheus", */
  /* 	channels: [ "General" ] */
  /* } */

  export let server;
  let playing_index = 0;
  let search = "";
  let allow_edit = true;

  function formatDuration(nano) {
		seconds = nano / 1000 / 1000 / 1000
    const s = ~~(seconds % 60);
    const m = ~~(seconds / 60);
    return `${m}:${s.toString().padStart(2, "0")}`;
  }

  $: playing = server.queue[playing_index];

  async function submit(e) {
    e.preventDefault();
    await fetch("http://localhost:4000/api/queue", {
      method: "post",
      body: JSON.stringify({
        guild_id: server_id,
        user_id: "asdfwe",
        url: search,
      }),
    });
    search = "";
  }
</script>

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
    margin-top: 1rem !important;
    margin-bottom: 0.5rem !important;
  }

  .playing-artist,
  .playing-queued {
    margin-top: 0.5rem !important;
    margin-bottom: 0.5rem !important;
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

<div class="section is-small">
  <div class="container max-width-desktop">
    <div class="columns">
      <div class="column is-4 player-side-container">
        <div class="player-side">
          <form on:submit={submit}>
            <input name="search" bind:value={search} placeholder="Play..." />
          </form>
          {#if playing}
            <div class="card">
              <div class="card-header">Currently Playing</div>
              <div class="card-item is-vertical">
                <img
                  class="playing-image"
                  src="https://bulma.io/images/placeholders/256x256.png" />
                <div class="playing-name subtitle">{playing.song.name}</div>
                <div class="playing-artist paragraph">
                  Artist: {playing.artist}
                </div>
                <div class="playing-queued paragraph">
                  Queued by: {playing.queued_by}
                </div>
                <div class="player">
                  <div class="player-progress">
                    <div class="player-progress-bar" style="width: 50%" />
                  </div>
                  <div class="player-buttons">
                    <button class="player-left button">prev</button>
                    <button class="player-toggle button">toggle</button>
                    <button class="player-right button">next</button>
                  </div>
                </div>
              </div>
            </div>
          {/if}
          <p class="paragraph">
            <a href="/servers/{server_id}/settings">
              Playing on channel '{channel}'
            </a>
          </p>
        </div>
      </div>
      <div class="column is-8">
        <div class="card" class:is-selectable={allow_edit}>
          <div class="card-header">Queue</div>
          {#each server.queue as item}
            <div
              class="song card-item is-vertical"
              class:is-playing={item == playing}>
              <div class="song-header">
                <div class="song-title">{item.song.name}</div>
                <div class="song-length">
                  {formatDuration(item.song.length)}
                </div>
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
