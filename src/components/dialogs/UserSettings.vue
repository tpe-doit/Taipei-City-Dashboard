<script setup>
import DialogContainer from './DialogContainer.vue';
import { useDialogStore } from '../../store/dialogStore';
import { useAuthStore } from '../../store/authStore';

const dialogStore = useDialogStore()
const authStore = useAuthStore()

function handleClose() {
    dialogStore.hideAllDialogs()
}
const accountType = ['Email用戶', '台北通', '台北on']

</script>

<template>
    <DialogContainer @onClose="handleClose" :dialog="`userSettings`">
        <h2>用戶設定</h2>
        <table class="usersettings">
            <tr>
                <th>名稱</th>
                <td colspan="3">{{ authStore.user.name }}</td>
            </tr>
            <tr>
                <th>類型</th>
                <td>{{ accountType[authStore.user.type] }}</td>
                <th>用戶代碼</th>
                <td>{{ authStore.user.id }}</td>
            </tr>
            <tr>
                <th>權限</th>
                <td>{{ authStore.user.status === 1 ? '管理員' : '一般用戶' }}</td>
                <th>帳戶狀態</th>
                <td>{{ authStore.user.status > 0 ? '啟用' : '停用' }}</td>
            </tr>
        </table>
    </DialogContainer>
</template>

<style scoped lang="scss">
.usersettings {
    margin: 1rem 0 0.5rem;
    width: 500px;
    border-spacing: 0;
    border: solid 1px var(--color-border);
    border-radius: 5px;

    td,
    th {
        border: solid 1px var(--color-border);
        height: 2rem;
        color: var(--color-complement-text);
    }

    th {
        width: 20%
    }

    td {
        width: 30%;
        padding-left: 0.5rem
    }
}
</style>