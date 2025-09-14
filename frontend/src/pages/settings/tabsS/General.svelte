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
    updateGamePath, // Добавляем updateGamePath для ручного обновления при первом запуске
  } from "../../../store/appSettings";
  import { get } from "svelte/store";
  import { Open } from "../../../../bindings/lce/backend/windows/settingswindow";
  import Radio from "./components/Radio.svelte"; // Импортируем компонент Radio
  import { t } from "svelte-i18n"; // Импортируем t для интернационализации текста
  import { tt } from "../../../lib/tooltip";
  import { Dialogs } from "@wailsio/runtime";

  let gamePathOptions = [];
  let selectedGamePath = ""; // Текущий выбранный путь к игре
  let isLoadingPaths = false; // Флаг для управления видимостью оверлея загрузки

  // Подписываемся на изменения appSettings для обновления списка путей
  appSettings.subscribe((settings) => {
    if (settings.all_paths && settings.all_paths.length > 0) {
      gamePathOptions = settings.all_paths.map((path) => ({
        label: path,
        value: path,
      }));

      // Приоритет отдается сохраненному game_path, если он существует и валиден.
      if (
        settings.game_path &&
        settings.all_paths.includes(settings.game_path)
      ) {
        selectedGamePath = settings.game_path;
      } else {
        // Если settings.game_path пуст, недействителен или отсутствует в currentAllPaths,
        // то никакой путь не должен быть выбран в UI.
        selectedGamePath = "";
      }
    } else {
      gamePathOptions = [];
      selectedGamePath = "";
    }
  });

  onMount(async () => {
    const isFirstRun = get(appSettings).first_run;

    if (isFirstRun) {
      await Open(); // Открываем окно настроек, если это первый запуск
      isLoadingPaths = true; // Показываем оверлей
      try {
        const pathsFound = await CheckAndFindPaths(); // Ждем прямое получение результатов
        console.log("Найденные пути:", pathsFound);

        await updateAllPaths(pathsFound); // Обновляем all_paths в настройках
        // Если найден ровно один путь, устанавливаем его как game_path при первом запуске
        if (pathsFound.length === 1) {
          await updateGamePath(pathsFound[0]);
        }
        await updateFirstRun(false); // Обновляем first_run
      } catch (error) {
        console.error("Ошибка при поиске или обновлении путей:", error);
        await updateFirstRun(false);
      } finally {
        isLoadingPaths = false; // Скрываем оверлей в любом случае
      }
    }
  });

  // Функция для обработки удаления пути
  async function handleDeletePath(pathToDelete) {
    const currentSettings = get(appSettings);
    const currentAllPaths = currentSettings.all_paths;
    const currentGamePath = currentSettings.game_path;

    // Фильтруем пути, удаляя выбранный
    const updatedPaths = currentAllPaths.filter(
      (path) => path !== pathToDelete
    );

    // Обновляем all_paths в хранилище
    await updateAllPaths(updatedPaths);

    // Если удаленный путь был текущим game_path, обновляем game_path
    if (pathToDelete === currentGamePath) {
      let newGamePath = "";
      if (updatedPaths.length > 0) {
        newGamePath = updatedPaths[0]; // Устанавливаем первый оставшийся путь
      }
      await updateGamePath(newGamePath); // Обновляем game_path в настройках
    }
  }

  // NEW: Функция для обработки нажатия на кнопку "Добавить путь"
  async function handleAddPathClick() {
    console.log("Открытие системного диалога для выбора папки.");
    try {
      const selectedDirectory = await Dialogs.OpenFile({
        Title: $t("select_game_folder_title"), // Заголовок диалога
        CanChooseDirectories: true, // Разрешаем выбор директорий
        CanChooseFiles: false, // Запрещаем выбор файлов
        AllowsMultipleSelection: false, // Разрешаем выбор только одной директории
      });

      if (selectedDirectory) {
        // OpenFile возвращает строку (путь) или пустую строку, если отменено
        const newPath = selectedDirectory;
        const currentSettings = get(appSettings);
        const currentAllPaths = currentSettings.all_paths;

        // Проверяем, существует ли уже такой путь
        if (!currentAllPaths.includes(newPath)) {
          const updatedPaths = [...currentAllPaths, newPath];
          await updateAllPaths(updatedPaths); // Добавляем новый путь в all_paths
          console.log("Новый путь добавлен:", newPath);

          // Если это был первый добавленный путь или текущий game_path пуст, устанавливаем его как game_path
          if (currentAllPaths.length === 0 || !currentSettings.game_path) {
            await updateGamePath(newPath);
          }
        } else {
          console.warn("Путь уже существует:", newPath);
          // TODO: Возможно, показать уведомление пользователю
        }
      } else {
        console.log("Выбор директории отменен.");
      }
    } catch (error) {
      console.error("Ошибка при открытии диалога или выборе папки:", error);
      // TODO: Возможно, показать уведомление об ошибке
    }
  }

  // Функция для запуска сканирования (для кнопки)
  async function handleRunScanner() {
    isLoadingPaths = true;
    try {
      const pathsFound = await CheckAndFindPaths();
      console.log("Найденные пути после сканирования:", pathsFound);

      await updateAllPaths(pathsFound);
      if (pathsFound.length === 1) {
        await updateGamePath(pathsFound[0]);
      }
      // При ручном сканировании, если это был первый запуск, отмечаем его как завершенный
      const isFirstRun = get(appSettings).first_run;
      if (isFirstRun) {
        await updateFirstRun(false);
      }
    } catch (error) {
      console.error("Ошибка при сканировании путей:", error);
      // Возможно, отобразить сообщение об ошибке пользователю
    } finally {
      isLoadingPaths = false;
    }
  }
</script>

{#if isLoadingPaths}
  <div class="loading-overlay">
    <div class="loading-content">
      <div class="spinner"></div>
      <span class="loading-text">{$t("scanning")}</span>
    </div>
  </div>
{/if}

{#if !isLoadingPaths}
  <div class="general-settings">
    {#if gamePathOptions.length > 0}
      <h3 class="choose-text">{$t("select_path")}</h3>
      <Radio
        options={gamePathOptions}
        name="game-path-selector"
        bind:selectedValue={selectedGamePath}
        onDelete={handleDeletePath}
      />
    {:else}
      <p class="not-found">
        {$t("paths_not_found")}
        <button class="run-scanner-button" on:click={handleRunScanner}>
          {$t("run_scanner")}
        </button>
      </p>
    {/if}

    <!-- Плавающая кнопка для добавления пути -->
    <button
      class="fab-add-button"
      on:click={handleAddPathClick}
      use:tt={{ content: $t("add_folder_tooltip") }}
    >
      +
    </button>
  </div>
{/if}

<style>
  .general-settings {
    padding: 20px;
    overflow-x: hidden;
    align-items: center;
  }

  .choose-text {
    text-align: center;
  }

  .not-found {
    text-align: center;
    position: absolute;
    top: 45%;
    left: 55%;
    transform: translate(-50%, -50%);
    display: flex; /* Используем flexbox для центрирования кнопки */
    flex-direction: column; /* Элементы располагаются вертикально */
    align-items: center; /* Центрирование по горизонтали */
    gap: 15px; /* Отступ между текстом и кнопкой */
  }

  /* Стили для оверлея загрузки */
  .loading-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.7); /* Затемненный фон */
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000; /* Поверх всех остальных элементов */
  }

  .loading-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    color: white;
    font-size: 1.2em;
  }

  .spinner {
    border: 4px solid rgba(255, 255, 255, 0.3);
    border-top: 4px solid #fff;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .loading-text {
    margin-top: 10px;
  }

  /* Стили для плавающей кнопки */
  .fab-add-button {
    position: absolute; /* Относительно .general-settings */
    bottom: 20px;
    right: 20px;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background-color: #3ba475; /* Цвет кнопки */
    color: white;
    font-size: 2em;
    border: none;
    box-shadow:
      0 6px 12px rgba(0, 0, 0, 0.4),
      0 0 0 3px rgba(0, 0, 0, 0.2); /* Более отчетливая тень */
    cursor: pointer;
    display: flex;
    justify-content: center;
    align-items: center;
    transition:
      background-color 0.3s ease,
      transform 0.3s ease,
      box-shadow 0.3s ease;
    z-index: 999; /* Чтобы была поверх контента */
  }

  .fab-add-button:hover {
    background-color: #2e8b57;
    transform: scale(1.1);
    box-shadow:
      0 8px 16px rgba(0, 0, 0, 0.5),
      0 0 0 4px rgba(255, 255, 255, 0.3); /* Увеличиваем тень при наведении */
  }

  /* Стили для кнопки "Run Scanner" */
  .run-scanner-button {
    padding: 10px 20px;
    background-color: #007bff; /* Примерный синий цвет */
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 1em;
    transition:
      background-color 0.3s ease,
      transform 0.2s ease;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .run-scanner-button:hover {
    background-color: #0056b3;
    transform: translateY(-2px);
  }
</style>
