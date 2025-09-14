<script>
  //@ts-nocheck
  import { Close } from "../../../bindings/lce/backend/windows/settingswindow";
  import { CloseIc } from "../../lib/icons.js";
  import { tt } from "../../lib/tooltip.js";
  import { t } from "svelte-i18n";
</script>

<div class="titlebar">
  <div class="titlebar-center">
    <span class="span-text">{$t("settings")}</span>
  </div>
  <div class="buttons">
    <button
      class="titlebar-button close"
      id="titlebar-close"
      on:click={Close}
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
    background: var(--settings-titlebar-bg-color);
    user-select: none;
    display: flex;
    justify-content: space-between;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
  }

  .span-text {
    color: var(--settings-titlebar-text-color);
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

  .buttons {
    position: absolute;
    right: 0;
  }

  .titlebar-button {
    --wails-draggable: no-drag;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    width: 40px;
    height: 30px;
    background: var(--settings-titlebar-close-button-bg-color);
    border: none;
    transition: 0.3s;
  }

  .button-icon {
    color: var(--settings-titlebar-close-button-color);
  }

  .close:hover {
    background: var(--settings-titlebar-close-button-bg-color-hover);
  }

  .close:hover .close-svg {
    color: var(--settings-titlebar-close-button-color-hover);
  }
</style>
