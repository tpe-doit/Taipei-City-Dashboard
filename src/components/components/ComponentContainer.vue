<script setup>
import { useContentStore } from '../../store/contentStore';
import { useDialogStore } from '../../store/dialogStore';
import ComponentTag from '../utilities/ComponentTag.vue';
import { computed, ref } from 'vue';
import { chartTypes } from '../../assets/configs/apexcharts/chartTypes';

const contentStore = useContentStore()
const dialogStore = useDialogStore()

const props = defineProps({
    content: { type: Object },
    notMoreInfo: { type: Boolean, default: true },
    isMapLayer: { type: Boolean, default: false },
})

const dataTime = computed(() => {
    if (!props.content.time_from) {
        return '固定資料'
    }
    if (!props.content.time_to) {
        return props.content.time_from.slice(0, 10)
    }
    return `${props.content.time_from.slice(0, 10)} ~ ${props.content.time_to.slice(0, 10)}`
})

const updateFreq = computed(() => {
    const unitRef = {
        minute: "分鐘",
        hour: "小時",
        day: "天",
        week: "週",
        month: "月",
        year: "年"
    }
    if (!props.content.update_freq) {
        return '不定期更新'
    }
    return `每${props.content.update_freq}${unitRef[props.content.update_freq_unit]}更新`
})

const activeChart = ref(props.content.chart_config.types[0])

function changeActiveChart(chartName) {
    activeChart.value = chartName
}

</script>

<template>
    <div :class="{ 'componentcontainer': true, 'moreinfostyle': !notMoreInfo, 'maplayerstyle': isMapLayer }">
        <div class="componentcontainer-header">
            <div>
                <h3>{{ content.name }}</h3>
                <h4>{{ `${content.source} | ${dataTime}` }}</h4>
            </div>
            <div v-if="notMoreInfo">
                <button title="回報問題"><span>flag</span></button>

                <!-- If you wish to implement the "delete component" function, uncomment the following -->
                <!-- <button v-if="contentStore.currentDashboard.type === 'customized'"
                    @click="contentStore.deleteComponent(content.topic_id, content.id, content.name)"><span>delete</span></button> -->

            </div>
        </div>
        <div class="componentcontainer-control" v-if="props.content.chart_config.types.length > 1">
            <button v-for="item in props.content.chart_config.types" @click="changeActiveChart(item)">{{
                chartTypes[item] }}</button>
        </div>
        <div :class="{ 'componentcontainer-chart': true, 'maplayerstyle-chart': isMapLayer }" v-if="content.chart_data">
            <component v-for="item in content.chart_config.types" :activeChart="activeChart" :is="item"
                :chart_config="content.chart_config" :series="content.chart_data">
            </component>
        </div>
        <div class="componentcontainer-footer">
            <div>
                <ComponentTag :icon="``" :text="updateFreq" />
                <ComponentTag v-if="!content.map_config" :icon="`wrong_location`" :text="`沒有地圖`" />
            </div>
            <button v-if="notMoreInfo && !isMapLayer" @click="dialogStore.showMoreInfo(content)">
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
    position: relative;

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

    &-control {
        position: absolute;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: 4rem;
        left: 0;
        z-index: 8;

        button {
            background-color: rgb(77, 77, 77);
            padding: 4px 4px;
            border-radius: 5px;
            transition: color 0.2s, opacity 0.2s;
            font-size: var(--font-s);
            margin: 0 4px;
            color: var(--color-complement-text);
            opacity: 0.25;
            text-align: center;

            &:hover {
                color: white;
                opacity: 1;
            }
        }
    }

    &-chart {
        height: 75%;
        overflow-y: scroll;
        position: relative;
        padding-top: 5%;

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

.maplayerstyle {
    height: 180px;
    max-height: 180px;

    @media (min-width: 1050px) {
        height: 210px;
        max-height: 210px;
    }


    @media (min-width: 1650px) {
        height: 225px;
        max-height: 225px;
    }

    @media (min-width: 2200px) {
        height: 275px;
        max-height: 275px;
    }

    &-chart {
        height: 60%;
    }
}
</style>