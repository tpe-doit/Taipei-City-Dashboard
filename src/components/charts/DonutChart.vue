<script setup>
import { computed, defineProps, ref } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const steps = ref(6) // How many data points to show before 'other'

const parsedSeries = computed(() => {
    const toParse = [...props.series[0].data]
    if (toParse.length <= steps.value) {
        return toParse.map((item) => item.y)
    }
    let output = []
    for (let i = 0; i < steps.value; i++) {
        output.push(toParse[i].y)
    }
    const toSum = toParse.splice(steps.value, toParse.length - steps.value)
    let sum = 0;
    toSum.forEach(element => sum += element.y)
    output.push(sum)
    return output
})
const parsedLabels = computed(() => {
    const toParse = [...props.series[0].data]
    if (toParse.length <= steps.value) {
        return toParse.map((item) => item.x)
    }
    let output = []
    for (let i = 0; i < steps.value; i++) {
        output.push(toParse[i].x)
    }
    output.push('其他')
    return output
})
const sum = computed(() => {
    return Math.round(parsedSeries.value.reduce((a, b) => a + b) * 100) / 100
})
const options = ref({
    labels: parsedLabels,
    legend: {
        show: false
    },
    stroke: {
        show: true,
        colors: ['#282a2c'],
        width: 3,
    },
    tooltip: {
        followCursor: false,
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.labels[seriesIndex] + '</h6>' +
                '<span>' + series[seriesIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        },
    },
    dataLabels: {
        formatter: function (val, { seriesIndex, dataPointIndex, w }) {
            let value = w.globals.labels[seriesIndex]
            return value.length > 7 ? value.slice(0, 6) + "..." : value
        }
    },
    plotOptions: {
        pie: {
            dataLabels: {
                offset: 15,
            },
            donut: {
                size: '80%',
            }
        }
    },
    colors: props.series.length >= steps.value ? [...props.chart_config.color, '#848c94'] : props.chart_config.color,
    chart: {
        offsetY: 15,
    }
})
</script>

<template>
    <div v-if="activeChart === 'DonutChart'" class="donutchart">
        <apexchart width="100%" type="donut" :options="options" :series="parsedSeries"></apexchart>
        <div class="donutchart-title">
            <h5>總合</h5>
            <h6>{{ sum }}</h6>
        </div>
    </div>
</template>

<style scoped lang="scss">
.donutchart {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    width: 100%;
    overflow-y: visible;
    position: relative;

    &-title {
        position: absolute;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;

        h5 {
            color: var(--color-complement-text)
        }

        h6 {
            font-size: var(--font-m);
            font-weight: 400;
            color: var(--color-complement-text)
        }
    }
}
</style>