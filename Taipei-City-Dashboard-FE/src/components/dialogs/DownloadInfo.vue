<script setup>
import { useDialogStore } from "../../store/dialogStore";
import { useContentStore } from "../../store/contentStore";
import { useAuthStore } from "../../store/authStore";

import SideBarTab from "../utilities/miscellaneous/SideBarTab.vue";

const dialogStore = useDialogStore();
const contentStore = useContentStore();
const authStore = useAuthStore();
</script>

<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div
        v-if="dialogStore.dialogs.downloadInfo"
        class="dialogcontainer"
      >
        <div
          class="dialogcontainer-background"
          @click="dialogStore.hideAllDialogs"
        />
        <div class="dialogcontainer-dialog">
          <div class="downloadinfo">
            <div v-if="authStore.token">
				<h2>「臺北防災立即go」防災手冊</h2>
				<div>
					<SideBarTab
						:title="`地震篇`"
						:expanded="true"
						@click="handleClick(dialogStore, 'earthquake')"
					/>
					<SideBarTab
						:title="`火災篇`"
						:expanded="true"
						@click="handleClick(dialogStore, 'fire')"
					/>
					<SideBarTab
						:title="`颱風篇`"
						:expanded="true"
						@click="handleClick(dialogStore, 'typhoon')"
					/>
				</div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.dialogcontainer {
	width: 200vw;
	height: 100vh;
	height: calc(var(--vh) * 100);
	display: flex;
	align-items: center;
	justify-content: center;
	position: fixed;
	top: 0;
	right: 0;
	opacity: 1;
	z-index: 2;

	&-dialog {
		width: fit-content;
		height: fit-content;
		position: absolute;
		top: 110px;
		right: 45px;
		padding: var(--font-m);
		border: solid 1px var(--color-border);
		border-radius: 5px;
		background-color: rgb(30, 30, 30);
		transform: translateY(0);
		z-index: 5;
	}

	&-background {
		width: 100vw;
		height: 100vh;
		height: calc(var(--vh) * 100);
		position: absolute;
		top: 0;
		left: 0;
		background-color: rgba(0, 0, 0, 0.66);
		z-index: 4;
	}
}

.downloadinfo {
	width: 280px;
	max-height: 350px;
	overflow-y: scroll;
}

// Classes that are provided by vue transitions. Read the official docs for more instructions.
// https://vuejs.org/guide/built-ins/transition.html
.dialog-enter-from,
.dialog-leave-to {
	opacity: 0;

	.dialogcontainer-dialog {
		transform: translateY(-2.25rem);
	}
}

.dialog-enter-active,
.dialog-leave-active {
	transition: opacity 0.3s ease;

	.dialogcontainer-dialog {
		transition: transform 0.3s ease;
	}
}
</style>

<script>
export default {
  	methods: {
		handleClick(dialogStore, eventType) {
			let fileURL = '';
			let fileName = '';
			if (eventType === 'earthquake') {
				fileURL = 'https://www.eoc.gov.taipei/PropagandaManual/Earthquake/mobile/download/%E3%80%8C%E8%87%BA%E5%8C%97%E9%98%B2%E7%81%BD%E7%AB%8B%E5%8D%B3go%E3%80%8D%E9%98%B2%E7%81%BD%E6%89%8B%E5%86%8A_%E5%9C%B0%E9%9C%87.pdf';
				// fileURL = '../../assets/infoFiles/「臺北防災立即go」防災手冊_地震.pdf';
				fileName = '台北防災手冊_地震.pdf';
			}
			else if (eventType === 'fire') {
				fileURL = 'https://www.eoc.gov.taipei/PropagandaManual/Fire/mobile/download/%E3%80%8C%E8%87%BA%E5%8C%97%E9%98%B2%E7%81%BD%E7%AB%8B%E5%8D%B3go%E3%80%8D%E9%98%B2%E7%81%BD%E6%89%8B%E5%86%8A_%E7%81%AB%E7%81%BD.pdf';
				fileName = '台北防災手冊_火災.pdf';
			}
			else if (eventType === 'typhoon') {
				fileURL = 'https://www.eoc.gov.taipei/PropagandaManual/Typhoon/mobile/download/%E3%80%8C%E8%87%BA%E5%8C%97%E9%98%B2%E7%81%BD%E7%AB%8B%E5%8D%B3go%E3%80%8D%E9%98%B2%E7%81%BD%E6%89%8B%E5%86%8A_%E9%A2%B1%E9%A2%A8.pdf';
				fileName = '台北防災手冊_颱風.pdf';
			}

			const link = document.createElement('a');
			link.href = fileURL;
			link.download = fileName;
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
			dialogStore.hideAllDialogs();
		}
	}
}
</script>