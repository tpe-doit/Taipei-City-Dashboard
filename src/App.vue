<script setup>
import NavBar from './components/NavBar.vue';
import SideBar from './components/SideBar.vue';
import SettingsBar from './components/SettingsBar.vue'
import { useAuthStore } from './store/authStore';
import { onBeforeMount } from 'vue';
import NotificationBar from './components/dialogs/NotificationBar.vue';

const authStore = useAuthStore()

onBeforeMount(() => {
  authStore.setUser();
})
</script>

<template>
  <div class="app-container">
    <NotificationBar />
    <NavBar />
    <div class="app-content">
      <SideBar />
      <div class="app-content-main">
        <SettingsBar />
        <RouterView></RouterView>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.app {
  &-container {
    max-width: 100vw;
    max-height: 100vh;
  }

  &-content {
    width: 100vw;
    max-width: 100vw;
    display: flex;
    height: calc(100vh - 60px);

    &-main {
      display: flex;
      flex-direction: column;
      width: 100%;
    }
  }
}

@media (max-width: 760px) {
  .app-content {
    height: calc(100vh - 120px);
  }
}
</style>