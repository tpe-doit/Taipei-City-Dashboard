<script setup>
import SideBarTab from './utilities/SideBarTab.vue';
import { ref } from 'vue';
import { useContentStore } from '../store/contentStore';
import { useDialogStore } from '../store/dialogStore';
import { useMapStore } from '../store/mapStore';
import AddDashboard from './dialogs/AddDashboard.vue';

const contentStore = useContentStore()
const dialogStore = useDialogStore()
const mapStore = useMapStore()

const isExpanded = ref(true);
function toggleExpand() {
    isExpanded.value = isExpanded.value ? false : true
    if (!isExpanded.value) {
        mapStore.resizeMap();
    }
}
</script>

<template>
    <div :class="{ sidebar: true, collapse: !isExpanded }">
        <div>
            <div class="sidebar-sub">
                <div class="sidebar-sub-add">
                    <h2>{{ isExpanded ? `儀表板列表` : `列表` }}</h2>
                    <button v-if="isExpanded" @click="dialogStore.showDialog('addDashboard')">新增</button>
                    <AddDashboard />
                </div>
                <SideBarTab v-for="item in contentStore.dashboards.filter((item) => item.index !== 'map-layers')"
                    :icon="item.icon" :title="item.name" :index="item.index" :key="item.index" :expanded="isExpanded" />
            </div>
            <div class="sidebar-sub">
                <h2>{{ isExpanded ? `基本地圖圖層` : `圖層` }}</h2>
                <SideBarTab :icon="`public`" :title="`圖資資訊`" :expanded="isExpanded" index="map-layers" />
            </div>
        </div>
        <button class="sidebar-collapse" @click="toggleExpand"><span>{{ isExpanded ? "keyboard_double_arrow_left" :
            "keyboard_double_arrow_right"
        }}</span></button>
    </div>
</template>

<style scoped lang="scss">
.sidebar {
    width: 190px;
    border-right: 1px solid var(--color-border);
    padding: 0 var(--font-m);
    margin-top: 20px;
    max-height: calc(100vh - 80px);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    transition: width 0.2s ease-out;
    overflow-x: hidden;

    h2 {
        color: var(--color-complement-text);
        font-weight: 400;
    }

    &-sub {
        margin-bottom: var(--font-s);
        min-width: 190px;

        &-add {
            display: flex;

            button {
                color: var(--color-highlight);
                margin-left: 3rem;
            }
        }
    }

    &-collapse {
        padding: 5px;
        border-radius: 5px;
        transition: background-color 0.2s;
        height: fit-content;
        align-self: flex-end;
        margin-bottom: 10px;

        &:hover {
            background-color: var(--color-component-background);
        }

        span {
            font-family: var(--font-icon);
            font-size: var(--font-l);
        }
    }
}

.collapse {
    width: 45px;

    h2 {
        margin-left: 5px;
    }
}
</style>