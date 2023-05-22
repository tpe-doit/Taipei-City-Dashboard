<script setup>
import { defineProps, ref, defineAsyncComponent } from 'vue';
import { useMapStore } from '../../store/mapStore';
import LoadingImage from '../LoadingImage.vue';
import { indexToChartType } from '../../assets/configs/temp';

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
        mapStore.addToMapLayerList(props.content.map_config[0])
    } else {
        mapStore.turnOffMapLayerVisibility(props.content.map_config[0])
    }
}

</script>

<template>
    <div :class="{ componentmapchart: true, checked: checked }">
        <div class="componentmapchart-header">
            <div>
                <div>
                    <h3>{{ content.name }}</h3>
                    <span v-if="!content.map_config">wrong_location</span>
                    <span v-if="content.calculation_config">insights</span>
                </div>
                <h4 v-if="checked">{{ `${content.source_from} | (顯示資料之時間)` }}</h4>
            </div>
            <label class="toggleswitch">
                <input type="checkbox" @change="handleToggle" v-model="checked">
                <span class="toggleswitch-slider"></span>
            </label>
        </div>
        <div class="componentmapchart-chart" v-if="checked">
            <LoadingImage v-if="content.chartData.length === 0" />
            <component v-else-if="indexToChartType[content.index]" :is="indexToChartType[content.index]" :content="content">
            </component>
            <p v-else>{{ content.index }}{{ content.chartData }}</p>
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