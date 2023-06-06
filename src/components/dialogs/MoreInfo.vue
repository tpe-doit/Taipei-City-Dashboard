<!-- Cleaned -->

<script setup>
import { useDialogStore } from '../../store/dialogStore';

import DialogContainer from './DialogContainer.vue';
import ComponentContainer from '../components/ComponentContainer.vue';
import HistoryChart from '../utilities/HistoryChart.vue';

const dialogStore = useDialogStore()
</script>

<template>
    <DialogContainer :dialog="`moreInfo`" @on-close="dialogStore.hideAllDialogs">
        <div class="moreinfo">
            <ComponentContainer :content="dialogStore.moreInfoContent" :not-more-info="false" />
            <div class="moreinfo-info">
                <h3>組件說明</h3>
                <p>{{ dialogStore.moreInfoContent.long_desc }}</p>
                <h3>使用情境</h3>
                <p>{{ dialogStore.moreInfoContent.use_case }}</p>
                <div v-if="dialogStore.moreInfoContent.links">
                    <h3>相關連結</h3>
                    <a>連結</a>
                </div>
                <div v-if="dialogStore.moreInfoContent.history_data">
                    <h3>歷史軸</h3>
                    <h4>*點擊並拉動以檢視細部區間資料</h4>
                    <HistoryChart :chart_config="dialogStore.moreInfoContent.chart_config"
                        :series="dialogStore.moreInfoContent.history_data" />
                </div>
                <!-- <div class="moreinfo-info-control">
                    <button><span>flag</span>回報問題</button>
                    <button><span>download</span>輸出</button>
                </div> -->
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
        display: flex;
        flex-direction: column;
        padding: 1rem;
        border-top: solid 1px var(--color-border);

        p {
            margin-bottom: 0.75rem;
            text-align: justify;
        }

        h4 {
            color: var(--color-complement-text);
            font-weight: 400;
            font-size: 10px;
        }

        @media (min-width: 820px) {
            border-left: solid 1px var(--color-border);
            border-top: none;
        }

        &-control {
            display: flex;
            align-items: flex-end;
            justify-content: flex-end;
            flex: 1;

            span {
                margin-right: 4px;
                font-family: var(--font-icon);
                font-size: calc(var(--font-m) * var(--font-to-icon));
            }

            button {
                display: flex;
                align-items: center;
                margin-left: 8px;
                padding: 2px 4px;
                border-radius: 5px;
                background-color: var(--color-highlight);
                font-size: var(--font-m);
            }
        }
    }
}
</style>