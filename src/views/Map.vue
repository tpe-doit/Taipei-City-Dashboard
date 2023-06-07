<!-- Cleaned -->

<!-- Map charts will be hidden in mobile mode and be replaced with the mobileLayers dialog -->

<script setup>
import { useContentStore } from '../store/contentStore'
import { useDialogStore } from '../store/dialogStore';

import ComponentMapChart from '../components/components/ComponentMapChart.vue'
import MapContainer from '../components/map/MapContainer.vue'

const contentStore = useContentStore()
const dialogStore = useDialogStore()
</script>

<template>
    <div class="map">
        <div class="map-charts hide-if-mobile">
            <!-- If the dashboard is map layers -->
            <ComponentMapChart v-if="contentStore.currentDashboard.index === 'map-layers'"
                v-for="item in contentStore.currentDashboard.content" :content="item" :key="`map-layer-${item.index}`"
                :is-map-layer="true" />
            <!-- other dashboards that have components -->
            <div v-else-if="contentStore.currentDashboard.content.length !== 0" class="map-charts">
                <ComponentMapChart v-for="item in contentStore.currentDashboard.content" :content="item"
                    :key="item.index" />
                <h2>基本圖層</h2>
                <ComponentMapChart v-for="item in contentStore.mapLayers" :content="item" :key="`map-layer-${item.index}`"
                    :is-map-layer="true" />
            </div>
            <!-- other dashboards that don't have components -->
            <div v-else class="map-charts-nodashboard">
                <span>sentiment_very_dissatisfied</span>
                <h2>尚未加入組件</h2>
                <button @click="dialogStore.showDialog('addComponent')">加入您的第一個組件</button>
            </div>
        </div>
        <MapContainer />
    </div>
</template>

<style scoped lang="scss">
.map {
    height: calc(100vh - 127px);
    display: flex;
    margin: var(--font-m) var(--font-m);

    &-charts {
        width: 360px;
        max-height: 100%;
        height: fit-content;
        display: grid;
        row-gap: var(--font-m);
        margin-right: var(--font-s);
        border-radius: 5px;
        overflow-y: scroll;

        @media (min-width: 1000px) {
            width: 370px;
        }

        @media (min-width: 2000px) {
            width: 400px;
        }

        &-nodashboard {
            width: 100%;
            height: calc(100vh - 127px);
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;

            span {
                margin-bottom: 1rem;
                font-family: var(--font-icon);
                font-size: 2rem;
            }

            button {
                color: var(--color-highlight);
            }
        }
    }
}
</style>