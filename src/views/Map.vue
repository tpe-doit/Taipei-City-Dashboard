<script setup>
import ComponentMapChart from '../components/components/ComponentMapChart.vue'
import MapContainer from '../components/map/MapContainer.vue'
import { useContentStore } from '../store/contentStore'

const contentStore = useContentStore()

</script>

<template>
    <div class="map">
        <div class="map-charts">
            <ComponentMapChart v-if="contentStore.currentDashboard.content"
                v-for="item in contentStore.currentDashboard.content" :content="item" :key="item.index" />
            <div v-else class="map-charts-nodashboard">
                <span>sentiment_very_dissatisfied</span>
                <h2>尚未加入組件</h2>
                <p>加入您的第一個組件</p>
            </div>
        </div>
        <MapContainer />
    </div>
</template>

<style scoped lang="scss">
.map {
    margin: var(--font-m) var(--font-m);
    display: flex;
    height: calc(100vh - 127px);

    &-charts {
        width: 300px;
        max-height: 100%;
        height: fit-content;
        display: grid;
        overflow-y: scroll;
        row-gap: var(--font-m);
        margin-right: var(--font-s);
        border-radius: 5px;

        @media (min-width: 1000px) {
            width: 350px;
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
                font-family: var(--font-icon);
                font-size: 2rem;
                margin-bottom: 1rem;
            }
        }
    }
}
</style>