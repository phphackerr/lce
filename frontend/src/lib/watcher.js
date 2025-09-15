import { Events } from "@wailsio/runtime";
import {
  StartWatching,
  StopWatching,
} from "../../bindings/lce/backend/modules/config_watcher/configwatcher";
import { appSettings } from "../store/appSettings";

let currentPath = null;

// --- API для UI --- //
export function onConfigChanged(callback) {
  Events.On("config-changed", (event) => {
    const filePath = Array.isArray(event.data) ? event.data[0] : event.data;
    console.log("⚡ Config changed:", filePath);
    callback(filePath);
  });
}

export async function startWatcher(path) {
  if (!path) {
    console.warn("⚠ startWatcher: пустой путь");
    return;
  }

  // Останавливаем предыдущий вотчер, если был
  if (currentPath) {
    try {
      await StopWatching();
      console.log("🛑 Previous watcher stopped");
    } catch (err) {
      console.warn("⚠ Ошибка при остановке предыдущего вотчера:", err);
    } finally {
      currentPath = null;
    }
  }

  try {
    await StartWatching(path);
    currentPath = path;
    console.log("👀 Watcher started for:", path);
  } catch (err) {
    console.error("Ошибка запуска вотчера:", err);
  }
}

export async function stopWatcher() {
  try {
    await StopWatching();
    console.log("🛑 Watcher stopped");
  } catch (err) {
    console.error("Ошибка остановки вотчера:", err);
  } finally {
    currentPath = null;
  }
}

// --- Автоподписка на смену game_path --- //
appSettings.subscribe((settings) => {
  // Укажем именно файл, а не папку
  const newPath = settings.game_path
    ? settings.game_path + "\\config.lod.ini"
    : null;

  if (newPath !== currentPath) {
    if (newPath) {
      startWatcher(newPath);
    } else if (currentPath) {
      stopWatcher();
    }
  }
});
