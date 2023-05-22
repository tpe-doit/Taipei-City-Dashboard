<script setup>
import { onMounted, ref } from 'vue';
import DialogContainer from './DialogContainer.vue';
import CustomCheckBox from '../utilities/CustomCheckBox.vue'
import { useDialogStore } from '../../store/dialogStore';
import { useContentStore } from '../../store/contentStore';
import { validateStrInput } from '../../assets/utilityFunctions/validate'

const dialogStore = useDialogStore();
const contentStore = useContentStore();
const newName = ref('')
const error = ref(null)
const deleteConfirm = ref(false)

function handleSubmit() {
    if (validateStrInput(newName.value) !== true) {
        error.value = validateStrInput(newName.value)
        return
    }
    contentStore.changeCurrentDashboardName(newName.value);
    handleClose()
    dialogStore.hideAllDialogs();
}
function handleClose() {
    newName.value = ''
    error.value = null
    deleteConfirm.value = false
    dialogStore.hideAllDialogs();
}
function handleDelete() {
    contentStore.deleteCurrentDashboard();
    handleClose()
}
</script>

<template>
    <DialogContainer :dialog="`dashboardSettings`" @on-close="handleClose">
        <div class="dashboardsettings">
            <h2>儀表板設定</h2>
            <div class="dashboardsettings-input">
                <p v-if="error">{{ error }}</p>
                <label for="name">
                    更改名稱
                </label>
                <input name="name" v-model="newName" :placeholder="contentStore.currentDashboard.name" />
                <div :style="{ marginTop: '0.5rem' }">
                    <input type="checkbox" id="delete" :value="true" v-model="deleteConfirm" class="custom-check-input"
                        :style="{ display: 'none' }" />
                    <CustomCheckBox for="delete">啟動刪除儀表板功能</CustomCheckBox>
                </div>
            </div>

            <div class="dashboardsettings-button">
                <button class="dashboardsettings-button-simple" @click="handleClose">取消</button>
                <button class="dashboardsettings-button-fancy" @click="handleSubmit">確定更改</button>
                <button v-if="deleteConfirm" class="dashboardsettings-button-delete" @click="handleDelete">刪除儀表板</button>
            </div>
        </div>
    </DialogContainer>
</template>

<style scoped lang="scss">
.dashboardsettings {
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

        &-delete {
            background-color: rgb(192, 67, 67);
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