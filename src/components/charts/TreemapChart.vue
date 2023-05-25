<script setup>
import { defineProps, ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const sum = computed(() => {
    let sum = 0;
    props.series.forEach(item => sum += item.y)
    return Math.round(sum * 100) / 100
})

const options = ref({
    legend: {
        show: false
    },
    stroke: {
        show: true,
        colors: ['#282a2c'],
        width: 2,
    },
    tooltip: {
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.categoryLabels[dataPointIndex] + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ' km2' + '</span>' +
                '</div>'
        }
    },
    dataLabels: {
        formatter: function (val, { seriesIndex, dataPointIndex, w }) {
            return dataPointIndex > 5 ? '' : val
        }
    },
    plotOptions: {
        treemap: {
            distributed: true,
            shadeIntensity: 0,
        }
    },
    colors: props.chart_config.color,
    grid: {
        show: false,
    },
    chart: {
        toolbar: {
            show: false
        }
    },
    xaxis: {
        type: 'category',
        labels: {
            show: false,
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false
        }
    },
})

</script>

<template>
    <div v-if="activeChart === 'TreemapChart'" class="treemapchart">
        <div class="treemapchart-title">
            <h5>面積總合</h5>
            <h6>{{ sum }} km2</h6>
        </div>
        <apexchart width="100%" type="treemap" :options="options" :series="[{ data: series }]"></apexchart>
    </div>
</template>

<style scoped lang="scss">
.treemapchart {

    &-title {
        display: flex;
        justify-content: center;
        flex-direction: column;
        margin: 0.5rem 0 -0.5rem;

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