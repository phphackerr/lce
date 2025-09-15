<script>
  import { onMount } from "svelte";
  import { onConfigChanged } from "../../lib/watcher";

  let show = false;
  let changedFile = "";

  onMount(() => {
    onConfigChanged((filePath) => {
      changedFile = filePath;
      show = true;
    });
  });

  function reloadConfig() {
    console.log("🔄 Перезагрузка конфига:", changedFile);
    show = false;
    // Здесь можно дернуть loadSettings() или кастомный reloadConfig()
  }

  function ignoreChange() {
    show = false;
  }
</script>
