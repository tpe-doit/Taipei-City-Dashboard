<script setup>
const VITE_APP_TITLE = import.meta.env.VITE_APP_TITLE
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../store/authStore';
import { useFullscreen } from '@vueuse/core'
import { useDialogStore } from '../store/dialogStore';
import UserSettings from './dialogs/UserSettings.vue';

const route = useRoute();
const authStore = useAuthStore()
const dialogStore = useDialogStore()
const { isFullscreen, toggle } = useFullscreen()

const linkQuery = computed(() => {
    const query = route.query
    return `?id=${query.id}&type=${query.type}`
})

</script>

<template>
    <div class="navbar">
        <div class="navbar-logo">
            <img src="../assets/images/TUIC.svg" alt="tuic logo" />
            <div>
                <h1>{{ VITE_APP_TITLE }}</h1>
                <h2>Taipei City Dashboard</h2>
            </div>
        </div>
        <div class="navbar-tabs" v-if="authStore.auth > 0">
            <router-link :to="`/Dashboard${linkQuery}`">儀表板總覽</router-link>
            <router-link :to="`/Mapview${linkQuery}`">地圖交叉比對</router-link>
        </div>
        <div class="navbar-user" v-if="authStore.auth > 0">
            <button><span>school</span></button>
            <button @click="toggle"><span>{{ isFullscreen ? 'fullscreen_exit' : 'fullscreen' }}</span></button>
            <div class="navbar-user-user">
                <button>{{ authStore.user.name }}</button>
                <ul>
                    <li><button @click="dialogStore.showDialog('userSettings')">用戶設定</button></li>
                    <li><button @click="authStore.handleLogout">登出</button></li>
                </ul>
                <teleport to="body">
                    <user-settings />
                </teleport>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.navbar {
    height: 60px;
    width: 100vw;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--color-border);
    background-color: var(--color-component-background);

    &-logo {
        display: flex;

        h1 {
            font-weight: 500;
        }

        h2 {
            font-size: var(--font-s);
            font-weight: 400;
        }

        img {
            filter: invert(1);
            height: 45px;
            margin: 0 var(--font-m)
        }
    }

    &-tabs {
        display: flex;

        a {
            display: flex;
            align-items: center;
            height: 59px;
            margin-left: var(--font-s)
        }

        .router-link-active {
            color: var(--color-highlight);
            border-bottom: solid 3px var(--color-highlight);
        }
    }

    &-user {
        display: flex;
        align-items: center;

        button {
            margin-right: var(--font-m);
            font-size: var(--font-m);
            display: flex;
            align-items: center;
            padding: 2px 4px;
            border-radius: 4px;
            transition: background-color 0.25s;
        }

        button:hover {
            background-color: var(--color-complement-text);
        }

        span {
            font-family: var(--font-icon);
            font-size: calc(var(--font-l) * var(--font-to-icon));
        }

        &-user:hover ul {
            display: block;
            opacity: 1;
        }

        &-user {
            height: 60px;
            display: flex;
            align-items: center;

            ul {
                position: absolute;
                z-index: 10;
                background-color: rgb(85, 85, 85);
                padding: 8px;
                border-radius: 5px;
                min-width: 100px;
                display: none;
                opacity: 0;
                transition: opacity 0.25s;
                right: 20px;
                top: 55px;

                li {
                    cursor: pointer;
                    padding: 8px 4px;
                    transition: background-color 0.25s;
                    border-radius: 5px;
                }

                li:hover {
                    background-color: var(--color-complement-text)
                }
            }
        }
    }


}
</style>