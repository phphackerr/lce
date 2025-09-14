<script>
  import { updateGamePath } from "../../../../store/appSettings";
  import { DeleteIc } from "../../../../lib/icons";
  import { tt } from "../../../../lib/tooltip";

  export let options = [];
  export let selectedValue;
  export let name = "custom-radio";
  export let onDelete = (value) => {};

  async function handleSelectionChange(newValue) {
    selectedValue = newValue;
    await updateGamePath(newValue);
  }
</script>

<div class="custom-radio-group">
  {#each options as option (option.value)}
    <label
      class="custom-radio-container"
      use:tt={{ content: option.label, placement: "top" }}
    >
      <input
        type="radio"
        {name}
        value={option.value}
        checked={selectedValue === option.value}
        on:change={() => handleSelectionChange(option.value)}
      />
      <span class="custom-radio-checkmark"></span>
      <span class="radio-label-text">{option.label}</span>
      <button
        class="delete-button"
        on:click|stopPropagation={() => onDelete(option.value)}
      >
        <DeleteIc />
      </button>
    </label>
  {/each}
</div>

<style>
  .custom-radio-group {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: calc(100% - 30px);
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.2);
    padding: 16px;
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.3);
  }
  .custom-radio-container {
    position: relative;
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 12px 20px;
    border-radius: 8px;
    background-color: rgba(255, 255, 255, 0.2);
    transition:
      background-color 0.3s ease,
      transform 0.3s ease,
      box-shadow 0.3s ease;
    font-size: 16px;
    color: #333333;
    user-select: none;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  .custom-radio-container:hover {
    background-color: rgba(255, 255, 255, 0.3);
    transform: scale(1.03);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.4);
  }
  .custom-radio-container input[type="radio"] {
    opacity: 0;
    position: absolute;
  }
  .custom-radio-checkmark {
    position: relative;
    height: 24px;
    width: 24px;
    border: 2px solid #ffffff;
    border-radius: 50%;
    background-color: rgba(0, 0, 0, 0.3);
    transition:
      background-color 0.4s ease,
      transform 0.4s ease;
    margin-right: 12px;
    display: inline-block;
    vertical-align: middle;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
  }
  .custom-radio-container
    input[type="radio"]:checked
    + .custom-radio-checkmark {
    background-color: #ffffff;
    border-color: #007bff;
    box-shadow: 0 0 0 8px rgba(0, 123, 255, 0.2);
    transform: scale(1.2);
    animation: pulse 0.6s forwards;
  }
  .custom-radio-checkmark::after {
    content: "";
    position: absolute;
    display: none;
  }
  .custom-radio-container
    input[type="radio"]:checked
    + .custom-radio-checkmark::after {
    display: block;
    left: 50%;
    top: 50%;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: #007bff;
    transform: translate(-50%, -50%);
  }

  /* NEW: Стили для текста метки и кнопки удаления */
  .radio-label-text {
    flex-grow: 1; /* Позволяет тексту занимать все доступное пространство */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis; /* Добавляет многоточие, если текст слишком длинный */
  }

  .delete-button {
    width: 35px;
    height: 35px;
    background: none;
    border: none;
    color: #ff4d4d; /* Красный цвет для кнопки удаления */
    font-size: 1.2em;
    cursor: pointer;
    margin-left: 10px; /* Отступ от текста */
    padding: 0;
    line-height: 1; /* Убираем лишний отступ */
    transition: color 0.2s ease;
  }

  .delete-button:hover {
    color: #cc0000; /* Темно-красный при наведении */
    transform: scale(1.1);
  }

  @keyframes pulse {
    0% {
      transform: scale(1.2);
    }
    50% {
      transform: scale(1.4);
    }
    100% {
      transform: scale(1.2);
    }
  }
</style>
