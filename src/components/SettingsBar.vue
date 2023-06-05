<script setup>
import { ref, computed, onMounted } from 'vue';
import { useContentStore } from '../store/contentStore';
import { useDialogStore } from '../store/dialogStore';
import AddComponent from './dialogs/AddComponent.vue';
import DashboardSettings from './dialogs/DashboardSettings.vue';
import { useRoute, useRouter } from 'vue-router';
import MobileNavigation from './dialogs/MobileNavigation.vue';

const contentStore = useContentStore()
const dialogStore = useDialogStore()

const route = useRoute();
const router = useRouter();

const isDashboard = ref(false)
function handleToggle() {
    if (isDashboard.value) {
        router.replace({ name: 'mapview', query: { index: route.query.index } })
    } else {
        router.replace({ name: 'dashboard', query: { index: route.query.index } })
    }
}

onMounted(() => {
    setTimeout(() => {
        isDashboard.value = route.path === '/dashboard' ? false : true
    }, 100)
})

</script>

<template>
    <div class="settingsbar">
        <div class="settingsbar-title">
            <span>{{ contentStore.currentDashboard.icon }}</span>
            <h2>{{ contentStore.currentDashboard.name }}
            </h2>
            <button @click="dialogStore.showDialog('mobileNavigation')" class="show-if-mobile">
                <span class="settingsbar-title-navigation">arrow_drop_down_circle</span>
            </button>
            <MobileNavigation />
        </div>
        <div class="settingsbar-settings hide-if-mobile" v-if="contentStore.currentDashboard.index !== 'map-layers'">
            <button @click="dialogStore.showDialog('addComponent')"><span>add_chart</span>
                新增組件
            </button>
            <AddComponent />
            <button @click="dialogStore.showDialog('dashboardSettings')"><span>settings</span>
                設定
            </button>
            <DashboardSettings />
        </div>
        <div class="settingsbar-navigation show-if-mobile">
            <p>圖表</p>
            <label class="toggleswitch">
                <input type="checkbox" @change="handleToggle" v-model="isDashboard">
                <span class="toggleswitch-slider"></span>
            </label>
            <p>地圖</p>
        </div>
    </div>
</template>

<style scoped lang="scss">
.settingsbar {
    width: calc(100% - 2*var(--font-m));
    min-height: 1.6rem;
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

        &-navigation {
            color: var(--color-complement-text);
            margin-left: 4px;
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

    &-navigation {
        display: flex;
        align-items: center;

        p {
            color: var(--color-complement-text)
        }
    }

}

.toggleswitch {
    position: relative;
    display: inline-block;
    height: 1rem;
    width: 2rem;
    margin: 0 4px;

    input {
        opacity: 0;
        width: 0;
        height: 0;
    }

    &-slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: var(--color-complement-text);
        -webkit-transition: .4s;
        transition: .4s;
        border-radius: var(--font-m);
    }

    &-slider::before {
        position: absolute;
        content: "";
        height: 0.7rem;
        width: 0.7rem;
        left: 0.15rem;
        bottom: 0.15rem;
        background-color: var(--color-border);
        -webkit-transition: .4s;
        transition: .4s;
        border-radius: 50%;
    }

    // input:checked+&-slider {
    //     background-color: var(--color-highlight);
    // }

    // input:focus+&-slider {
    //     box-shadow: 0 0 1px var(--color-highlight);
    // }

    input:checked+&-slider:before {
        -webkit-transform: translateX(1rem);
        -ms-transform: translateX(1rem);
        transform: translateX(1rem);
        // background-color: white;
    }
}
</style>