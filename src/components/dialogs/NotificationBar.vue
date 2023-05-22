<script setup>
import { useDialogStore } from '../../store/dialogStore';

const dialogStore = useDialogStore()
</script>

<template>
    <Teleport to="body">
        <Transition name="notification">
            <div class="notificationcontainer" v-if="dialogStore.dialogs.notificationBar">
                <div class="notificationcontainer-notification">
                    <span
                        :class="{ success: dialogStore.notification.status === 'success', fail: dialogStore.notification.status !== 'success' }">{{
                            dialogStore.notification.status === 'success' ? 'check_circle' : 'error' }}</span>
                    <h5
                        :class="{ success: dialogStore.notification.status === 'success', fail: dialogStore.notification.status !== 'success' }">
                        {{ dialogStore.notification.message }}</h5>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<style scoped lang="scss">
.notificationcontainer {
    position: fixed;
    width: 100vw;
    top: 20px;
    display: flex;
    justify-content: center;
    padding-bottom: 1rem;

    &-notification {
        height: 3rem;
        width: fit-content;
        background-color: rgb(63, 63, 63);
        border: solid 1px var(--color-border);
        border-radius: 5px;
        display: flex;
        align-items: center;
        padding: 0 1rem;
        box-shadow: 0px 5px 10px black;

        span {
            font-family: var(--font-icon);
            font-size: var(--font-l);
            margin-right: 10px;
        }

        h5 {
            font-weight: 400;
        }
    }
}

.success {
    color: greenyellow
}

.fail {
    color: rgb(237, 90, 90)
}

.notification-enter-from,
.notification-leave-to {
    top: -60px;
}

.notification-enter-active,
.notification-leave-active {
    transition: top 0.3s ease;
}
</style>