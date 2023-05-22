<script setup>
import { useDialogStore } from '../../store/dialogStore'

const dialogStore = useDialogStore()
const props = defineProps({ dialog: String })
const emits = defineEmits(['onClose'])
</script>

<template>
    <Teleport to="body">
        <Transition name="dialog">
            <div class="dialogcontainer" v-if="dialogStore.dialogs[dialog]">
                <div class="dialogcontainer-background" @click="$emit('onClose')"></div>
                <div class="dialogcontainer-dialog">
                    <slot></slot>
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

    &-dialog {
        position: relative;
        z-index: 2;
        background-color: var(--color-component-background);
        border: solid 1px var(--color-border);
        border-radius: 5px;
        padding: var(--font-m);
        width: fit-content;
        height: fit-content;
        transform: translateY(0);
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
</style>