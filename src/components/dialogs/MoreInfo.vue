<script setup>
import DialogContainer from './DialogContainer.vue';
import { useDialogStore } from '../../store/dialogStore';
import ComponentContainer from '../components/ComponentContainer.vue';

const dialogStore = useDialogStore()

</script>

<template>
    <DialogContainer :dialog="`moreInfo`" @on-close="dialogStore.hideAllDialogs">
        <div class="moreinfo">
            <ComponentContainer :content="dialogStore.moreInfoContent" :not-more-info="false" />
            <div class="moreinfo-info">
                <h3>組件說明</h3>
                <p>本組件搭配啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉。</p>
                <h3>使用情境</h3>
                <p>本組件搭配啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉，啦啦啦拉拉。</p>
                <h3>相關連結</h3>
                <p>連結</p>
                <div class="moreinfo-info-control">
                    <button><span>flag</span>回報問題</button>
                    <button><span>download</span>輸出</button>
                </div>
            </div>
        </div>

    </DialogContainer>
</template>

<style scoped lang="scss">
.moreinfo {
    height: fit-content;
    width: 400px;
    display: grid;

    @media (min-width: 820px) {
        width: 700px;
        height: 410px;
        grid-template-columns: 3fr 2fr;
    }

    @media (min-width: 1200px) {
        height: 440px;
        width: 800px;
    }

    @media (min-width: 2200px) {
        height: 550px;
        width: 900px;
    }

    &-info {
        padding: 1rem;
        border-top: solid 1px var(--color-border);
        display: flex;
        flex-direction: column;

        p {
            text-align: justify;
            margin-bottom: 0.75rem;
        }

        @media (min-width: 820px) {
            border-left: solid 1px var(--color-border);
            border-top: none;
        }

        &-control {
            flex: 1;
            display: flex;
            align-items: flex-end;
            justify-content: flex-end;

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
                background-color: var(--color-highlight);
                padding: 2px 4px;
            }

        }
    }
}
</style>