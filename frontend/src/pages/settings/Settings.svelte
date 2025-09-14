<script>
  import TabsS from "./TabsS.svelte";
  import TitleBarS from "./TitleBarS.svelte";
</script>

<div class="main">
  <TitleBarS />
  <div class="tabs">
    <TabsS />
  </div>
</div>

<style>
  .main {
    width: 100vw;
    height: 100vh;
    margin: 0;
    padding: 0;
    color: var(--settings-text-color);
    background: var(--settings-bg-color);
    overflow: hidden;
  }

  .tabs {
    margin-top: 30px;
    height: calc(100vh - 30px);
    box-sizing: border-box;
    overflow: hidden;
  }
</style>
