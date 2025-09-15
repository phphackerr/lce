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

  let changedFile = "";
  let currentPath = null;
  let unsubscribe;
  let configListener;

  onMount(() => {
    // --- Подписка на изменения game_path ---
    unsubscribe = appSettings.subscribe((settings) => {
      const newPath = settings.game_path
        ? settings.game_path + "\\config.lod.ini"
        : null;

      if (newPath !== currentPath) {
        if (newPath) {
          LoadConfig()
            .then(() => startWatcher(newPath))
            .catch((err) => console.error("Ошибка загрузки конфига:", err));
        } else if (currentPath) {
          stopWatcher();
        }
        currentPath = newPath;
      }
    });

    // --- Подписка на события изменения конфига ---
    configListener = onConfigChanged(async (filePath) => {
      if (!get(isInternalChange)) {
        changedFile = filePath;
        const diff = await CheckConfigDiff();
        console.log("🔍 External changes:", diff);
      }
    });
  });

  // --- Отписка при unmount ---
  onDestroy(() => {
    if (unsubscribe) unsubscribe();
    if (configListener) configListener.off();
  });

  function reloadConfig() {
    isInternalChange.mark();
    SetConfigValue("GAMEOPTIONS", "WideScreen", "false");
  }

  function ignoreChange() {}
</script>

<div class="config-alert">
  <p>Конфиг изменён: {changedFile}</p>
  <button on:click={reloadConfig}>Перезагрузить</button>
  <button on:click={ignoreChange}>Игнорировать</button>
</div>

<style>
  .config-alert {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: #333;
    color: #fff;
    padding: 15px;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .config-alert button {
    padding: 5px 10px;
    cursor: pointer;
  }
</style>
