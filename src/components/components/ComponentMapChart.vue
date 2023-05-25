<script setup>
import { ref, computed } from 'vue';
import { useMapStore } from '../../store/mapStore';
import LoadingImage from '../LoadingImage.vue';
import { chartTypes } from '../../assets/configs/apexcharts/chartTypes';

const props = defineProps({
    content: { type: Object }
})
const mapStore = useMapStore()

const checked = ref(false)

function handleToggle() {
    if (!props.content.map_config) {
        return
    }
    if (checked.value) {
        mapStore.addToMapLayerList(props.content.map_config)
    } else {
        mapStore.turnOffMapLayerVisibility(props.content.map_config)
    }
}

const dataTime = computed(() => {
    if (!props.content.time_from) {
        return '固定資料'
    }
    if (!props.content.time_to) {
        return props.content.time_from.slice(0, 10)
    }
    return `${props.content.time_from.slice(0, 10)} ~ ${props.content.time_to.slice(0, 10)}`
})

const activeChart = ref(props.content.chart_config.types[0])

function changeActiveChart(chartName) {
    activeChart.value = chartName
}

</script>

<template>
    <div :class="{ componentmapchart: true, checked: checked }">
        <div class="componentmapchart-header">
            <div>
                <div>
                    <h3>{{ content.name }}</h3>
                    <span v-if="!content.map_config">wrong_location</span>
                </div>
                <h4 v-if="checked">{{ `${content.source} | ${dataTime}` }}</h4>
            </div>
            <label class="toggleswitch">
                <input type="checkbox" @change="handleToggle" v-model="checked">
                <span class="toggleswitch-slider"></span>
            </label>
        </div>
        <div class="componentmapchart-control" v-if="props.content.chart_config.types.length > 1 && checked">
            <button v-for="item in props.content.chart_config.types" @click="changeActiveChart(item)">{{
                chartTypes[item] }}</button>
        </div>
        <div class="componentmapchart-chart" v-if="checked && content.chart_data">
            <component v-for="item in content.chart_config.types" :activeChart="activeChart" :is="item"
                :chart_config="content.chart_config" :series="content.chart_data">
            </component>
        </div>

        <div class="componentmapchart-footer"></div>
    </div>
</template>

<style scoped lang="scss">
.checked {
    max-height: 300px;
    height: 300px;
}

.componentmapchart {
    background-color: var(--color-component-background);
    border-radius: 5px;
    width: calc(100% - var(--font-m) *2);
    max-width: calc(100% - var(--font-m) *2);
    padding: var(--font-m);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    position: relative;

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

        div:first-child {
            div {
                display: flex;
                align-items: center;
            }

            span {
                font-family: var(--font-icon);
                color: var(--color-complement-text);
                margin-left: 8px;
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
        z-index: 10;

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
        height: 80%;
        overflow-y: scroll;
        position: relative;

        p {
            color: var(--color-border)
        }
    }


}

.toggleswitch {
    position: relative;
    display: inline-block;
    height: 1rem;
    width: 2rem;
    margin-top: 4px;

    input {
        opacity: 0;
        width: 0;
        height: 0;
    }

    &-slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: var(--color-complement-text);
        -webkit-transition: .4s;
        transition: .4s;
        border-radius: var(--font-m);
    }

    &-slider::before {
        position: absolute;
        content: "";
        height: 0.7rem;
        width: 0.7rem;
        left: 0.15rem;
        bottom: 0.15rem;
        background-color: var(--color-border);
        -webkit-transition: .4s;
        transition: .4s;
        border-radius: 50%;
    }

    input:checked+&-slider {
        background-color: var(--color-highlight);
    }

    input:focus+&-slider {
        box-shadow: 0 0 1px var(--color-highlight);
    }

    input:checked+&-slider:before {
        -webkit-transform: translateX(1rem);
        -ms-transform: translateX(1rem);
        transform: translateX(1rem);
        background-color: white;
    }
}
</style>