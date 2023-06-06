<!-- Cleaned -->

<script setup>
import { ref } from 'vue';
import { useDialogStore } from '../../store/dialogStore'

import DialogContainer from './DialogContainer.vue'
import { validateStrInput } from '../../assets/utilityFunctions/validate'

const dialogStore = useDialogStore()

// Stores the inputted dashboard name
const name = ref('')
const errorMessage = ref(null)

function handleSubmit() {
    if (validateStrInput(name.value) !== true) {
        errorMessage.value = validateStrInput(name.value)
        return
    }
    // If a backend is connected, uncomment the following to implement the create dashboard function
    // contentStore.createNewDashboard(name.value);
    dialogStore.showNotification('fail', '尚未新增新增儀表板功能，無法新增儀表板')
    handleClose();
}
function handleClose() {
    name.value = '';
    errorMessage.value = null
    dialogStore.hideAllDialogs();
}
</script>

<template>
    <DialogContainer :dialog="`addDashboard`" @onClose="handleClose">
        <div class="adddashboard">
            <h2>新增自訂儀表板</h2>
            <div class="adddashboard-input">
                <p v-if="errorMessage">{{ errorMessage }}</p>
                <label for="name">
                    請輸入名稱
                </label>
                <input name="name" v-model="name" />
            </div>
            <div class="adddashboard-control">
                <button class="adddashboard-control-cancel" @click="handleClose">取消</button>
                <button class="adddashboard-control-confirm" @click="handleSubmit">確定</button>
            </div>
        </div>
    </DialogContainer>
</template>

<style scoped lang="scss">
.adddashboard {
    width: 300px;

    &-input {
        display: flex;
        flex-direction: column;
        margin: 1rem 0 1.5rem;

        label {
            margin-bottom: 0.5rem;
            font-size: var(--font-s);
        }

        p {
            color: rgb(216, 52, 52);
        }

        input {
            padding: 4px 6px;
            border: solid 1px var(--color-border);
            border-radius: 5px;
            background-color: transparent;
            font-size: var(--font-m);

            &:focus {
                outline: none;
                border: solid 1px var(--color-highlight)
            }
        }
    }

    &-control {
        display: flex;
        justify-content: flex-end;

        &-cancel {
            margin: 0 2px;
            padding: 4px 6px;
            border-radius: 5px;
            transition: color 0.2s;

            &:hover {
                color: var(--color-highlight)
            }
        }

        &-confirm {
            margin: 0 2px;
            padding: 4px 10px;
            border-radius: 5px;
            background-color: var(--color-highlight);
            transition: opacity 0.2s;

            &:hover {
                opacity: 0.8;
            }
        }
    }
}
</style>