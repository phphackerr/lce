//@ts-nocheck

import { writable } from "svelte/store";
import {
  GetSettings,
  UpdateSettings,
  GetOption,
} from "../../bindings/lce/backend/modules/app_settings/appsettings";

// Создаем store для настроек
export const appSettings = writable({
  width: 1600,
  height: 900,
  language: "en",
  game_path: "",
  first_run: true, // NEW
  all_paths: [], // NEW
  theme: "default", // NEW
});

// Функция для загрузки настроек
export async function loadSettings() {
  try {
    const settings = await GetSettings();
    appSettings.set(settings);
    console.log("Настройки загружены:", settings);
    return settings;
  } catch (error) {
    console.error("Ошибка загрузки настроек:", error);
    // Используем значения по умолчанию при ошибке
    const defaultSettings = {
      width: 1600,
      height: 900,
      language: "en",
      game_path: "",
      first_run: true, // NEW
      all_paths: [], // NEW
      theme: "default", // NEW
    };
    appSettings.set(defaultSettings);
    return defaultSettings;
  }
}

// Функция для обновления настроек
export async function updateSettings(newSettings) {
  try {
    appSettings.update((current) => ({ ...current, ...newSettings }));

    const updatedSettings = await UpdateSettings(newSettings);

    appSettings.set(updatedSettings);
    console.log("Настройки обновлены:", updatedSettings);
    return updatedSettings;
  } catch (error) {
    console.error("Ошибка обновления настроек:", error);
    await loadSettings();
    throw error;
  }
}

// Получение конкретной опции
export async function getSetting(key) {
  try {
    const value = await GetOption(key);
    return value;
  } catch (error) {
    console.error(`Ошибка получения настройки ${key}:`, error);
    return null;
  }
}

// Сброс настроек
export async function resetSettings() {
  try {
    const defaultSettings = {
      width: 1600,
      height: 900,
      language: "en",
      game_path: "",
      first_run: true, // NEW
      all_paths: [], // NEW
      theme: "default", // NEW
    };

    const updatedSettings = await UpdateSettings(defaultSettings);

    appSettings.set(updatedSettings);
    console.log("Настройки сброшены к значениям по умолчанию");
    return updatedSettings;
  } catch (error) {
    console.error("Ошибка сброса настроек:", error);
    throw error;
  }
}

// Хелперы
export async function updateWindowSize(width, height) {
  return updateSettings({ width, height });
}

export async function updateLanguage(language) {
  console.log(`Язык изменен на: ${language}`);
  return updateSettings({ language });
}

export async function updateGamePath(gamePath) {
  console.log(`Путь к игре изменен на: ${gamePath}`);
  return updateSettings({ game_path: gamePath });
}

export async function updateFirstRun(firstRun) {
  console.log(`First run установлен в: ${firstRun}`);
  return updateSettings({ first_run: firstRun });
}

export async function updateAllPaths(allPaths) {
  console.log(`Все пути обновлены на: ${allPaths}`);
  return updateSettings({ all_paths: allPaths });
}

export async function updateTheme(theme) {
  console.log(`Тема изменена на: ${theme}`);
  return updateSettings({ theme: theme });
}
