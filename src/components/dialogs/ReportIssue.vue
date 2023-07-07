<!-- Cleaned -->

<script setup>
import { ref } from 'vue';
import { useDialogStore } from '../../store/dialogStore';
import DialogContainer from './DialogContainer.vue';

const dialogStore = useDialogStore();

const allInputs = ref({
    type: "組件基本資訊有誤",
    description: "",
    name: "",
    contact: "",
})
const issueTypes = ["組件基本資訊有誤", "組件資料有誤或未更新", "系統問題", "其他建議"]

function handleSubmit() {
    let openedWindow
    openedWindow = window.open(`https://docs.google.com/forms/d/e/1FAIpQLScgY1xXYuGTb1daf3rFg0yyH-x1A6CSCuwF2P1GBENzx7XPug/formResponse?usp=pp_url&entry.2090216374=${dialogStore.issue.id}&entry.1790274053=${dialogStore.issue.name}&entry.1025623950=${allInputs.value.type}&entry.11266555=${allInputs.value.description}&entry.1057102353=${allInputs.value.name}&entry.502116166=${allInputs.value.contact}&submit=Submit`)
    setTimeout(() => {
        openedWindow.close()
        dialogStore.showNotification('success', '感謝您的回報，我們將儘速處理')
        handleClose();
    }, 200)

}
function handleClose() {
    allInputs.value = {
        type: "組件基本資訊有誤",
        description: "",
        name: "",
        contact: "",
    }
    dialogStore.dialogs.reportIssue = false
}
</script>

<template>
    <DialogContainer dialog="reportIssue" @on-close="handleClose">
        <div class="reportissue">
            <h2>回報問題</h2>
            <h3>問題種類 (必填)</h3>
            <div v-for="item in issueTypes">
                <input class="reportissue-radio" type="radio" v-model="allInputs.type" :value="item" :id="item" />
                <label :for="item">
                    <div></div>
                    {{ item }}
                </label>
            </div>
            <h3>問題描述 (必填)</h3>
            <textarea v-model="allInputs.description"></textarea>
            <h3>姓名 (非必填)</h3>
            <input class="reportissue-input" type="text" v-model="allInputs.name" />
            <h3>任何聯絡方式 (非必填)</h3>
            <input class="reportissue-input" type="text" v-model="allInputs.contact" />
            <div class="reportissue-control">
                <button class="reportissue-control-cancel" @click="handleClose">取消</button>
                <button v-if="allInputs.description" class="reportissue-control-confirm" @click="handleSubmit">送出</button>
            </div>
        </div>
    </DialogContainer>
</template>

<style scoped lang="scss">
.reportissue {
    width: 300px;
    display: flex;
    flex-direction: column;

    h3 {
        margin: 0.5rem 0;
        font-size: var(--font-s);
        font-weight: 400;
    }

    &-radio {
        display: none;

        &:checked+label {
            color: white;

            div {
                background-color: var(--color-highlight);
            }
        }
    }

    label {
        position: relative;
        display: flex;
        align-items: center;
        font-size: var(--font-s);
        color: var(--color-complement-text);
        transition: color 0.2s;
        cursor: pointer;

        div {
            width: calc(var(--font-s) / 2);
            height: calc(var(--font-s) / 2);
            margin-right: 4px;
            padding: calc(var(--font-s) / 4);
            border-radius: 50%;
            border: 1px solid var(--color-border);
            transition: background-color 0.2s;
        }
    }

    textarea {
        height: 100px;
        padding: 4px 6px;
        border: solid 1px var(--color-border);
        border-radius: 5px;
        background-color: transparent;
        font-size: var(--font-m);
        resize: none;

        &:focus {
            outline: none;
            border: solid 1px var(--color-highlight);
        }
    }

    &-input {
        padding: 4px 6px;
        border: solid 1px var(--color-border);
        border-radius: 5px;
        background-color: transparent;
        font-size: var(--font-m);

        &:focus {
            outline: none;
            border: solid 1px var(--color-highlight);
        }
    }

    &-control {
        display: flex;
        justify-content: flex-end;
        margin-top: 1rem;

        &-cancel {
            margin: 0 2px;
            padding: 4px 6px;
            border-radius: 5px;
            transition: color 0.2s;

            &:hover {
                color: var(--color-highlight);
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