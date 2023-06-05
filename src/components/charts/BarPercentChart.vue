<script setup>
import { ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const options = ref({
    legend: {
        show: props.series.length > 2 ? true : false,
        position: "top",
        offsetY: 20
    },
    stroke: {
        show: true,
        colors: ['#282a2c'],
        width: 2,
    },
    tooltip: {
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.seriesNames[seriesIndex] + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        }
    },
    dataLabels: {
        textAnchor: 'start',
    },
    plotOptions: {
        bar: {
            horizontal: true,
            borderRadius: 5,
        }
    },
    colors: props.series.length > 2 ? props.chart_config.color : [props.chart_config.color[0], "#777"],
    grid: {
        show: false,
    },
    chart: {
        toolbar: {
            show: false
        },
        stacked: true,
        stackType: '100%'
    },
    xaxis: {
        type: 'category',
        categories: props.chart_config.categories ? props.chart_config.categories : [],
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
    return `${40 + props.series[0].data.length * 30}`
})

</script>

<template>
    <div v-if="activeChart === 'BarPercentChart'">
        <apexchart width="100%" :height="height" type="bar" :options="options" :series="series"></apexchart>
    </div>
</template>

<style scoped></style>