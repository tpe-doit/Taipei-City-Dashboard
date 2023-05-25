<script setup>
import { defineProps, ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const options = ref({
    legend: {
        show: false
    },
    stroke: {
        show: true,
        colors: ['#282a2c'],
        width: 0,
    },
    tooltip: {
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.labels[dataPointIndex] + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + '</span>' +
                '</div>'
        }
    },
    dataLabels: {
        textAnchor: 'start',
        offsetX: 20
    },
    plotOptions: {
        bar: {
            horizontal: true,
            borderRadius: 2,
            distributed: true
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

const height = computed(() => {
    return `${40 + props.series.length * 16}`
})

</script>

<template>
    <div v-if="activeChart === 'HorizontalBarChart'">
        <apexchart width="100%" :height="height" type="bar" :options="options" :series="[{ data: series }]"></apexchart>
    </div>
</template>

<style scoped></style>