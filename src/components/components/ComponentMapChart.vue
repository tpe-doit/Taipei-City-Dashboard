<!-- Cleaned -->

<!-- This component has two modes 'normal mapview charts' / 'basic map layers' -->
<!-- The different modes are controlled by the prop "isMapLayer" (default false) -->

<script setup>
import { ref, computed } from 'vue';
import { useMapStore } from '../../store/mapStore';

import { chartTypes } from '../../assets/configs/apexcharts/chartTypes';

const mapStore = useMapStore()

const props = defineProps({
    // The complete config (incl. chart data) of a dashboard component will be passed in
    content: { type: Object },
    isMapLayer: { type: Boolean, default: false },
})

// The default active chart is the first one in the list defined in the dashboard component
const activeChart = ref(props.content.chart_config.types[0])
// Stores whether the component is toggled on or not
const checked = ref(false)

// Parses time data into display format
const dataTime = computed(() => {
    if (!props.content.time_from) {
        return '固定資料'
    }
    if (!props.content.time_to) {
        return props.content.time_from.slice(0, 10)
    }
    return `${props.content.time_from.slice(0, 10)} ~ ${props.content.time_to.slice(0, 10)}`
})

// Open and closes the component as well as communicates to the mapStore to turn on and off map layers
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
// Toggles between chart types defined in the dashboard component
function changeActiveChart(chartName) {
    activeChart.value = chartName
}
</script>

<template>
    <div :class="{ componentmapchart: true, checked: checked, 'maplayer': isMapLayer && checked }">
        <div class="componentmapchart-header">
            <div>
                <div>
                    <h3>{{ content.name }}</h3>
                    <span v-if="content.map_config">map</span>
                    <span v-if="content.history_data">insights</span>
                </div>
                <h4 v-if="checked">{{ `${content.source} | ${dataTime}` }}</h4>
            </div>
            <!-- The class "toggleswitch" could be edited in /assets/styles/toggleswitch.css -->
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
            <!-- The components referenced here can be edited in /components/charts -->
            <component v-for="item in content.chart_config.types" :activeChart="activeChart" :is="item"
                :chart_config="content.chart_config" :series="content.chart_data">
            </component>
        </div>
    </div>
</template>

<style scoped lang="scss">
.componentmapchart {
    width: calc(100% - var(--font-m) *2);
    max-width: calc(100% - var(--font-m) *2);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    position: relative;
    padding: var(--font-m);
    border-radius: 5px;
    background-color: var(--color-component-background);

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
                margin-left: 8px;
                color: var(--color-complement-text);
                font-family: var(--font-icon);
                user-select: none;
                pointer-events: none;
            }
        }
    }

    &-control {
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: 4rem;
        left: 0;
        z-index: 10;

        button {
            margin: 0 4px;
            padding: 4px 4px;
            border-radius: 5px;
            background-color: rgb(77, 77, 77);
            opacity: 0.25;
            color: var(--color-complement-text);
            font-size: var(--font-s);
            text-align: center;
            ;
            transition: color 0.2s, opacity 0.2s;

            &:hover {
                opacity: 1;
                color: white;
            }
        }
    }

    &-chart {
        height: 80%;
        position: relative;
        overflow-y: scroll;

        p {
            color: var(--color-border)
        }
    }
}

.checked {
    max-height: 300px;
    height: 300px;
}

.maplayer {
    height: 200px;
    max-height: 200px;
    padding-bottom: 0;
}
</style>