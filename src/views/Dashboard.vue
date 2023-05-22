<script setup>

import { useContentStore } from '../store/contentStore'
import ComponentContainer from '../components/components/ComponentContainer.vue'
import MoreInfo from '../components/dialogs/MoreInfo.vue';
const contentStore = useContentStore()

</script>

<template>
    <div :class="{ dashboard: true, nodashboard: !contentStore.currentDashboard.content }">
        <ComponentContainer v-if="contentStore.currentDashboard.content"
            v-for="item in contentStore.currentDashboard.content" :content="item" />
        <div v-else class="dashboard-nodashboard">
            <span>sentiment_very_dissatisfied</span>
            <h2>尚未加入組件</h2>
            <p>加入您的第一個組件</p>
        </div>
        <MoreInfo />
    </div>
</template>

<style scoped lang="scss">
.dashboard {
    margin: var(--font-m) var(--font-m);
    display: grid;
    row-gap: var(--font-s);
    column-gap: var(--font-s);
    max-height: calc(100vh - 127px);
    overflow-y: scroll;

    @media (min-width: 820px) {
        grid-template-columns: 1fr 1fr;
    }

    @media (min-width: 1200px) {
        grid-template-columns: 1fr 1fr 1fr;
    }

    @media (min-width: 1800px) {
        grid-template-columns: 1fr 1fr 1fr 1fr;
    }

    @media (min-width: 2200px) {
        grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
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

.nodashboard {
    grid-template-columns: 1fr;
}
</style>