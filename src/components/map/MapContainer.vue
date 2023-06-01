<script setup>

import { onMounted } from 'vue';
import { useMapStore } from '../../store/mapStore';

const mapStore = useMapStore()

onMounted(() => {
    mapStore.initializeMapBox()
})
</script>

<template>
    <div class="mapcontainer">
        <div id="mapboxBox" />
        <div class="mapcontainer-controls">
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
    height: calc(100% - var(--font-m));
    flex: 1;

    &-controls {
        margin-top: 8px;
        display: flex;

        @media (max-width: 1000px) {
            display: none;
        }

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