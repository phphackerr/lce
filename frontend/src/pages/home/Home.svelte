<script>
  import TitleBar from "./TitleBar.svelte";
  import Watcher from "./Watcher.svelte";
</script>

<div class="main">
  <TitleBar />
  <Watcher />
</div>

<style>
  .main {
    width: 100vw;
    height: 100vh;
    background: var(--background-color);
    color: var(--text-color);
  }
</style>
