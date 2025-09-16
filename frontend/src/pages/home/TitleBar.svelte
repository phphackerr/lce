<script>
  //@ts-nocheck
  import { Window, Application } from "@wailsio/runtime";
  import { Open } from "../../../bindings/lce/backend/windows/settingswindow";
  import icon from "/favicon.png";
  import {
    SettingsIc,
    MinimizeIc,
    MaximizeIc,
    CloseIc,
  } from "../../lib/icons.js";
  import { tt } from "../../lib/tooltip";
  import { t } from "svelte-i18n";
</script>

<div class="titlebar" on:dblclick={Window.ToggleMaximise} role="presentation">
  <div class="logo">
    <img src={icon} class="icon" alt="icon" />
    <span class="span-text">{$t("window_title")}</span>
  </div>
  <div class="titlebar-center">
    <div class="container-input">
      <label class="search-label">
        <input
          type="text"
          name="text"
          class="input"
          required=""
          placeholder="Type here..."
        />
        <kbd class="slash-icon">Ctrl + F</kbd>
        <svg
          class="search-icon"
          xmlns="http://www.w3.org/2000/svg"
          version="1.1"
          xmlns:xlink="http://www.w3.org/1999/xlink"
          width="512"
          height="512"
          x="0"
          y="0"
          viewBox="0 0 56.966 56.966"
          style="enable-background:new 0 0 512 512"
          xml:space="preserve"
        >
          <g>
            <path
              d="M55.146 51.887 41.588 37.786A22.926 22.926 0 0 0 46.984 23c0-12.682-10.318-23-23-23s-23 10.318-23 23 10.318 23 23 23c4.761 0 9.298-1.436 13.177-4.162l13.661 14.208c.571.593 1.339.92 2.162.92.779 0 1.518-.297 2.079-.837a3.004 3.004 0 0 0 .083-4.242zM23.984 6c9.374 0 17 7.626 17 17s-7.626 17-17 17-17-7.626-17-17 7.626-17 17-17z"
              fill="currentColor"
              data-original="#000000"
              class=""
            ></path>
          </g>
        </svg>
      </label>
    </div>
  </div>
  <div class="buttons">
    <button
      type="button"
      class="titlebar-button"
      on:click={Open}
      use:tt={{ content: $t("settings_tooltip") }}
      aria-label="Settings"
    >
      <div class="button-icon">
        <SettingsIc />
      </div>
    </button>
    <button
      class="titlebar-button"
      id="titlebar-minimize"
      on:click={Window.Minimise}
      use:tt={{ content: $t("minimize_tooltip") }}
      aria-label="Minimize"
    >
      <div class="button-icon">
        <MinimizeIc />
      </div>
    </button>
    <button
      class="titlebar-button"
      id="titlebar-maximize"
      on:click={Window.ToggleMaximise}
      use:tt={{ content: $t("maximize_tooltip") }}
      aria-label="Maximize"
    >
      <div class="button-icon">
        <MaximizeIc />
      </div>
    </button>
    <button
      class="titlebar-button close"
      id="titlebar-close"
      on:click={Application.Quit}
      aria-label="Close"
      use:tt={{ content: $t("close_tooltip") }}
    >
      <div class="button-icon close-svg">
        <CloseIc />
      </div>
    </button>
  </div>
</div>

<style>
  .titlebar {
    --wails-draggable: drag;
    height: 30px;
    background: var(--titlebar-bg-color);
    user-select: none;
    display: flex;
    justify-content: space-between;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
  }

  .logo {
    display: flex;
    align-items: center;
    padding-left: 4px;
  }

  .icon {
    width: 26px;
    height: 26px;
  }

  .span-text {
    color: var(--titlebar-logo-text-color);
    margin-left: 10px;
    text-align: center;
    font-size: 16px;
    font-weight: 600 !important;
    white-space: nowrap;
  }

  .titlebar-center {
    position: absolute;
    left: 50%;
    top: 0;
    height: 100%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    pointer-events: none; /* чтобы не мешать drag-region */
    width: 300px; /* или auto, если нужно */
    justify-content: center;
  }

  .container-input {
    margin: 0 auto;
    position: relative;
    display: flex;
    align-items: center;
    pointer-events: auto; /* чтобы input был кликабелен */
    justify-content: center;
  }

  .search-label {
    --wails-draggable: no-drag;
    display: flex;
    height: 25px;
    align-items: center;
    box-sizing: border-box;
    position: relative;
    border: 1px solid var(--titlebar-search-border-color);
    border-radius: 12px;
    overflow: hidden;
    background: var(--titlebar-search-bg-color);
    padding: 7px;
    cursor: text;
  }

  .search-label:hover {
    border-color: var(--titlebar-search-border-color-hover);
  }

  .search-label:focus-within {
    background: var(--titlebar-search-bg-color-focus);
    border-color: var(--titlebar-search-border-color-focus);
  }

  .search-label input {
    outline: none;
    width: 100%;
    border: none;
    background: none;
    color: var(--titlebar-search-input-text-color);
  }

  .search-label input:focus + .slash-icon {
    display: none;
  }

  .search-label input:valid ~ .search-icon {
    display: block;
  }

  .search-label input:valid {
    width: calc(100% - 22px);
    transform: translateX(20px);
  }

  .search-label svg,
  .slash-icon {
    position: absolute;
    color: var(--titlebar-search-hotkey-text-color);
  }

  .search-icon {
    display: none;
    width: 12px;
    height: auto;
  }

  .slash-icon {
    right: 7px;
    border: 1px solid var(--titlebar-search-hotkey-border-color);
    background: var(--titlebar-search-hotkey-bg-color);
    display: flex; /* Добавляем flexbox */
    align-items: center; /* Выравнивание по вертикали по центру */
    justify-content: center; /* Выравнивание по горизонтали по центру */
    border-radius: 3px;
    box-shadow: var(--titlebar-search-hotkey-box-shadow);
    cursor: text;
    font-size: 12px;
    width: fit-content;
    height: 17px;
    padding: 0 3px;
  }

  .slash-icon:active {
    box-shadow: var(--titlebar-search-hotkey-box-shadow-active);
    text-shadow: var(--titlebar-search-hotkey-text-shadow-active);
    color: var(--titlebar-search-hotkey-text-color-active);
  }

  .titlebar-button {
    --wails-draggable: no-drag;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    width: 40px;
    height: 30px;
    background: var(--titlebar-button-bg-color);
    border: none;
    transition: 0.3s;
  }

  .button-icon {
    color: var(--titlebar-button-color);
    width: 24px;
    height: 24px;
  }

  .titlebar-button:hover {
    background: var(--titlebar-button-bg-color-hover);
  }

  .titlebar-button:hover .button-icon {
    color: var(--titlebar-button-color-hover);
  }

  .close:hover {
    background: var(--titlebar-close-button-bg-color-hover);
  }

  .close:hover .close-svg {
    color: var(--titlebar-close-button-color-hover);
  }
</style>
