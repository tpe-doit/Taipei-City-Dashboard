<script setup>
import { onMounted } from 'vue';
import { useMapStore } from '../../store/mapStore';
import { useDialogStore } from '../../store/dialogStore';
import { useContentStore } from '../../store/contentStore';
import MobileLayers from '../dialogs/MobileLayers.vue';

const mapStore = useMapStore()
const dialogStore = useDialogStore()
const contentStore = useContentStore()

onMounted(() => {
    mapStore.initializeMapBox()
})
</script>

<template>
    <div class="mapcontainer">
        <div id="mapboxBox">
            <button class="mapcontainer-layers show-if-mobile"
                @click="dialogStore.showDialog('mobileLayers')"><span>layers</span></button>
            <KeepAlive>
                <MobileLayers :key="contentStore.currentDashboard.index" />
            </KeepAlive>
        </div>
        <div class="mapcontainer-controls hide-if-mobile">
            <button @click="mapStore.easeToLocation([121.536609, 25.044808], 12.5, 4000, 0, 0)">返回預設</button>
            <button
                @click="mapStore.easeToLocation([121.56719891052944, 25.036536501522235], 15.5, 4000, 67.5, 45)">台北市政府</button>
            <button
                @click="mapStore.easeToLocation([121.55676686461811, 25.038745980403505], 14.5, 4000, 67, 35)">東區</button>
            <button @click="mapStore.easeToLocation([121.52231934552663, 25.09131257257215], 15, 4000, 50, 30)">士林</button>
            <button
                @click="mapStore.easeToLocation([121.5228860863223, 25.068515771454074], 15.75, 4000, 60, 130)">花博公園</button>
            <button
                @click="mapStore.easeToLocation([121.51013649269608, 25.055775323336594], 16, 4000, 60, 50)">大稻埕</button>
            <button
                @click="mapStore.easeToLocation([121.55180953513364, 25.12898719167495], 13.25, 4000, 60, 0)">陽明山</button>
        </div>
    </div>
</template>

<style scoped lang="scss">
.mapcontainer {
    width: 100%;
    height: calc(100%);
    flex: 1;

    &-controls {
        margin-top: 8px;
        display: flex;

        button {
            margin-right: 6px;
            color: var(--color-complement-text);
            height: 1.5rem;
            width: fit-content;
            padding: 4px;
            border-radius: 5px;
            background-color: var(--color-component-background);

            &:focus {
                animation-name: colorfade;
                animation-duration: 4s;
            }
        }
    }

    &-layers {
        background-color: white;
        width: 1.75rem;
        height: 1.75rem;
        position: absolute;
        z-index: 1;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        right: 10px;
        top: 108px;

        span {
            font-family: var(--font-icon);
            color: var(--color-component-background);
            font-size: 1.2rem;
        }
    }
}

#mapboxBox {
    width: 100%;
    height: calc(100% - 32px);
    border-radius: 5px;

    @media (max-width: 1000px) {
        height: 100%;
    }
}

@keyframes colorfade {
    0% {
        color: var(--color-highlight)
    }

    75% {
        color: var(--color-highlight)
    }

    100% {
        color: var(--color-complement-text)
    }
}
</style>