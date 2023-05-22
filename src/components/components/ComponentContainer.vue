<script setup>
import { defineProps, defineAsyncComponent } from 'vue';
import { useContentStore } from '../../store/contentStore';
import { useDialogStore } from '../../store/dialogStore';
import ComponentTag from '../utilities/ComponentTag.vue';
import LoadingImage from '../LoadingImage.vue';
import { indexToChartType } from '../../assets/configs/temp';

const contentStore = useContentStore()
const dialogStore = useDialogStore()

const props = defineProps({
    content: { type: Object },
    notMoreInfo: { type: Boolean, default: true },
})

</script>

<template>
    <div :class="{ 'componentcontainer': true, 'moreinfostyle': !notMoreInfo }">
        <div class="componentcontainer-header">
            <div>
                <h3>{{ content.name }}</h3>
                <h4>{{ `${content.source_from} | (顯示資料之時間)` }}</h4>
            </div>
            <div v-if="notMoreInfo">
                <button title="回報問題"><span>flag</span></button>
                <button v-if="contentStore.currentDashboard.type === 'customized'"
                    @click="contentStore.deleteComponent(content.topic_id, content.id, content.name)"><span>delete</span></button>
            </div>
        </div>

        <div class="componentcontainer-chart">
            <LoadingImage v-if="content.chartData.length === 0" />
            <component v-else-if="indexToChartType[content.index]" :is="indexToChartType[content.index]" :content="content">
            </component>
            <p v-else>{{ content.index }}{{ content.chartData }}</p>
        </div>
        <div class="componentcontainer-footer">
            <div>
                <ComponentTag :icon="``" :text="`更新頻率`" />
                <ComponentTag v-if="!content.map_config" :icon="`wrong_location`" :text="`沒有地圖`" />
                <ComponentTag v-if="content.calculation_config" :icon="`insights`" :text="`有歷史軸`" />
            </div>
            <button v-if="notMoreInfo" @click="dialogStore.showMoreInfo(content)">
                <p>組件資訊</p>
                <span>arrow_circle_right</span>
            </button>
        </div>
    </div>
</template>

<style scoped lang="scss">
.componentcontainer {
    background-color: var(--color-component-background);
    border-radius: 5px;
    height: 330px;
    max-height: 330px;
    width: calc(100% - var(--font-m) *2);
    max-width: calc(100% - var(--font-m) *2);
    padding: var(--font-m);
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    @media (min-width: 1050px) {
        height: 370px;
        max-height: 370px;
    }


    @media (min-width: 1650px) {
        height: 400px;
        max-height: 400px;
    }

    @media (min-width: 2200px) {
        height: 500px;
        max-height: 500px;
    }

    &-header {
        display: flex;
        justify-content: space-between;

        h3 {
            font-size: var(--font-m);
        }

        h4 {
            color: var(--color-complement-text);
            font-size: var(--font-s);
            font-weight: 400;
        }

        button span {
            font-family: var(--font-icon);
            font-size: calc(var(--font-l) * var(--font-to-icon));
            color: var(--color-complement-text);
            transition: color 0.2s;

            &:hover {
                color: white
            }
        }
    }

    &-chart {
        height: 70%;
        overflow-y: scroll;
        position: relative;

        p {
            color: var(--color-border)
        }
    }

    &-footer {
        display: flex;
        align-items: center;
        justify-content: space-between;

        div {
            display: flex;
            align-items: center;
        }

        button {
            display: flex;
            align-items: center;

            span {
                font-family: var(--font-icon);
                color: var(--color-highlight);
                margin-left: 4px
            }

            p {
                color: var(--color-highlight);
            }
        }

    }
}

.moreinfostyle {
    height: 350px;
    max-height: 350px;

    @media (min-width: 820px) {
        height: 380px;
        max-height: 380px;
    }

    @media (min-width: 1200px) {
        height: 420px;
        max-height: 420px;
    }

    @media (min-width: 2200px) {
        height: 520px;
        max-height: 520px;
    }

}
</style>