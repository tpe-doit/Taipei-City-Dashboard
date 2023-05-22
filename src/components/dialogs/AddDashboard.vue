<script setup>
import { ref } from 'vue';
import DialogContainer from './DialogContainer.vue'
import { useDialogStore } from '../../store/dialogStore'
import { useContentStore } from '../../store/contentStore';
import { validateStrInput } from '../../assets/utilityFunctions/validate'

const dialogStore = useDialogStore()
const contentStore = useContentStore()
const name = ref('')
const error = ref(null)
function handleSubmit() {
    if (validateStrInput(name.value) !== true) {
        error.value = validateStrInput(name.value)
        return
    }
    contentStore.createNewDashboard(name.value);
    dialogStore.hideAllDialogs();
}
function handleClose() {
    name.value = '';
    error.value = null
    dialogStore.hideAllDialogs();
}
</script>

<template>
    <DialogContainer :dialog="`addDashboard`" @onClose="handleClose">
        <div class="adddashboard">
            <h2>新增自訂儀表板</h2>
            <div class="adddashboard-input">
                <p v-if="error">{{ error }}</p>
                <label for="name">
                    請輸入名稱
                </label>
                <input name="name" v-model="name" />
            </div>

            <div class="adddashboard-button">
                <button class="adddashboard-button-simple" @click="handleClose">取消</button>
                <button class="adddashboard-button-fancy" @click="handleSubmit">確定</button>
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
            font-size: var(--font-s);
            margin-bottom: 0.5rem
        }

        p {
            color: rgb(216, 52, 52);
        }

        input {
            background-color: transparent;
            border: solid 1px var(--color-border);
            border-radius: 5px;
            font-size: var(--font-m);
            padding: 4px 6px;


            &:focus {
                outline: none;
                border: solid 1px var(--color-highlight)
            }
        }
    }

    &-button {
        display: flex;
        justify-content: flex-end;

        &-simple {
            padding: 4px 6px;
            border-radius: 5px;
            margin: 0 2px;
            transition: color 0.2s;

            &:hover {
                color: var(--color-highlight)
            }
        }

        &-fancy {
            background-color: var(--color-highlight);
            padding: 4px 10px;
            border-radius: 5px;
            margin: 0 2px;
            transition: opacity 0.2s;

            &:hover {
                opacity: 0.8;
            }
        }
    }
}
</style>