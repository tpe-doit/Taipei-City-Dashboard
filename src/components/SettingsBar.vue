<script setup>
import { useContentStore } from '../store/contentStore';
import { useDialogStore } from '../store/dialogStore';
import AddComponent from './dialogs/AddComponent.vue';
import DashboardSettings from './dialogs/DashboardSettings.vue';
const contentStore = useContentStore()
const dialogStore = useDialogStore()

</script>

<template>
    <div class="settingsbar">
        <div class="settingsbar-title">
            <span>{{ contentStore.currentDashboard.icon }}</span>
            <h2>{{ contentStore.currentDashboard.name === 'defaultFav' ? '我的最愛' : contentStore.currentDashboard.name }}</h2>
            <p>|</p>
            <h3>{{ contentStore.currentDashboard.type === 'fixed' ? '官方儀表板' : '自訂儀表板' }}</h3>
        </div>
        <div class="settingsbar-settings">
            <button @click="dialogStore.showDialog('addComponent')"><span>add_chart</span>
                新增組件
            </button>
            <AddComponent />
            <button @click="dialogStore.showDialog('dashboardSettings')"><span>settings</span>
                設定
            </button>
            <DashboardSettings />
        </div>
    </div>
</template>

<style scoped lang="scss">
.settingsbar {
    width: calc(100% - 2*var(--font-m));
    border-bottom: solid 1px var(--color-border);
    display: flex;
    justify-content: space-between;
    padding-bottom: 0.5rem;
    margin: 20px var(--font-m) 0;

    &-title {
        display: flex;
        align-items: center;

        span {
            font-family: var(--font-icon);
            font-size: calc(var(--font-m) * var(--font-to-icon));
        }

        p {
            font-size: var(--font-l);
            margin-left: var(--font-s);
            font-weight: 200;
            line-height: var(--font-l);
        }

        h2,
        h3 {
            font-weight: 400;
            font-size: var(--font-m);
            margin-left: var(--font-s);
        }

        h3,
        p {
            color: var(--color-complement-text)
        }
    }

    &-settings {
        display: flex;
        align-items: center;

        span {
            font-family: var(--font-icon);
            font-size: calc(var(--font-m) * var(--font-to-icon));
            margin-right: 4px
        }

        button {
            display: flex;
            align-items: center;
            margin-left: 8px;
            font-size: var(--font-m);
            border-radius: 5px;

            &:first-child {
                background-color: var(--color-highlight);
                padding: 2px 4px;
            }

            &:last-child {
                color: var(--color-highlight);
                border: solid 1px var(--color-highlight);
                padding: 1px 3px;

                span {
                    color: var(--color-highlight)
                }
            }
        }
    }
}
</style>