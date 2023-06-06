<!-- Cleaned -->

<script setup>
import { ref, computed } from 'vue';

const props = defineProps(['chart_config', 'activeChart', 'series'])

const chartOptions = ref({
    chart: {
        stacked: true,
        stackType: '100%',
        toolbar: {
            show: false
        },
    },
    colors: props.series.length > 2 ? props.chart_config.color : [props.chart_config.color[0], "#777"],
    dataLabels: {
        textAnchor: 'start',
    },
    grid: {
        show: false,
    },
    legend: {
        offsetY: 20,
        position: "top",
        show: props.series.length > 2 ? true : false,
    },
    plotOptions: {
        bar: {
            borderRadius: 5,
            horizontal: true,
        }
    },
    stroke: {
        colors: ['#282a2c'],
        show: true,
        width: 2,
    },
    tooltip: {
        // The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.seriesNames[seriesIndex] + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        },
    },
    xaxis: {
        axisBorder: {
            show: false
        },
        axisTicks: {
            show: false,
        },
        categories: props.chart_config.categories ? props.chart_config.categories : [],
        labels: {
            show: false,
        },
        type: 'category',
    },
})

const chartHeight = computed(() => {
    return `${40 + props.series[0].data.length * 30}`
})

</script>

<template>
    <div v-if="activeChart === 'BarPercentChart'">
        <apexchart width="100%" :height="chartHeight" type="bar" :options="chartOptions" :series="series"></apexchart>
    </div>
</template>