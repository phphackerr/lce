//@ts-nocheck
import "./lib/icons";
import "./lib/tooltip";
import App from "./App.svelte";
import { initGoI18n } from "./store/i18n";
import { loadSettings, appSettings } from "./store/appSettings";
import { applyTheme } from "./lib/theming";
import { get } from "svelte/store";

async function initialiseApp() {
  // Создаем асинхронную функцию
  await initGoI18n(); // Ждем инициализации i18n
  await loadSettings().then(() => {
    const theme = get(appSettings).theme; // ИСПРАВЛЕНО
    applyTheme(theme);
  }); // Загружаем настройки

  const app = new App({
    target: document.getElementById("app"),
  });

  return app;
}

const app = initialiseApp(); // Вызываем асинхронную функцию
export default app;
