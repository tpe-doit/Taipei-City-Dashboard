<script setup>
import { useDialogStore } from '../../store/dialogStore'
import { useContentStore } from '../../store/contentStore';
import SideBarTab from '../utilities/SideBarTab.vue';

const dialogStore = useDialogStore()
const contentStore = useContentStore()
</script>

<template>
    <Teleport to="body">
        <Transition name="dialog">
            <div class="dialogcontainer" v-if="dialogStore.dialogs.mobileNavigation">
                <div class="dialogcontainer-background" @click="dialogStore.hideAllDialogs"></div>
                <div class="dialogcontainer-dialog">
                    <div class="mobilenavigation">
                        <h2>儀表板列表</h2>
                        <SideBarTab v-for="item in contentStore.dashboards.filter((item) => item.index !== 'map-layers')"
                            :icon="item.icon" :title="item.name" :index="item.index" :key="item.index" :expanded="true"
                            @click="dialogStore.hideAllDialogs" />
                        <h2>基本地圖圖層</h2>
                        <SideBarTab :icon="`public`" :title="`圖資資訊`" index="map-layers" :expanded="true"
                            @click="dialogStore.hideAllDialogs" />
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<style scoped lang="scss">
.dialogcontainer {
    width: 100vw;
    height: 100vh;
    top: 0;
    left: 0;
    position: fixed;
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 1;
    z-index: 10;

    &-dialog {
        position: absolute;
        z-index: 2;
        background-color: rgb(30, 30, 30);
        border: solid 1px var(--color-border);
        border-radius: 5px;
        padding: var(--font-m);
        width: fit-content;
        height: fit-content;
        transform: translateY(0);
        top: 110px;
        left: 45px
    }

    &-background {
        position: absolute;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        background-color: rgba(0, 0, 0, 0.66);
    }


}

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

.mobilenavigation {
    width: 170px;
    max-height: 350px;
    overflow-y: scroll;
}
</style>