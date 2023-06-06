<script setup>

import { useContentStore } from '../store/contentStore'
import { useDialogStore } from '../store/dialogStore';
import ComponentContainer from '../components/components/ComponentContainer.vue'
import MoreInfo from '../components/dialogs/MoreInfo.vue';

const contentStore = useContentStore()
const dialogStore = useDialogStore()

</script>

<template>
    <div v-if="contentStore.currentDashboard.index === 'map-layers'" class="dashboard">
        <ComponentContainer v-for="item in contentStore.currentDashboard.content" :content="item" :is-map-layer="true"
            :key="item.index" />
    </div>
    <div v-else-if="contentStore.currentDashboard.content.length !== 0" class="dashboard">
        <ComponentContainer v-for="item in contentStore.currentDashboard.content" :content="item" :key="item.index" />
        <MoreInfo />
    </div>
    <div v-else class="dashboard nodashboard">
        <div class="dashboard-nodashboard">
            <span>sentiment_very_dissatisfied</span>
            <h2>尚未加入組件</h2>
            <button @click="dialogStore.showDialog('addComponent')">加入您的第一個組件</button>
        </div>
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

    @media (min-width: 720px) {
        grid-template-columns: 1fr 1fr;
    }

    @media (min-width: 1150px) {
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

        button {
            color: var(--color-highlight)
        }
    }
}

.nodashboard {
    grid-template-columns: 1fr;
}
</style>