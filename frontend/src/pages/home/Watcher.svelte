<script>
  import { onMount, onDestroy } from "svelte";
  import { appSettings } from "../../store/appSettings";
  import {
    startWatcher,
    stopWatcher,
    onConfigChanged,
  } from "../../lib/watcher";
  import { isInternalChange } from "../../store/internalChange";
  import {
    LoadConfig,
    SetConfigValue,
    CheckConfigDiff,
  } from "../../../bindings/lce/backend/modules/config_editor/configeditor";
  import { get } from "svelte/store";
  import { t } from "svelte-i18n";

  let diff = {};
  let currentPath = null;
  let changedFile = "";
  let unsubscribe;
  let configListener;
  let isHovered = false;

  onMount(() => {
    unsubscribe = appSettings.subscribe((settings) => {
      const newPath = settings.game_path
        ? settings.game_path + "\\config.lod.ini"
        : null;
      if (newPath !== currentPath) {
        if (newPath) LoadConfig().then(() => startWatcher(newPath));
        else if (currentPath) stopWatcher();
        currentPath = newPath;
      }
    });

    configListener = onConfigChanged(async (filePath) => {
      if (!get(isInternalChange)) {
        changedFile = filePath;
        diff = await CheckConfigDiff();
      }
    });
  });

  onDestroy(() => {
    if (unsubscribe) unsubscribe();
    if (configListener) configListener.off();
  });

  function acceptChanges() {
    isInternalChange.mark();
    SetConfigValue("GAMEOPTIONS", "WideScreen", "false");
    diff = {};
    changedFile = "";
  }

  function restoreChanges() {
    diff = {};
    changedFile = "";
  }

  function handleMouseEnter() {
    isHovered = true;
  }

  function handleMouseLeave() {
    isHovered = false;
  }
</script>

<div class="floating-panel">
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    class="panel"
    on:mouseenter={handleMouseEnter}
    on:mouseleave={handleMouseLeave}
    style="
      height: {isHovered
      ? Object.keys(diff).length
        ? '300px'
        : '60px'
      : '60px'};
      width: {isHovered
      ? Object.keys(diff).length
        ? '400px'
        : '250px'
      : '60px'};
      border-radius: {isHovered ? '5px' : '50%'};"
  >
    <div
      class="indicator"
      style="
      background-color: {Object.keys(diff).length ? 'orange' : 'green'};
      height: {Object.keys(diff).length
        ? isHovered
          ? '93%'
          : '40px'
        : '40px'};
      border-radius: {Object.keys(diff).length
        ? isHovered
          ? '5px'
          : '50%'
        : '50%'};"
    >
      <div class="icon">
        {Object.keys(diff).length ? "⚠" : "✔"}
      </div>
    </div>
    <div class="content" class:visible={isHovered}>
      {#if Object.keys(diff).length}
        <div class="headline">{$t("external_changes")}:</div>
        <hr />
        {#each Object.entries(diff) as [section, keys]}
          <div class="section">
            <strong>{section}</strong>
            <ul>
              {#each Object.entries(keys) as [key, info]}
                <li>
                  {key}: {info.old ?? "<deleted>"} → {info.new ?? "<deleted>"} ({info.status})
                </li>
              {/each}
            </ul>
          </div>
        {/each}
      {:else}
        Im working, all fine
      {/if}
    </div>
  </div>
</div>

<style>
  .floating-panel {
    height: fit-content;
    max-height: 300px;
    max-width: 400px;
    padding: 10px;
    border-radius: 5px;
    position: fixed;
    bottom: 5px;
    left: 5px;
    z-index: 999;
  }

  .panel {
    width: 60px;
    height: 60px;
    border-radius: 5px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
    display: flex;
    background-color: rgba(0, 0, 0, 0.5);
    flex-direction: row;
    gap: 10px;
    transition: all 0.3s ease;
  }

  .indicator {
    width: 40px;
    height: 40px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    align-self: flex-end;
    border-radius: 50%;
    margin-left: 10px;
    margin-bottom: 10px;
    transition: all 0.3s ease;
  }

  .icon {
    font-size: 22px;
  }

  .content {
    opacity: 0;
    visibility: hidden;
    max-height: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    align-items: center;
    flex-grow: 1;
    transition:
      opacity 0.3s ease,
      visibility 0.3s ease,
      max-height 0.3s ease,
      border-color 0.3s ease;
  }

  .content.visible {
    opacity: 1;
    visibility: visible;
    max-height: 100%;
    padding: 5px;
    overflow-y: auto; /* <--- ДОБАВИТЬ ЭТУ СТРОКУ */
  }

  .content hr {
    width: 90%;
  }

  .section {
    margin-bottom: 10px;
  }

  .section strong {
    display: block;
    margin-bottom: 3px;
  }

  .section ul {
    list-style: none;
    padding-left: 10px;
    margin: 0;
  }

  .section li {
    font-size: 0.85em;
  }
</style>
