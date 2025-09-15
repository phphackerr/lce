<script>
  export let show = false;
  export let text = "Загрузка..."; // Текст по умолчанию, можно переопределить
</script>

{#if show}
  <div class="loading-overlay">
    <div class="loading-content">
      <div class="spinner"></div>
      <span class="loading-text">{text}</span>
    </div>
  </div>
{/if}

<style>
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
</style>
