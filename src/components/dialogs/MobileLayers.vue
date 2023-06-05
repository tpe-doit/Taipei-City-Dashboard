<script setup>
import { computed, ref } from 'vue';
import { useDialogStore } from '../../store/dialogStore'
import { useContentStore } from '../../store/contentStore';
import MobileLayerTab from '../utilities/MobileLayerTab.vue';

const dialogStore = useDialogStore()
const contentStore = useContentStore()

const filteredMapLayers = computed(() => {
    if (!contentStore.currentDashboard.content) {
        return []
    }
    return contentStore.currentDashboard.content.filter((element) => element.map_config)
})
</script>

<template>
    <Teleport to="body">
        <div :class="{ dialogcontainer: true, 'show-dialog': dialogStore.dialogs.mobileLayers === true }">
            <div class="dialogcontainer-background" @click="dialogStore.hideAllDialogs"></div>
            <div class="dialogcontainer-dialog">
                <div class="mobilelayers">
                    <MobileLayerTab v-if="contentStore.currentDashboard.index === 'map-layers'"
                        v-for="item in contentStore.currentDashboard.content" :content="item"
                        :key="`map-layer-${item.index}`" />
                    <div v-else-if="filteredMapLayers.length !== 0">
                        <MobileLayerTab v-for="item in filteredMapLayers" :content="item" :key="item.index" />
                        <h2>基本圖層</h2>
                        <MobileLayerTab v-for="item in contentStore.mapLayers" :content="item"
                            :key="`map-layer-${item.index}`" />
                    </div>
                    <div v-else>
                        <h2>基本圖層</h2>
                        <MobileLayerTab v-for="item in contentStore.mapLayers" :content="item"
                            :key="`map-layer-${item.index}`" />
                    </div>
                </div>
            </div>
        </div>
    </Teleport>
</template>

<style scoped lang="scss">
.dialogcontainer {
    width: 100vw;
    height: 100vh;
    top: 0;
    left: 0;
    position: fixed;
    opacity: 0;
    z-index: -1;

    &-dialog {
        position: absolute;
        background-color: rgb(30, 30, 30);
        border: solid 1px var(--color-border);
        border-radius: 5px;
        padding: var(--font-m) 11px var(--font-m) var(--font-m);
        width: fit-content;
        height: fit-content;
        transform: translateY(0);
        left: 16px;
        top: 135px
    }

    &-background {
        position: absolute;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        background-color: rgba(0, 0, 0, 0.5);
    }

}

@keyframes opacity-transition {
    0% {
        opacity: 0
    }

    100% {
        opacity: 1
    }
}

@keyframes transform-transition {
    0% {
        transform: translateY(-2.25rem);
    }

    100% {
        transform: translateY(0);
    }
}

.mobilelayers {
    width: 80px;
    max-height: 350px;
    overflow-y: scroll;
    padding: 0;

    h2 {
        margin-bottom: 8px;
    }
}

.show-dialog {
    z-index: 10;
    opacity: 1;
    animation: opacity-transition 0.3s ease;

    .dialogcontainer-dialog {
        animation: transform-transition 0.3s ease;
    }
}
</style>