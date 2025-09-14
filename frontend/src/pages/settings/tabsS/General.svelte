<script context="module">
  export const tabMetadata = {
    order: 1,
  };
</script>

<script>
  import { onMount } from "svelte";
  import { CheckAndFindPaths } from "../../../../bindings/lce/backend/modules/paths_scanner/scanner";
  import {
    appSettings,
    updateFirstRun,
    updateAllPaths,
  } from "../../../store/appSettings";
  import { get } from "svelte/store";
  import { Open } from "../../../../bindings/lce/backend/windows/settingswindow";

  onMount(async () => {
    const isFirstRun = get(appSettings).first_run;

    if (isFirstRun) {
      await Open();
      try {
        const pathsFound = await CheckAndFindPaths(); // Ждем прямое получение результатов
        console.log("Найденные пути:", pathsFound);

        await updateAllPaths(pathsFound); // Обновляем all_paths в настройках
        await updateFirstRun(false); // Обновляем first_run
      } catch (error) {
        console.error("Ошибка при поиске или обновлении путей:", error);
        // Обработка ошибки, возможно, установить first_run в false все равно,
        // или предоставить пользователю возможность повторить
        await updateFirstRun(false);
      }
    }
  });
</script>

<style></style>
