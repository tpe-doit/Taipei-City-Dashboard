<!-- Cleaned -->

<script setup>
import { onBeforeMount, onMounted } from 'vue';
import { useAuthStore } from './store/authStore';
import { useDialogStore } from './store/dialogStore';

import NavBar from './components/NavBar.vue';
import SideBar from './components/SideBar.vue';
import SettingsBar from './components/SettingsBar.vue';
import NotificationBar from './components/dialogs/NotificationBar.vue';
import InitialWarning from './components/dialogs/InitialWarning.vue';

const authStore = useAuthStore();
const dialogStore = useDialogStore();

onBeforeMount(() => {
	authStore.setUser();
	authStore.checkIfMobile();
});
onMounted(() => {
	const showInitialWarning = localStorage.getItem('initialWarning');
	if (!showInitialWarning) {
		dialogStore.showDialog('initialWarning');
	}
});
</script>

<template>
	<div class="app-container">
		<NotificationBar />
		<NavBar />
		<div :class="{ 'app-content': true, 'app-mobile': authStore.isMobileDevice }">
			<SideBar />
			<div class="app-content-main">
				<SettingsBar />
				<RouterView></RouterView>
			</div>
		</div>
		<InitialWarning />
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
		height: calc(100vh - 60px);
		display: flex;

		&-main {
			width: 100%;
			display: flex;
			flex-direction: column;
		}
	}

	&-mobile {
		height: calc(100vh - 120px);
	}
}
</style>