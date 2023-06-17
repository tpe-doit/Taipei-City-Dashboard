<!-- Cleaned -->

<script setup>
import { ref } from 'vue';
import { useDialogStore } from '../../store/dialogStore'
import { useContentStore } from '../../store/contentStore'

import DialogContainer from './DialogContainer.vue'
import { validateStrInput } from '../../assets/utilityFunctions/validate'

const dialogStore = useDialogStore()
const contentStore = useContentStore()

// Stores the inputted dashboard name
const name = ref('')
// Stores the inputted icon
const icon = ref('')
const icons = ['shopping_cart', 'info', 'language', 'event', 'paid', 'account_balance', 'work', 'gavel', 'build_circle', 'circle_notifications', 'accessible', 'health_and_safety', 'science', 'coronavirus', 'luggage', 'flash_on', 'call', 'place', 'park', 'directions_car', 'lunch_dining', 'traffic', 'attractions', 'star', 'help', 'warning', 'lightbulb', 'notifications_active']
const errorMessage = ref(null)

function handleSubmit() {
    if (validateStrInput(name.value) !== true) {
        errorMessage.value = validateStrInput(name.value)
        return
    }
    // createNewDashboard is currently a dummy function to demonstrate what creating a new dashboard may look like
    // Connect a backend to actually implement the function or remove altogether
    contentStore.createNewDashboard(name.value, icon.value);
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
                <h3>
                    請輸入名稱
                </h3>
                <input type="text" v-model="name" />
            </div>
            <h3>請選擇圖示</h3>
            <div class="adddashboard-icon">
                <div v-for="item in icons">
                    <input type="radio" v-model="icon" :id="item" :value="item" />
                    <label :for="item">{{ item }}</label>
                </div>
            </div>
            <div class="adddashboard-control">
                <button class="adddashboard-control-cancel" @click="handleClose">取消</button>
                <button v-if="name && icon" class="adddashboard-control-confirm" @click="handleSubmit">確定</button>
            </div>
        </div>
    </DialogContainer>
</template>

<style scoped lang="scss">
.adddashboard {
    width: 300px;

    h3 {
        margin-bottom: 0.5rem;
        font-size: var(--font-s);
        font-weight: 400;
    }

    &-input {
        display: flex;
        flex-direction: column;
        margin: 1rem 0 0.5rem;

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

    &-icon {
        display: grid;
        grid-template-columns: 26px 26px 26px 26px 26px 26px 26px 26px 26px 26px;
        column-gap: 4px;
        row-gap: 4px;
        margin-bottom: 0.5rem;

        input {
            display: none;
            transition: border 0.2s;

            &:checked+label {
                border: solid 1px var(--color-highlight)
            }
        }

        label {
            width: 1.5rem;
            height: 1.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            border: solid 1px transparent;
            border-radius: 5px;
            font-size: 1.2rem;
            font-family: var(--font-icon);
            cursor: pointer;

            &:hover {
                border: solid 1px var(--color-border)
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